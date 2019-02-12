package main

import (
	"fmt"
	"os"

	"github.com/FuzzyStatic/blizzard"
	"github.com/ccod/go-bnet"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func loadConfig() Config {
	var c Config
	err := godotenv.Load()
	if err != nil {
		fmt.Println("something went wrong")
		return c
	}

	c.domain = "http://localhost"
	if fromEnv := os.Getenv("DOMAIN"); fromEnv != "" {
		c.domain = fromEnv
	}

	c.port = "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		c.port = fromEnv
	}

	c.blizzClientID = os.Getenv("BLIZZ_CLIENT_ID")
	c.blizzClientSecret = os.Getenv("BLIZZ_CLIENT_SECRET")
	c.oauthSalt = os.Getenv("OAUTH_SALT")
	c.jwtSecret = os.Getenv("JWT_SECRET")
	// not in use just yet. going to be where I redirect to after creating the jwt token
	c.clientDomain = os.Getenv("CLIENT_DOMAIN")

	return c
}

func (c *Config) serverSetup(db *gorm.DB) Server {
	blizz := blizzard.NewClient(c.blizzClientID, c.blizzClientSecret, blizzard.US, blizzard.Locale("enUS"))
	err := blizz.AccessTokenReq()
	if err != nil {
		fmt.Println(err)
	}

	return Server{
		&oauth2.Config{
			ClientID:     c.blizzClientID,
			ClientSecret: c.blizzClientSecret,
			Scopes:       []string{"sc2.profile"},
			RedirectURL:  "http://localhost:8080/bnet_oauth_cb",
			Endpoint:     bnet.Endpoint("us"),
		},
		c.oauthSalt,
		blizz,
		*c,
		db,
	}
}
