package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"rest-api/pkg/client/mongodb"
	"rest-api/pkg/logging"

	"rest-api/internal/config"
	"rest-api/internal/user"
	"rest-api/internal/user/db"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler — обработчик запросов
func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))

}

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
		cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	users, err := storage.FindAll(context.Background())
	fmt.Println(users)

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)

}

// Настройка и запуск сервера
func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenError error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(filepath.ErrBadPattern)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("create unix socket")
		listener, listenError = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket %s", socketPath)

	} else {
		logger.Info("tcp")
		listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenError != nil {
		logger.Fatal(listenError)
	}

	//Настройка таймаутов и запуск HTTP-сервера
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatalln(server.Serve(listener))
}
