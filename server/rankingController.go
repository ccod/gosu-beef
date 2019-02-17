package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) listRanking(w http.ResponseWriter, r *http.Request) {
	var rankings []Ranking
	s.db.Find(&rankings)

	for i := 0; i < len(rankings); i++ {
		var player Player
		s.db.Model(rankings[i]).Related(&player)
		rankings[i].Player = player
	}

	respondJSON(w, r, rankings)
}

func (s *Server) getRanking(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("failed to convert id to int: %s", err)
		return
	}
	var ranking Ranking
	var player Player

	s.db.First(&ranking, id)
	s.db.Model(&ranking).Related(&player)
	ranking.Player = player

	respondJSON(w, r, ranking)
}

// createRanking is for strictly for new or replacing ranked member
func (s *Server) createRanking(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newRanking, prevRanking Ranking
	err := decoder.Decode(&newRanking)
	if err != nil {
		fmt.Printf("failed to decode body: %s", err)
		return
	}

	// avoid having two players with the same rank
	s.db.Where("rank = ?", newRanking.Rank).First(&prevRanking)
	if prevRanking.ID != 0 {
		s.db.Delete(&prevRanking)
	}

	s.db.Create(&newRanking)

	respondJSON(w, r, newRanking)
}

// promoteRanking is for normal challenge promotion, takes rank of the challenged, and rotates players below down a rank
// TODO: still need to remove previous player rank if exists
func (s *Server) promoteRanking(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newRanking Ranking
	var rankings []Ranking
	err := decoder.Decode(&newRanking)
	if err != nil {
		fmt.Printf("failed to decode body: %s", err)
		return
	}

	s.db.Where("rank >= ?", newRanking.Rank).Find(&rankings)
	s.db.Where("rank >= ?", newRanking.Rank).Delete(Ranking{})
	for i := 0; i < len(rankings); i++ {
		var player Player
		if rankings[i].PlayerID == newRanking.PlayerID {
			continue
		}

		nextRanking := Ranking{
			Rank:     rankings[i].Rank + 1,
			PlayerID: rankings[i].PlayerID,
		}

		s.db.Create(&nextRanking)
		s.db.Model(&nextRanking).Related(&player)
		nextRanking.Player = player
		rankings[i] = nextRanking
	}

	s.db.Create(&newRanking)

	rankings = append(rankings, newRanking)
	respondJSON(w, r, rankings)
}

func (s *Server) deleteRanking(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("failed to convert id to int: %s", err)
		return
	}

	var ranking Ranking
	s.db.First(&ranking, id)
	s.db.Delete(&ranking)

	respondJSON(w, r, ranking)
}
