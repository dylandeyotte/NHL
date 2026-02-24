package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dylandeyotte/nhl/internal/auth"
	"github.com/dylandeyotte/nhl/internal/database"
)

func (cfg *apiConfig) handlerHome(w http.ResponseWriter, r *http.Request) {
	// Get bearer token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve token", err)
		return
	}
	// Validate JWT
	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve token", err)
		return
	}
	// Get list of followed players
	playerList, err := cfg.database.GetFollowedPlayers(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "Error finding followed players", err)
			return
		}
	}
	// Build list of players with stats
	returnPlayerList, err := cfg.buildPlayerlistWithStats(playerList)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating player stats", err)
		return
	}

	// Get followed team
	returnTeamList := []Team{}
	team, err := cfg.database.GetFollowedTeam(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "Error finding followed team", err)
			return
		}
	} else {
		// Build standings if a team is followed
		returnTeamList, err = cfg.buildStandings(team)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error building standings", err)
			return
		}
	}
	// Build JSON payload for return
	type JSONData struct {
		Standings []Team   `json:"standings"`
		Players   []Player `json:"players"`
	}
	payload := JSONData{
		Standings: returnTeamList,
		Players:   returnPlayerList,
	}

	// Respond with JSON
	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerGetFollows(w http.ResponseWriter, r *http.Request) {
	// Get token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retreive token", err)
		return
	}
	// Validate token
	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}
	// Get list of followed players
	playerList, err := cfg.database.GetFollowedPlayers(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "Error finding followed team", err)
			return
		}
	}
	// Build player list of names
	returnPlayerList := []string{}
	for _, player := range playerList {
		returnPlayerList = append(returnPlayerList, player.PlayerName)
	}
	// Get followed team
	followedTeam, err := cfg.database.GetFollowedTeam(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "Error finding followed team", err)
			return
		}
	}
	// Build JSON payload for return
	type JSONData struct {
		Team    string   `json:"team"`
		Players []string `json:"players"`
	}
	payload := JSONData{
		Team:    followedTeam.TeamName,
		Players: returnPlayerList,
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerSearchPlayers(w http.ResponseWriter, r *http.Request) {
	// Get name and limit from query
	name := r.URL.Query().Get("player")
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	// Assemble URL
	URL := "https://search.d3.nhle.com/api/v1/search/player?culture=en-us&"
	params := url.Values{}
	params.Add("active", "true")
	params.Add("q", name)
	params.Add("limit", limit)

	fullURL := URL + params.Encode()

	// Make HTTP request
	resp, err := http.Get(fullURL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error making HTTP request", err)
		return
	}
	defer resp.Body.Close()

	// Decode JSON
	psr := PlayerSearchResults{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&psr); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, psr)
}

func (cfg *apiConfig) handlerUnfollowPlayer(w http.ResponseWriter, r *http.Request) {
	// Retrieve and parse ID
	playerID := r.PathValue("playerID")

	// Get token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retreive token", err)
		return
	}
	// Validate token
	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}
	// Convert ID from string to int
	playerIDInt, err := strconv.Atoi(playerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error converting string to int", err)
		return
	}
	// Unfollow player
	if err := cfg.database.UnfollowPlayer(r.Context(), database.UnfollowPlayerParams{
		PlayerID: int32(playerIDInt),
		UserID:   userID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error unfollowing player", err)
		return
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, "Unfollow successful")
}

func (cfg *apiConfig) handlerFollowPlayer(w http.ResponseWriter, r *http.Request) {
	// Retrieve and parse ID
	playerID := r.PathValue("playerID")

	// Get token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retreive token", err)
		return
	}
	// Validate token
	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}
	// Assemble URL
	URL := "https://api-web.nhle.com/v1/player"
	fullURL := fmt.Sprintf("%v/%v/landing", URL, playerID)

	// Make HTTP request for player info
	resp, err := http.Get(fullURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid player ID", err)
		return
	}
	defer resp.Body.Close()

	// Decode JSON
	playerInfo := PlayerStats{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&playerInfo); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	// Create player in database
	player, err := cfg.database.CreatePlayer(r.Context(), database.CreatePlayerParams{
		ID:         int32(playerInfo.PlayerID),
		PlayerName: playerInfo.FirstName.Default + " " + playerInfo.LastName.Default,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create player in database", err)
		return
	}
	// Create followed player in database
	followedPlayer, err := cfg.database.FollowPlayer(r.Context(), database.FollowPlayerParams{
		UserID:     userID,
		PlayerID:   player.ID,
		PlayerName: player.PlayerName,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not follow player", err)
		return
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, followedPlayer)
}
