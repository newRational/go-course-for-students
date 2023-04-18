package main

import (
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	"homework9/internal/ports/httpgin"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), userrepo.New()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
