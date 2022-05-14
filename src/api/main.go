package main

//each and every file in go must have a package name
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Pokemon struct {
	Id   int
	Name string
}

var Pokemons []Pokemon

func main() {

	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/api/pokemons", GetAllPokemons)
	myRouter.HandleFunc("/api/pokemons/{id}", GetPokemonById)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", myRouter))
}

func createPokemonsList(data [][]string) []Pokemon {

	// convert records to array of structs
	pokemonsList := []Pokemon{}

	str := ""

	n, _ := strconv.Atoi(str)

	for i, line := range data {
		if i > 0 { // omit header line
			var rec Pokemon
			for j, field := range line {
				if j == 0 {
					n, _ = strconv.Atoi(field)
					rec.Id = n
				} else if j == 1 {
					rec.Name = field
				}
			}
			pokemonsList = append(pokemonsList, rec)
		}
	}
	return pokemonsList
}

func createPokemon(data [][]string, id int) Pokemon {
	pokemonList := createPokemonsList(data)

	pokemon := Pokemon{}

	for _, p := range pokemonList {
		if p.Id == id {
			pokemon = p
		}
	}

	fmt.Println(pokemon)

	return pokemon
}

func readCsv(csvFile string) [][]string {

	// open file
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func GetAllPokemons(w http.ResponseWriter, r *http.Request) {
	pokemonsList := createPokemonsList(readCsv("pokemons.csv"))
	json.NewEncoder(w).Encode(pokemonsList)
}

func GetPokemonById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	n, _ := strconv.Atoi(key)

	pokemon := createPokemon(readCsv("pokemons.csv"), n)
	json.NewEncoder(w).Encode(pokemon)
}
