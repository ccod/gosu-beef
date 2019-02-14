package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
)

func (s *Server) handlerPlayer(w http.ResponseWriter, r *http.Request) {
	accountID := context.Get(r, JAccID).(int)
	fmt.Printf("AccountID in PlayerHandler: %v", accountID)
	w.Header().Set("Content-type", "application/json")

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

func (s *Server) handlerPlayers(w http.ResponseWriter, r *http.Request) {
	var players []Player
	s.db.Find(&players)

	// This is a temp measure for testing usability of the client
	if len(players) == 1 {
		temp := []string{
			"scrubbles", "Lost Jinjo", "Dan", "Jamary", "Kirby", "Pastry", "Guildlin",
			"Thor", "Mac", "Silverknight", "Hypno", "Nonickname", "Water",
		}

		for i := 0; i < len(temp); i++ {
			player := players[0]
			player.AccountID = i + 1
			player.DisplayName = temp[i]
			s.db.Create(&player)
			players = append(players, player)
		}
	}

	respondJSON(w, r, players)
}
