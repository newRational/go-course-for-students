package main

import (
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	"homework9/internal/ports/httpgin"
)

const port = ":18080"

func main() {
	server := httpgin.NewHTTPServer(port, app.NewApp(adrepo.New(), userrepo.New()))
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
