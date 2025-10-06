package main

import (
	""
	"fmt"
	"net/http"
	"study-REST-API-PSUTI/internal/handler"
	"study-REST-API-PSUTI/internal/logger"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// прочитать конфиг

	// активируем логгер
	log, _ := logger.New("log/log.txt", "debug")

	// создать подключение к бд

	// расписать роутер и адреса апишки
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.TODO)
	})

	srv := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info(fmt.Sprintf("Listening and serving HTTP on port %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server: " + err.Error())
		}
	}()
}
