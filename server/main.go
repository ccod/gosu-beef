package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ccod/blizzard"

	bnet "github.com/ccod/go-bnet"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
)

// Player is going to be the main identifier for this service
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

type config struct {
	domain            string
	port              string
	blizzClientID     string
	blizzClientSecret string
	oauthSalt         string
	jwtSecret         string
	clientDomain      string
}

type server struct {
	oauthCfg         *oauth2.Config
	oauthStateString string
	blizz            *blizzard.Client
	env              config
	db               *gorm.DB
}

// Player is the return struct for collecting profile information
// type Player struct {
// 	Profile sc2c.Player  `json:"profile"`
// 	Summary sc2c.Profile `json:"ladders"`
// }

func (c *config) load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("something went wrong")
		return
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
}

const htmlIndex = `<html><body>
Log in with <a href="/login">Bnet</a>
</body></html>
`

func (s *server) handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func (s *server) handlBnetLogin(w http.ResponseWriter, r *http.Request) {
	url := s.oauthCfg.AuthCodeURL(s.oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// TODO: first print a jwt with user info, later make it a redirect to webapp with jwt in the url or something
func (s *server) handleBnetCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != s.oauthStateString {
		fmt.Printf("invalid oauth state expected '%s', go '%s'\n", s.oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := s.oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := s.oauthCfg.Client(oauth2.NoContext, token)
	client := bnet.NewClient("us", oauthClient)
	user, _, err := client.UserInfo()
	if err != nil {
		fmt.Printf("client.Profile().SC2() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	claims := &jwt.StandardClaims{
		Id:     strconv.Itoa(user.ID),
		Issuer: "gosu-beef",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := jwtToken.SignedString([]byte(s.env.jwtSecret))
	if err != nil {
		fmt.Printf("jwt signing failed: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("UserInfo is: %v\n", user)
	fmt.Printf("Jwt token: %s", tokenString)
	http.Redirect(w, r, s.env.clientDomain+"/callback#"+tokenString, http.StatusTemporaryRedirect)
}

// check claim from jwt, check db for row, if exists return profile data
// if empty, pull data from blizzard api,
func (s *server) handlerSC2Player(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	reqToken = strings.Split(reqToken, "Bearer ")[1]

	w.Header().Set("Content-type", "application/json")

	token, err := jwt.ParseWithClaims(reqToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.env.jwtSecret), nil
	})

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !(ok && token.Valid) {
		fmt.Printf("error with Parse with Claims: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	accountID, err := strconv.Atoi(claims.Id)
	if err != nil {
		fmt.Printf("Atoi call failed: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	var player Player
	s.db.First(&player, accountID)
	if player.AccountID != 0 {
		fmt.Println("I succeeded in finding the saved player")

		var ladderRecord LadderRecord
		s.db.Model(&player).Related(&ladderRecord)
		player.LadderRecord = ladderRecord

		response, err := json.Marshal(player)
		if err != nil {
			fmt.Printf("JSON encoding failed: %s", err)
			w.Write([]byte("{\"failure\":true}"))
			return
		}

		// fmt.Printf("response: %s", response)
		w.Write(response)
		return
	}

	// Player doesn't exist in DB, so preparing to call API to populate it
	s.blizz.TokenValidation()

	playerData, _, err := s.blizz.SC2Player(accountID)
	if err != nil {
		fmt.Printf("SC2Player call failed: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	profileID, err := strconv.Atoi((*playerData)[0].ProfileID)
	if err != nil {
		fmt.Printf("profileID conversion call failed: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	sc2Profile, _, err := s.blizz.SC2Profile(1, 1, profileID)
	if err != nil {
		fmt.Printf("SC2Profile call failed: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	sc2LegacyProfile, _, err := s.blizz.SC2LegacyProfile(1, 1, profileID)
	if err != nil {
		fmt.Printf("SC2LegacyProfile call failed: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	player = Player{
		AccountID:   accountID,
		ProfileID:   profileID,
		ProfileURL:  (*playerData)[0].ProfileURL,
		AvatarURL:   (*playerData)[0].AvatarURL,
		ClanTag:     sc2LegacyProfile.ClanTag,
		ClanName:    sc2LegacyProfile.ClanName,
		DisplayName: (*playerData)[0].Name,
		RealmID:     (*playerData)[0].RealmID,
		RegionID:    (*playerData)[0].RegionID,
		LadderRecord: LadderRecord{
			Rank:          sc2Profile.Snapshot.SeasonSnapshot.OneV1.Rank,
			League:        sc2Profile.Snapshot.SeasonSnapshot.OneV1.LeagueName.(string),
			Wins:          sc2Profile.Snapshot.SeasonSnapshot.OneV1.TotalWins,
			Total:         sc2Profile.Snapshot.SeasonSnapshot.OneV1.TotalGames,
			PreferredRace: sc2LegacyProfile.Career.PrimaryRace,
		},
	}

	s.db.Create(&player)

	response, err := json.Marshal(player)
	if err != nil {
		fmt.Printf("JSON encoding failed: %s", err)
		w.Write([]byte("{\"failure\":true}"))
		return
	}

	w.Write(response)
}

func main() {
	c := config{}
	c.load()

	db, _ := gorm.Open("sqlite3", "./gorm.db")
	defer db.Close()

	db.AutoMigrate(&Player{})
	db.AutoMigrate(&LadderRecord{})

	blizz := blizzard.NewClient(c.blizzClientID, c.blizzClientSecret, blizzard.US, blizzard.Locale("enUS"))
	err := blizz.AccessTokenReq()
	if err != nil {
		fmt.Println(err)
	}

	s := server{
		&oauth2.Config{
			ClientID:     c.blizzClientID,
			ClientSecret: c.blizzClientSecret,
			Scopes:       []string{"sc2.profile"},
			RedirectURL:  "http://localhost:8080/bnet_oauth_cb",
			Endpoint:     bnet.Endpoint("us"),
		},
		c.oauthSalt,
		blizz,
		c,
		db,
	}

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,

		Debug: true,
	})

	r := mux.NewRouter()

	r.HandleFunc("/", s.handleMain)
	r.HandleFunc("/login", s.handlBnetLogin)
	r.HandleFunc("/bnet_oauth_cb", s.handleBnetCallback)
	r.HandleFunc("/accountID", s.handlerSC2Player)

	handler := cor.Handler(r)
	fmt.Print("Started running on " + c.domain + ":" + c.port + "\n")
	fmt.Println(http.ListenAndServe(":"+c.port, handler))
}
