package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/dylandeyotte/nhl/internal/cache"
	"github.com/dylandeyotte/nhl/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	database *database.Queries
	secret   string
	platform string
	cache    cache.Cache
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type Player struct {
	PlayerID          int    `json:"player_id"`
	Name              string `json:"name"`
	GamesPlayed       int    `json:"games_played"`
	Goals             int    `json:"goals"`
	Assists           int    `json:"assists"`
	Points            int    `json:"points"`
	PointPercentage   string `json:"p/gp"`
	Last5Games        string `json:"last_5_games_totals"`
	PlayingToday      bool   `json:"playing_today"`
	SweaterNumber     int    `json:"sweater_number"`
	Position          string `json:"position"`
	CurrentTeamAbbrev string `json:"team_abbrev"`
	BirthDate         string `json:"birth_date"`
	BirthCity         string `json:"birth_city"`
	BirthCountry      string `json:"birth_country"`
	DraftYear         int    `json:"draft_year"`
	DraftPosition     int    `json:"draft_pos"`
	Height            int    `json:"height"`
	Weight            int    `json:"weight"`
}

type Team struct {
	Name                   string `json:"name"`
	GamesPlayed            int    `json:"games_played"`
	Wins                   int    `json:"wins"`
	Losses                 int    `json:"losses"`
	OvertimeLosses         int    `json:"otl"`
	Points                 int    `json:"points"`
	RegulationWins         int    `json:"rw"`
	RegulationOvertimeWins int    `json:"row"`
	GoalDifferential       int    `json:"goal_differential"`
	Last10                 string `json:"last_10"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load data from ENV
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("No DB_URL set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("No secret set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("No platform set")
	}
	// Create db handle
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	// Create cache
	cache := cache.NewCache(60 * time.Minute)

	// Set up config
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		database: dbQueries,
		secret:   secret,
		platform: platform,
		cache:    cache,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	mux.HandleFunc("DELETE /api/teams/{tricode}/follow", apiCfg.handlerUnfollowTeam)
	mux.HandleFunc("POST /api/teams/{tricode}/follow", apiCfg.handlerFollowTeam)

	mux.HandleFunc("POST /api/players/{playerID}/follow", apiCfg.handlerFollowPlayer)
	mux.HandleFunc("DELETE /api/players/{playerID}/follow", apiCfg.handlerUnfollowPlayer)
	mux.HandleFunc("GET /api/players/search", apiCfg.handlerSearchPlayers)

	mux.HandleFunc("GET /api/following", apiCfg.handlerGetFollows)
	mux.HandleFunc("GET /api/home", apiCfg.handlerHome)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := http.Server{
		Handler: corsMiddleware(mux),
		Addr:    ":8080",
	}
	defer server.Close()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
