package main

import (
	"encoding/json"
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

type Movie struct {
	Title string
	Year string
	Plot string
}

var keys Keys

func main() {
	err := envconfig.Process("movietime", &keys)
	if err != nil {
		log.Fatal(err.Error())
	}

	app := &cli.App{
		Name: "reviews",
		Usage: "get aggregated reviews of a movie",
		Action: func(c *cli.Context) error {
			searchTitle := c.Args().Get(0)
			getMovie(searchTitle)
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getMovie(searchTitle string) {
	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&plot=full&apikey=%s&", searchTitle, keys.Omdb)
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

	movie := Movie{}
	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(movie.Title)
}
