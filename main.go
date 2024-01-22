package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Movie represents a movie entity
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director represents a director entity
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(res http.ResponseWriter, req *http.Request) {
	fmt.Print("getMovies called\n")
	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(movies)
}

func deleteMovie(res http.ResponseWriter, req *http.Request) {
	fmt.Print("deleteMovie called\n")
	res.Header().Set("Content-type", "application/json")
	id := chi.URLParam(req, "id")
	deleteMovieByID(id)
	json.NewEncoder(res).Encode(movies)
}

func getMovie(res http.ResponseWriter, req *http.Request) {
	fmt.Print("getMovieByID called\n")
	res.Header().Set("Content-type", "application/json")
	id := chi.URLParam(req, "id")
	movie := getMovieByID(id)
	if movie != nil {
		json.NewEncoder(res).Encode(movie)
		return
	}
	http.NotFound(res, req)
}

func createMovie(res http.ResponseWriter, req *http.Request) {
	fmt.Print("createMovie called\n")
	res.Header().Set("Content-type", "application/json")
	var movie Movie
	_ = json.NewDecoder(req.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	movies = append(movies, movie)
	json.NewEncoder(res).Encode(movie)
}

func updateMovie(res http.ResponseWriter, req *http.Request) {
	fmt.Print("updateMovie called\n")
	res.Header().Set("Content-type", "application/json")
	id := chi.URLParam(req, "id")
	deleteMovieByID(id)
	var movie Movie
	_ = json.NewDecoder(req.Body).Decode(&movie)
	movie.ID = id
	movies = append(movies, movie)
	json.NewEncoder(res).Encode(movie)
}

func deleteMovieByID(id string) {
	for index, value := range movies {
		if value.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
}

func getMovieByID(id string) *Movie {
	for _, value := range movies {
		if value.ID == id {
			return &value
		}
	}
	return nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	movies = append(movies, Movie{ID: "1", Isbn: "784935", Title: "first movie", Director: &Director{Lastname: "test", Firstname: "name"}})
	movies = append(movies, Movie{ID: "2", Isbn: "7834535", Title: "second movie", Director: &Director{Lastname: "test2", Firstname: "name2"}})

	r.Get("/movies", getMovies)
	r.Get("/movies/{id}", getMovie)
	r.Post("/movies", createMovie)
	r.Post("/movies/{id}", updateMovie)
	r.Delete("/movies/{id}", deleteMovie)

	fmt.Print("Starting server at port 8000\n")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
