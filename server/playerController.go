package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func (s *Server) handlerSC2Player(w http.ResponseWriter, r *http.Request) {
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
