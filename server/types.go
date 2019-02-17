package main

import (
	"github.com/FuzzyStatic/blizzard"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

// JwtKey is used to contain the account Id after verifying the JwtToken
type JwtKey int

// JAccID will be used to store accountID from JWT middleware
const JAccID JwtKey = 0

// Challenge another participant,
// possibly set cron job to check challenge after a week
// if no other challenge has been accepted by the recipient
type Challenge struct {
	gorm.Model
	Issuer    int  `gorm:"foreign_key:player"`
	Recipient int  `gorm:"foreign_key:player"`
	Accepted  bool `gorm:"defaul:false" json:"accepted"`
}

// Ranking the order in which I pupulate the settle the beef pyramid
type Ranking struct {
	gorm.Model
	Rank     int    `json:"rank"`
	Player   Player `gorm:"foreignkey:PlayerID" json:"player"`
	PlayerID int    `json:"playerId"`
}

// History struct{}

// Player is going to be basically my user struct
type Player struct {
	AccountID      int          `gorm:"primary_key" json:"accountId"`
	ProfileID      int          `json:"profileId"`
	ProfileURL     string       `json:"profileUrl"`
	AvatarURL      string       `json:"avatarUrl"`
	ClanTag        string       `json:"clanTag"`
	ClanName       string       `json:"clanName"`
	DisplayName    string       `json:"displayName"`
	RealmID        int          `json:"realmId"`
	RegionID       int          `json:"regionId"`
	Admin          bool         `gorm:"default:false" json:"admin"`
	LadderRecord   LadderRecord `json:"ladderRecord"`
	LadderRecordID int          `json:"-"`
}

// LadderRecord is only referring to 1v1
type LadderRecord struct {
	ID            int    `gorm:"AUTO_INCREMENT" json:"id"`
	Rank          int    `json:"rank"`
	League        string `json:"league"`
	Wins          int    `json:"wins"`
	Total         int    `json:"total"`
	PreferredRace string `json:"preferredRace"`
}

// Config is where I store all the variables that need to be pulled from the environment
type Config struct {
	domain            string
	port              string
	blizzClientID     string
	blizzClientSecret string
	oauthSalt         string
	jwtSecret         string
	clientDomain      string
}

// Server operates as a container for app pieces needed across the app
type Server struct {
	oauthCfg         *oauth2.Config
	oauthStateString string
	blizz            *blizzard.Client
	env              Config
	db               *gorm.DB
}
