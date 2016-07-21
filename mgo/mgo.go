package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	// fazer a conexao:
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// Sets a DB and collection to work with
	// By default if mongoDB doesnt have the named collection/database,
	// It will create automatically
	c := session.DB("nintendo").C("pokemons")

	// Inerts the named pokemons in the mgoDB by using the Pokemons struct
	err = c.Insert(
		&Pokemon{"charmander", 10, 100, "fire"},
		&Pokemon{"squirtle", 10, 100, "water"},
		&Pokemon{"bulbasaur", 10, 100, "grass"},
		&Pokemon{"pikachu", 10, 100, "electric"},
	)
	if err != nil {
		panic(err)
	}
	result := []Pokemon{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// Pokemons is a nice structure that will be used to make CRUD operations into
// the mongoDB collection
type Pokemon struct {
	Name string
	CP   int
	HP   int
	Type string
}
