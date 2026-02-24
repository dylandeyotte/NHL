package main

import (
	"net/http"

	"github.com/dylandeyotte/nhl/internal/auth"
	"github.com/dylandeyotte/nhl/internal/database"
)

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
