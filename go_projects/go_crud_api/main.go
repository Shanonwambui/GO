package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" //Used to create a new router and handle the HTTP requests.
)

type Movie struct {
	ID       string    `json: "id"`
	Isbn     string    `json: "isbn"`
	Title    string    `json: "title"`
	Director *Director `json: "director"`
}
type Director struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) //encodes a Go data structure into JSON format and writes it to an HTTP response.
			return
		}
	}

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	/*Converting JSON data into a Go struct or variable.
	In this case, the JSON data from the HTTP request body is being converted into the movie struct.
	The _ is used to ignore the error returned by the Decode method. In Go, the Decode method returns an error if the JSON decoding fails.
	*/
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	w.Header().Set("Content-Type", "application/json")
	//get the params
	params := mux.Vars(r)
	//loop through the movies
	for index, item := range movies {
		//check if the id matches the id in the params
		if item.ID == params["id"] {
			//delete the movie
			movies = append(movies[:index], movies[index+1:]...)
			//create a new movie
			var movie Movie
			//decode the json data and store it in the movie struct
			_ = json.NewDecoder(r.Body).Decode(&movie)
			//set the id to the id in the params
			movie.ID = params["id"]
			//append the movie to the movies slice
			movies = append(movies, movie)
			//encode the movie struct into json and write it to the response
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "448743", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "445555", Title: "Movie Two", Director: &Director{Firstname: "Shanon", Lastname: "Wambui"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
