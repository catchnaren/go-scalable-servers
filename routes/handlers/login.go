package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/catchnaren/go-scalable-servers/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauth2Config *oauth2.Config
var oauthStateString = "go-server"

func init() {
	googleOauth2Config = &oauth2.Config{
		ClientID:     config.Config.GoogleClientID,
		ClientSecret: config.Config.GoogleClientSecret,
		RedirectURL: config.Config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleLogin(ctx *gin.Context) {
	url := googleOauth2Config.AuthCodeURL("go-server", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusFound, url)
}

func HandleGoogleCallback(ctx *gin.Context) {
	// Validate the state
	state := ctx.Query("state")
	if state != oauthStateString {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}
	
	// Validate the code and get usertoken back
	code := ctx.Query("code")
	token, err := googleOauth2Config.Exchange(ctx, code)
	if err != nil {
		// log.Sugar.Errorf("Error while exchange code for token:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Get user info using userToken
	client := googleOauth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		// log.Sugar.Errorf("Error retrieving user info:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Send userinfo to frontend
		//  1. As an object
		// *2. As a token (which will expire) - JWT token
	var userInfo struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		// log.Sugar.Errorf("Error decoding user info:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}
	// log.Printf(userInfo.Email)
	jwtToken, err := generateJWT(userInfo.Email, userInfo.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"email": userInfo.Email,
		"name": userInfo.Name,
		"picture": userInfo.Picture,
	})
}

func generateJWT(email, name string) (string, error) {
	tokenInfo := jwt.MapClaims{
		"email": email,
		"name": name,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenInfo) 
	return token.SignedString([]byte(config.Config.JWTSaltKey))
}