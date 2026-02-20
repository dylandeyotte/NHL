package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

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
	returnPlayerList := cfg.buildPlayerlistWithStats(playerList)

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

func (cfg *apiConfig) handlerFollowTeam(w http.ResponseWriter, r *http.Request) {
	// Retrieve and parse ID
	triCode := r.PathValue("tricode")

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
	// Get team from database
	team, err := cfg.database.FetchTeamByTriCode(r.Context(), triCode)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to find team with given tri code", err)
		return
	}
	// Create followed team in database
	followedTeam, err := cfg.database.FollowTeam(r.Context(), database.FollowTeamParams{
		TeamName: team.TeamName,
		UserID:   userID,
		TeamID:   team.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error following team", err)
		return
	}
	// Repsond with JSON
	respondWithJSON(w, http.StatusOK, followedTeam)
}

func (cfg *apiConfig) handlerUnfollowTeam(w http.ResponseWriter, r *http.Request) {
	// Retrieve and parse ID
	triCode := r.PathValue("tricode")

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
	// Get team from database
	team, err := cfg.database.FetchTeamByTriCode(r.Context(), triCode)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to find team with given tri code", err)
		return
	}
	// Unfollow team
	if err := cfg.database.UnfollowTeam(r.Context(), database.UnfollowTeamParams{
		UserID: userID,
		TeamID: team.ID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error unfollowing team", err)
		return
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, "Unfollow successful")
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Create struct for data
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	// Deocde JSON into struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode JSON", err)
		return
	}
	// Fetch user
	user, err := cfg.database.FetchUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to find user", err)
		return
	}
	// Validate password
	match, err := auth.ValidatePassword(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error matching password and hash", err)
		return
	}
	if match != true {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Create JWT
	token, err := auth.MakeJWT(user.ID, cfg.secret, 60*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create JWT", err)
		return
	}
	// Create refresh token
	refreshTokenString := auth.MakeRefreshToken()
	_, err = cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(720 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	// Create user for JSON
	userFinal := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshTokenString,
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, userFinal)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Create struct for data
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	// Deocde JSON into struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode JSON", err)
		return
	}
	// Hash password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}
	// Create user in database
	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}
	// Create user for JSON
	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusOK, newUser)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Get token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retreive token", err)
		return
	}
	// Get user from token
	user, err := cfg.database.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error finding user for given token", err)
		return
	}
	// Make new JWT
	token, err := auth.MakeJWT(user.ID, cfg.secret, 60*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}
	// Respond with JSON
	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retreive token", err)
		return
	}
	// Revoke token
	_, err = cfg.database.RevokeToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error revoking token", err)
		return
	}
	// Respond with JSON
	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Check if dev
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", nil)
		return
	}
	// Delete users
	if err := cfg.database.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		return
	}
	// Delete players
	if err := cfg.database.DeletePlayers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting players", err)
		return
	}
	// Delete player follows
	if err := cfg.database.DeletePlayerFollows(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting follows", err)
		return
	}
	// Delete team follows
	if err := cfg.database.DeleteTeamFollows(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting follows", err)
		return
	}
	// Write confirmation
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users deleted"))
}
