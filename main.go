package main

import (
	"fmt"
	"go-mongo/controller"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Go-MongoAPI")
	fmt.Println("Server is getting stated...")

	router := mux.NewRouter()

	// Health route

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Health check route... Working")
	}).Methods("GET")

	// MiddleWare handle -- Source code renders

	router.HandleFunc("/panic/", panicDemo)
	router.HandleFunc("/debug/", controller.SourceCodeHandler).Methods("GET")

	router.HandleFunc("/api/movies", controller.GetMyAllMovies).Methods("GET")
	router.HandleFunc("/api/movie/{id}", controller.GetAMovieById).Methods("GET")

	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT")

	router.HandleFunc("/api/movie/{id}", controller.DeleteAMovie).Methods("DELETE")
	router.HandleFunc("/api/deleteallmovie", controller.DeleteAllMovies).Methods("DELETE")

	fmt.Println("Listening is at PORT : 4000")
	log.Fatal(http.ListenAndServe(":4000", devMw(router)))
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	panic("Oh no!")
}
func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				// log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, controller.MakeLinks(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}
