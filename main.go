package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies { //Como forEach do JavaScript
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = returnLastIdFromMovies()
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {

			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)

			break
		}
	}
}

func returnLastIdFromMovies() string {
	sort.SliceStable(movies, func(i, j int) bool {
		return movies[i].ID < movies[j].ID
	})

	lastMovie := movies[len(movies)-1]

	lastId, err := strconv.Atoi(lastMovie.ID)

	if err != nil {
		fmt.Println(err)

	} else {

		lastId = lastId + 1
	}

	return strconv.Itoa(lastId)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "91623", Title: "Titanic", Director: &Director{Firstname: "Mark", Lastname: "Scorcese"}})
	movies = append(movies, Movie{ID: "2", Isbn: "36987", Title: "Top Gun", Director: &Director{Firstname: "Luka", Lastname: "Modrich"}})
	movies = append(movies, Movie{ID: "3", Isbn: "12345", Title: "Don't Look Up", Director: &Director{Firstname: "Leonardo", Lastname: "DiCaprio"}})
	movies = append(movies, Movie{ID: "4", Isbn: "18178", Title: "Blade Trinity", Director: &Director{Firstname: "Wesley", Lastname: "Snipes"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/create-movie", createMovie).Methods("POST")
	r.HandleFunc("/edit-movie/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/delete-movie/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
