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

func (cfg *apiConfig) buildStandings(ft database.FollowedTeam) ([]Team, error) {
	standings := Standings{}
	URL := "https://api-web.nhle.com/v1/standings/now"

	// Check cache for data
	entry, ok := cfg.cache.Get(URL)
	if ok {
		if err := json.Unmarshal(entry, &standings); err != nil {
			return nil, err
		}
	} else {
		// Make HTTP request
		resp, err := http.Get(URL)
		if err != nil {
			return nil, err
		}
		// Get byte data
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// Cache data
		cfg.cache.Add(URL, data)

		// Unmarshal data
		if err := json.Unmarshal(data, &standings); err != nil {
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

func (cfg *apiConfig) buildPlayerHelper(i int, followedPlayer database.FollowedPlayer, output []Player, wait *sync.WaitGroup) {
	// Defer counter decrease
	defer wait.Done()

	// Build stats for player
	player, err := cfg.buildPlayerStats(followedPlayer)
	if err != nil {
		return // ERROR HANDLING??
	}
	// Add player to output list at index
	output[i] = player
}
func (cfg *apiConfig) buildPlayerlistWithStats(playerList []database.FollowedPlayer) []Player {
	// Create output list
	output := make([]Player, len(playerList))

	var wait sync.WaitGroup

	// Loop through player list
	for i, followedPlayer := range playerList {
		// Increase counter
		wait.Add(1)

		// Concurrently build list of player stats
		go cfg.buildPlayerHelper(i, followedPlayer, output, &wait)
	}
	// Wait for counter to zero
	wait.Wait()

	return output
}

func (cfg *apiConfig) buildPlayerStats(followedPlayer database.FollowedPlayer) (Player, error) {
	stats := PlayerStats{}
	// Assemble URL
	URL := fmt.Sprintf("https://api-web.nhle.com/v1/player/%v/landing", followedPlayer.PlayerID)

	// Check cache for data
	entry, ok := cfg.cache.Get(URL)
	if ok {
		if err := json.Unmarshal(entry, &stats); err != nil {
			return Player{}, err
		}
	} else {
		// Make HTTP request
		resp, err := http.Get(URL)
		if err != nil {
			return Player{}, err
		}
		defer resp.Body.Close()

		// Get byte data
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return Player{}, err
		}
		// Cache data
		cfg.cache.Add(URL, data)

		// Unmarshal data
		if err := json.Unmarshal(data, &stats); err != nil {
			return Player{}, err
		}
	}
	// Get stat line for last 5 games
	last5StatLine := buildLast5StatLine(stats)

	// Calculate points per game
	ppgFloat := float64(stats.FeaturedStats.RegularSeason.SubSeason.Points) / float64(stats.FeaturedStats.RegularSeason.SubSeason.GamesPlayed)
	ppg := fmt.Sprintf("%.2f", ppgFloat)

	// Get playing today status
	pt, err := playingToday(stats.CurrentTeamAbbrev)
	if err != nil {
		return Player{}, err
	}

	// Make player struct with stats
	player := Player{
		Name:            stats.FirstName.Default + " " + stats.LastName.Default,
		GamesPlayed:     stats.FeaturedStats.RegularSeason.SubSeason.GamesPlayed,
		Goals:           stats.FeaturedStats.RegularSeason.SubSeason.Goals,
		Assists:         stats.FeaturedStats.RegularSeason.SubSeason.Assists,
		Points:          stats.FeaturedStats.RegularSeason.SubSeason.Points,
		PointPercentage: ppg,
		Last5Games:      last5StatLine,
		PlayingToday:    pt,
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

func playingToday(teamAbbrev string) (bool, error) {
	// Assemble URL
	url := fmt.Sprintf("https://api-web.nhle.com/v1/club-schedule/%v/week/now", teamAbbrev)

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Decode JSON
	schedule := Schedule{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&schedule); err != nil {
		return false, err
	}
	// Loop through schedule
	gameCheck := false
	for _, game := range schedule.Games {
		// Find game date
		gameDate, err := time.Parse("2006-01-02", game.GameDate)
		if err != nil {
			return false, nil
		}
		// Compare today to game date
		nowYear, nowMonth, nowDay := time.Now().UTC().Date()
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
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Marshal data to JSON
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Marshalling error", err)
		return
	}
	// Write response
	w.WriteHeader(code)
	w.Write(dat)
}
