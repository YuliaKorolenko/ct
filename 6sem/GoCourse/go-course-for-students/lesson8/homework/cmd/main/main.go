package main

import (
	"homework8/internal/adapters/adrepo"
	"homework8/internal/adapters/usrepo"
	"homework8/internal/app"
	"homework8/internal/ports/httpgin"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), usrepo.NewUserRepo()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
