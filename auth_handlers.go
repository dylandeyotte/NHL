package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dylandeyotte/nhl/internal/auth"
	"github.com/dylandeyotte/nhl/internal/database"
)

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
