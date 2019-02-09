package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/FuzzyStatic/blizzard"
	"github.com/joho/godotenv"
)

type config struct {
	port              string
	blizzClientID     string
	blizzClientSecret string
	sampleAccoundID   int
}

func (c *config) load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("something went wrong")
		return
	}

	c.blizzClientID = os.Getenv("BLIZZ_CLIENT_ID")
	c.blizzClientSecret = os.Getenv("BLIZZ_CLIENT_SECRET")
	c.sampleAccoundID, err = strconv.Atoi(os.Getenv("SAMPLE_ACCOUNT_ID"))
	if err != nil {
		fmt.Printf("config convert account id to int failed: %s\n", err)
		return
	}
}

func main() {
	c := config{}
	c.load()

	blizz := blizzard.NewClient(c.blizzClientID, c.blizzClientSecret, blizzard.US, blizzard.Locale("en_US"))
	err := blizz.AccessTokenReq()
	if err != nil {
		fmt.Println(err)
	}

	dat, _, err := blizz.SC2Player(c.sampleAccoundID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", (*dat)[0].ProfileID)
}
