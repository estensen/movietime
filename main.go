package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli/v2"
)

type Keys struct {
	Omdb      string
	Translate string
}

type Rating struct {
	Source string
	Value  string
}

type Movie struct {
	Title   string
	Year    string
	Plot    string
	Ratings []Rating
}

var keys Keys

func main() {
	err := envconfig.Process("movietime", &keys)
	if err != nil {
		log.Fatal(err.Error())
	}

	var lang string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lang",
				Value:       "en",
				Usage:       "Language to translate plot text to",
				Destination: &lang,
			},
		},
		Name:  "reviews",
		Usage: "get aggregated reviews of a movie",
		Action: func(c *cli.Context) error {
			searchTitle := c.Args().Get(0)
			getMovie(searchTitle, lang)
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getMovie(searchTitle, lang string) {
	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&plot=full&apikey=%s&", searchTitle, keys.Omdb)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	movie := Movie{}
	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(movie.Title)

	printScore(movie)

	translatedPlot, err := translateText(lang, movie.Plot)
	if err != nil {
		log.Fatal(err)
	}
	println(translatedPlot)
}

func printScore(movie Movie) {
	for r := range movie.Ratings {
		if movie.Ratings[r].Source == "Rotten Tomatoes" {
			score, _ := strconv.Atoi(strings.Trim(movie.Ratings[r].Value, "%"))
			stars := strings.Repeat("⭐️", score/10)
			fmt.Println(stars)
		}
	}
}

func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey(keys.Translate))
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
