package main

import (
	"context"
	"net/http"

	"github.com/catchnaren/go-scalable-servers/config"
	"github.com/catchnaren/go-scalable-servers/db"
	"github.com/catchnaren/go-scalable-servers/routes"
)

func main() {
	handler := routes.MountRoutes()
	
	db.InitDB()
	
	server := &http.Server{
		Addr: config.Config.AppPort,
		Handler: handler,
	}
	
	defer db.DB.Close(context.Background()) //unmount
	
	server.ListenAndServe()
}