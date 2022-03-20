package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/when-are-you-free/backend/auth"
	"github.com/when-are-you-free/backend/meetup_endpoints"
	"github.com/when-are-you-free/backend/storage"
)

func main() {
	storage, errStorage := loadStorage()
	if errStorage != nil {
		return
	}
	defer persistStorage(storage)

	go serve(storage)

	log.Println("Server is now running. Press CTRL-C to exit.")
	killChannel := make(chan os.Signal, 1)
	signal.Notify(killChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-killChannel
	log.Println("Shutting down ...")
}

func loadStorage() (*storage.Storage, error) {
	log.Println("Loading storage ...")
	storage := storage.NewStorage()

	readOnlyPersistedStorage, errOpen := os.Open("storage.json")
	if errOpen != nil {
		//No data yet, this must be the first start.
		if os.IsNotExist(errOpen) {
			log.Println("No data could be found, assuming this is the first startup.")
			return storage, nil
		}

		log.Fatalln("Error opening file containing persisted storage:", errOpen)
		return nil, errOpen
	}

	if errLoad := storage.Load(readOnlyPersistedStorage); errLoad != nil {
		log.Fatalln("Error loading persisted storage:", errLoad)
		return nil, errLoad
	}

	log.Println("Storage loaded")
	return storage, nil
}

func persistStorage(storage *storage.Storage) {
	log.Println("Persisting storage ...")
	writablePersistedStorage, errOpen := os.OpenFile("storage.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if errOpen != nil {
		log.Fatalln("Error opening file containing persisted storage:", errOpen)
	}

	if errPersist := storage.Persist(writablePersistedStorage); errPersist != nil {
		log.Fatalln("Error persisting storage:", errPersist)
	}
	log.Println("Persisted storage")
}

func serve(storage *storage.Storage) {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "X-User-Token", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           0,
	}))
	router.Group(func(router chi.Router) {
		router.Use(auth.Auth)

		meetupHandler := meetup_endpoints.New(storage)
		router.Post("/meetup", meetupHandler.Post)
		router.Get("/meetup/{uuid}", meetupHandler.Get)
		router.Patch("/meetup/{uuid}", meetupHandler.Patch)
		router.Delete("/meetup/{uuid}", meetupHandler.Delete)
	})

	log.Println("Listening ...")
	http.ListenAndServe(":8080", router)
}
