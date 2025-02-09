package main

import (
	"log"
	"net/http"
	"os"

	"hungour-streaming-server/controller"
)

func main() {
	http.HandleFunc("/process/", controller.ProcessController)
	http.HandleFunc("/callback/", controller.CallbackController)
	http.HandleFunc("/start/", controller.StartController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	log.Println("サーバ起動")

	log.Fatal(http.ListenAndServe(":8087", nil))
}
