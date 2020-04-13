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
	fail(err)
	conf.InitValues()

	fmt.Println("Reading seeds...")
	content, err := ioutil.ReadFile(cfg.File)
	fail(err)

	var seeds struct {
		Items []products.Product `json:"items"`
	}

	fail(json.Unmarshal(content, &seeds))

	for _, item := range seeds.Items {
		fmt.Println(fmt.Sprintf("Seeding item %s", item.Name))
		itemAsJSON, _ := json.Marshal(item)
		_, err := http.Post(fmt.Sprintf("http://%s:%d/inventory/available", cfg.Host, cfg.Port), "application/json", bytes.NewReader(itemAsJSON))
		fail(err)
	}

	fmt.Println("Done seeding.")
}

func fail(err error) {
	if err != nil {
		panic(err.Error())
	}
}
