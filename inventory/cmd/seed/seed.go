package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BoRuDar/configuration"
	"github.com/felipesere/inventory/v0/pkg/products"
	"io/ioutil"
	"net/http"
)

func main() {

	cfg := struct {
		Host string `flag:"host" default:"localhost"`
		Port int64 	`flag:"port" default:"8080"`
		File string `flag:"file"`
	}{}

	conf, err := configuration.New(&cfg, []configuration.Provider{
		configuration.NewFlagProvider(&cfg),
		configuration.NewDefaultProvider(),
	},
		false,
		true,
	)
	if err != nil {
		panic(err.Error())
	}
	conf.InitValues()

	fmt.Println("Reading seeds...")
	content, err := ioutil.ReadFile(cfg.File)
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
		_, err := http.Post(fmt.Sprintf("http://%s:%d/inventory/available", cfg.Host, cfg.Port), "application/json", bytes.NewReader(itemAsJSON))
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println("Done seeding.")
}
