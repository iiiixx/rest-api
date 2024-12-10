package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"rest-api/internal/user"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler — обработчик запросов
func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))

}

func main() {
	log.Println("create router")
	router := httprouter.New()

	log.Println("register user handler")
	handler := user.NewHandler()
	handler.Register(router)

	start(router)

}

// Настройка и запуск сервера
func start(router *httprouter.Router) {
	log.Println("start application")

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	//Настройка таймаутов и запуск HTTP-сервера
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("server is listening port 0.0.0.0:1234")
	log.Fatalln(server.Serve(listener))
}
