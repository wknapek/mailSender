package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/madflojo/tasks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"vodeno/handlers"
	"vodeno/model"
	"vodeno/worker"
)

func main() {
	httpServerExitDone := &sync.WaitGroup{}
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()
	router := chi.NewRouter()
	log.Info().Msg("app running")
	dsn := "host=localhost user=adm_user password=S3cret dbname=vodeno port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Msg("failed to connect to database details: " + err.Error())
		return
	}
	InitTable(db)
	handl := handlers.New(db)
	router.Use(middleware.Logger)
	router.Post("/api/messages", handl.CreateEmail)
	router.Post("/api/messages/send", handl.SendEmail)
	router.Delete("/api/messages/{id}", handl.DeleteEmail)
	noSSLSrv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	httpServerExitDone.Add(1)
	go func() {
		defer httpServerExitDone.Done()
		if err = noSSLSrv.ListenAndServe(); err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	work := worker.NewWorker(db)
	scheduler := tasks.New()
	defer scheduler.Stop()
	id, err := scheduler.Add(&tasks.Task{
		TaskContext: tasks.TaskContext{},
		Interval:    time.Duration(5 * time.Minute),
		TaskFunc:    work.Deleter,
	})
	log.Info().Msg("start deleting task with id:" + id)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err = noSSLSrv.Shutdown(ctx); err != nil {
		panic(err)
	}
	httpServerExitDone.Wait()
	fmt.Println("server ending work")
}

func InitTable(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&model.Email{})
	if err != nil {
		panic(err)
	}
}
