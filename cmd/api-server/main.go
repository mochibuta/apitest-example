package main

import (
	"log"

	apiserver "github.com/mochibuta/apitest-example/cmd/api-server/server"
)

func main() {
	srv, err := apiserver.InitServer()
	if err != nil {
		log.Fatal(err)
	}

	srv.Run(":8080")

	defer apiserver.CloseDB()

}
