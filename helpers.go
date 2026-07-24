package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dylandeyotte/nhl/internal/database"
)

type PlayerCard int

const (
	Stats PlayerCard = iota
	Bio
)

func (cfg *apiConfig) buildStandings(ft database.FollowedTeam) ([]Team, error) {
	standings := Standings{}
	URL := "https://api-web.nhle.com/v1/standings/now"

	// Check cache for data
	entry, ok := cfg.cache.Get(URL)
	if ok {
		if err := json.Unmarshal(entry, &standings); err != nil {
			fmt.Printf("unmarshalling standings error from cache: %v\n", err)
			return nil, err
		}
	} else {
		// Make HTTP request
		resp, err := http.Get(URL)
		if err != nil {
			fmt.Printf("hhtp request error: %v\n", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Get byte data
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("byte read error: %v\n", err)
			return nil, err
		}
		// Cache data
		if err := cfg.cache.Add(URL, data); err != nil {
			return nil, err
		}

		// Unmarshal data
		if err := json.Unmarshal(data, &standings); err != nil {
			fmt.Printf("unmarshalling standings error: %v\n", err)
			return nil, err
		}
	}
	// Get team index and team div
	returnList := []Team{}
	var teamIndex int
	for i := range standings.Standings {
		if standings.Standings[i].TeamName.Default == ft.TeamName {
			teamIndex = i
		}
	}
	teamDiv := standings.Standings[teamIndex].DivisionAbbrev

	// Get standings for team division
	for _, team := range standings.Standings {
		if team.DivisionAbbrev == teamDiv {
			last10Record := fmt.Sprintf("%v-%v-%v", team.L10Wins, team.L10Losses, team.L10OtLosses)
			teamData := Team{
				Name:                   team.TeamName.Default,
				GamesPlayed:            team.GamesPlayed,
				Wins:                   team.Wins,
				Losses:                 team.Losses,
				OvertimeLosses:         team.OtLosses,
				Points:                 team.Points,
				RegulationWins:         team.RegulationWins,
				RegulationOvertimeWins: team.RegulationPlusOtWins,
				GoalDifferential:       team.GoalDifferential,
				Last10:                 last10Record,
			}
			returnList = append(returnList, teamData)
		}
	}
	return returnList, nil
}

func (cfg *apiConfig) buildPlayerHelper(i int, followedPlayer database.FollowedPlayer, output []Player, wait *sync.WaitGroup, errCh chan error, pc PlayerCard) {
	// Defer counter decrease
	defer wait.Done()

	// Build stats for player
	player, err := cfg.buildPlayerInfo(followedPlayer, pc)
	if err != nil {
		fmt.Printf("build player error: %v\n", err)
		errCh <- err
		return
	}
	// Add player to output list at index
	output[i] = player
}
func (cfg *apiConfig) buildPlayerlist(playerList []database.FollowedPlayer, pc PlayerCard) ([]Player, error) {
	// Create output list
	output := make([]Player, len(playerList))

	var wait sync.WaitGroup

	errCh := make(chan error, len(playerList))

	// Loop through player list
	for i, followedPlayer := range playerList {
		// Increase counter
		wait.Add(1)

		// Concurrently build list of player stats
		go cfg.buildPlayerHelper(i, followedPlayer, output, &wait, errCh, pc)
	}
	// Wait for counter to zero and close channel
	wait.Wait()
	close(errCh)

	// Loop through channel, checking for err
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func (cfg *apiConfig) buildPlayerInfo(followedPlayer database.FollowedPlayer, pc PlayerCard) (Player, error) {
	stats := PlayerStats{}
	// Assemble URL
	URL := fmt.Sprintf("https://api-web.nhle.com/v1/player/%v/landing", followedPlayer.PlayerID)

	// Check cache for data
	entry, ok := cfg.cache.Get(URL)
	if ok {
		if err := json.Unmarshal(entry, &stats); err != nil {
			fmt.Printf("unmarshal from cache error: %v\n", err)
			return Player{}, err
		}
	} else {
		// Make HTTP request
		resp, err := http.Get(URL)
		if err != nil {
			fmt.Printf("http request error: %v\n", err)
			return Player{}, err
		}
		defer resp.Body.Close()

		// Get byte data
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("byte read error: %v\n", err)
			return Player{}, err
		}
		// Cache data
		if err := cfg.cache.Add(URL, data); err != nil {
			fmt.Printf("Cannot cache bad data: %v\n", err)
			return Player{}, err
		}

		// Unmarshal data
		if err := json.Unmarshal(data, &stats); err != nil {
			fmt.Printf("unmarshal error: %v\n", err)
			return Player{}, err
		}
	}

	// Return player struct for bio
	if pc == Bio {
		return Player{
			PlayerID:          stats.PlayerID,
			Name:              stats.FirstName.Default + " " + stats.LastName.Default,
			CurrentTeamAbbrev: stats.CurrentTeamAbbrev,
			Position:          stats.Position,
			BirthDate:         stats.BirthDate,
			BirthCity:         stats.BirthCity.Default,
			BirthCountry:      stats.BirthCountry,
			DraftYear:         stats.DraftDetails.Year,
			DraftPosition:     stats.DraftDetails.OverallPick,
			Height:            stats.HeightInInches,
			Weight:            stats.WeightInPounds,
		}, nil
	}

	// Get stat line for last 5 games
	last5StatLine := buildLast5StatLine(stats)

	// Calculate points per game
	ppgFloat := float64(stats.FeaturedStats.RegularSeason.SubSeason.Points) / float64(stats.FeaturedStats.RegularSeason.SubSeason.GamesPlayed)
	ppg := fmt.Sprintf("%.2f", ppgFloat)
	if stats.FeaturedStats.RegularSeason.SubSeason.GamesPlayed == 0 {
		ppg = "-"
	}

	// Get playing today status
	pt, err := cfg.playingToday(stats.CurrentTeamAbbrev, "", nil)
	if err != nil {
		fmt.Printf("playing today error: %v", err)
		return Player{}, err
	}

	// Make player struct with stats
	player := Player{
		Name:              stats.FirstName.Default + " " + stats.LastName.Default,
		GamesPlayed:       stats.FeaturedStats.RegularSeason.SubSeason.GamesPlayed,
		Goals:             stats.FeaturedStats.RegularSeason.SubSeason.Goals,
		Assists:           stats.FeaturedStats.RegularSeason.SubSeason.Assists,
		Points:            stats.FeaturedStats.RegularSeason.SubSeason.Points,
		PointPercentage:   ppg,
		Last5Games:        last5StatLine,
		PlayingToday:      pt,
		SweaterNumber:     stats.SweaterNumber,
		Position:          stats.Position,
		CurrentTeamAbbrev: stats.CurrentTeamAbbrev,
	}
	return player, nil
}

func buildLast5StatLine(stats PlayerStats) string {
	// Set variables
	var goalTotal int
	var assistTotal int
	var pointTotal int

	// Loop over last 5 games
	for _, stat := range stats.Last5Games {
		goalTotal += stat.Goals
		assistTotal += stat.Assists
		pointTotal += stat.Points
	}
	last5StatLine := fmt.Sprintf("%v-%v-%v", goalTotal, assistTotal, pointTotal)
	return last5StatLine
}

func (cfg *apiConfig) playingToday(teamAbbrev, baseURL string, client *http.Client) (bool, error) {
	schedule := Schedule{}
	// Check URL and client for testing
	if baseURL == "" {
		baseURL = "https://api-web.nhle.com"
	}
	if client == nil {
		client = http.DefaultClient
	}
	// Assemble URL
	url := fmt.Sprintf("%v/v1/club-schedule/%v/week/now", baseURL, teamAbbrev)

	// Check cache for data
	entry, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(entry, &schedule); err != nil {
			fmt.Printf("unmarshal from cache error: %v\n", err)
			return false, err
		}
	} else {
		// Make HTTP request
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("http request error: %v\n", err)
			return false, err
		}
		defer resp.Body.Close()

		// Get byte data
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("byte data error for %v: %v\n", teamAbbrev, err)
			return false, err
		}
		// Cache data
		if err := cfg.cache.Add(url, data); err != nil {
			return false, err
		}

		// Unmarshal data
		if err := json.Unmarshal(data, &schedule); err != nil {
			fmt.Printf("Unmarshalling error for %v\n", teamAbbrev)
			return false, err
		}
	}

	// Loop through schedule
	gameCheck := false
	for _, game := range schedule.Games {
		// Find game date
		gameDate, err := time.Parse("2006-01-02", game.GameDate)
		if err != nil {
			fmt.Printf("parsing error: %v\n", err)
			return false, nil
		}
		// Compare today to game date
		nowYear, nowMonth, nowDay := time.Now().Date()
		gameYear, gameMonth, gameDay := gameDate.Date()
		if gameYear == nowYear && gameMonth == nowMonth && gameDay == nowDay {
			gameCheck = true
		}

	}
	return gameCheck, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
		fmt.Println(msg)
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Marshal data to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Marshalling error", err)
		return
	}
	// Write response
	w.WriteHeader(code)
	w.Write(data)
}
