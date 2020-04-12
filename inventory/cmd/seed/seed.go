package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/felipesere/inventory/v0/pkg/products"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Reading seeds...")
	content, err := ioutil.ReadFile("./seed.json")
	if err != nil {
		panic(err.Error())
	}

	var seeds struct {
		Items []products.Product `json:"items"`
	}

	err = json.Unmarshal(content, &seeds)

	if err != nil {
		panic(err.Error())
	}

	for _, item := range seeds.Items {
		fmt.Println(fmt.Sprintf("Seeding item %s", item.Name))
		itemAsJSON, _ := json.Marshal(item)
		_, err := http.Post("http://localhost:8080/inventory/available", "application/json", bytes.NewReader(itemAsJSON))
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println("Done seeding.")
}
