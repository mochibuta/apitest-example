package main

import (
	"context"
	"log"

	apiserver "github.com/mochibuta/apitest-example/cmd/api-server/server"
)

func main() {
	ctx := context.Background()

	srv, err := apiserver.InitServer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run(":8080")

	defer apiserver.CloseDB()

}
