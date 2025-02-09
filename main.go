package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hungour-streaming-server/controller"
)

func main() {
	http.HandleFunc("/process/", controller.ProcessController)
	http.HandleFunc("/callback/", controller.CallbackController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("初回通話開始")
	controller.StartController()

	log.Println("サーバ起動")
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", port), nil))
}
