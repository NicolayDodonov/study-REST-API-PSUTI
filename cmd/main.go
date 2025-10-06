package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"study-REST-API-PSUTI/internal/config"
	"study-REST-API-PSUTI/internal/storage"
	"syscall"

	"study-REST-API-PSUTI/internal/handler"
	"study-REST-API-PSUTI/internal/logger"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	configPath = "config/config.yaml"
)

func main() {
	// прочитать конфиг
	cnf := config.MustLoad(configPath)

	// активируем логгер
	log, _ := logger.New(cnf.Log.Path, cnf.Log.Level)

	// создать подключение к бд
	db, err := sqlx.Connect("postgres", cnf.Postgres.DSN())
	if err != nil {
		log.Fatal(err.Error())
	}

	// создать обработчик
	h := handler.New(storage.New(db), log)
	_ = h
	// расписать роутер и адреса апишки
	r := chi.NewRouter()
	r.Use(chimid.RequestID)                 // X-Request-ID + в контексте
	r.Use(chimid.RealIP)                    // реальный IP клиента
	r.Use(chimid.Recoverer)                 // panic → 500
	r.Use(chimid.Timeout(cnf.HTTP.Timeout)) // таймаут на обработку запроса
	r.Use(chimid.StripSlashes)              // нормализация путей
	r.Use(chimid.Compress(5))               // gzip/deflate с уровнем 5

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
		r.Get("/users", h.GetUser)
		r.Put("/add_info", h.UpdateUser)
		r.Delete("/delete_info", h.DeleteUser)
	})

	srv := http.Server{
		Addr:         cnf.HTTP.Host + ":" + cnf.HTTP.Port,
		Handler:      r,
		ReadTimeout:  cnf.HTTP.ReadTimeout,
		WriteTimeout: cnf.HTTP.WriteTimeout,
		IdleTimeout:  cnf.HTTP.IdleTimeout,
	}

	go func() {
		log.Info(fmt.Sprintf("Listening and serving HTTP on port %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server: " + err.Error())
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cnf.HTTP.Timeout)
	defer cancel()
	// передаём в сервер контекст для выхода
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown error")
	}
	log.Info("server shutting down")
}
