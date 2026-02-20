package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dylandeyotte/nhl/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type TeamInfo struct {
	Data []struct {
		ID          int    `json:"id"`
		FranchiseID int    `json:"franchiseId"`
		FullName    string `json:"fullName"`
		LeagueID    int    `json:"leagueId"`
		RawTricode  string `json:"rawTricode"`
		TriCode     string `json:"triCode"`
	} `json:"data"`
	Total int `json:"total"`
}

func main() {
	// Load data from ENV
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("No DB_URL set")
	}
	// Create db handle
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	// Make HTTP request
	url := "https://api.nhle.com/stats/rest/en/team"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	// Decode JSON
	teams := TeamInfo{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&teams); err != nil {
		log.Fatal(err)
	}
	// Load teams into database
	for _, team := range teams.Data {
		_, err := dbQueries.CreateTeam(context.Background(), database.CreateTeamParams{
			ID:       int32(team.ID),
			TeamName: team.FullName,
			TriCode:  team.TriCode,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Teams added to database")
}
