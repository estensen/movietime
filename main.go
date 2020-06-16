package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli/v2"
)

type Keys struct {
	Omdb string
}

func main() {
	var keys Keys
	err := envconfig.Process("movietime", &keys)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(keys.Omdb)

	app := &cli.App{
		Name: "reviews",
		Usage: "get aggregated reviews of a movie",
		Action: func(c *cli.Context) error {
			movie := c.Args().Get(0)
			fmt.Println(movie)
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://www.omdbapi.com/?t=batman+begins&plot=full&apikey=%s&", keys.Omdb)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
