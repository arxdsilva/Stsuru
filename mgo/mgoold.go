package main

import (
	"fmt"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	c := session.DB("nintendo").C("pokemons")
	fmt.Println(reflect.TypeOf(c))

	err = c.Insert(
		&Pokemon{"charmander", 10, 100, "fire"},
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
	fmt.Println(reflect.TypeOf(result))
	err = c.Remove(bson.M{}).All(&result)
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
}

type Pokemon struct {
	Name string
	CP   int
	HP   int
	Type string
}
