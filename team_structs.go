package main

import (
	"time"
)

type Standings struct {
	WildCardIndicator    bool      `json:"wildCardIndicator"`
	StandingsDateTimeUtc time.Time `json:"standingsDateTimeUtc"`
	Standings            []struct {
		ConferenceAbbrev         string  `json:"conferenceAbbrev"`
		ConferenceHomeSequence   int     `json:"conferenceHomeSequence"`
		ConferenceL10Sequence    int     `json:"conferenceL10Sequence"`
		ConferenceName           string  `json:"conferenceName"`
		ConferenceRoadSequence   int     `json:"conferenceRoadSequence"`
		ConferenceSequence       int     `json:"conferenceSequence"`
		Date                     string  `json:"date"`
		DivisionAbbrev           string  `json:"divisionAbbrev"`
		DivisionHomeSequence     int     `json:"divisionHomeSequence"`
		DivisionL10Sequence      int     `json:"divisionL10Sequence"`
		DivisionName             string  `json:"divisionName"`
		DivisionRoadSequence     int     `json:"divisionRoadSequence"`
		DivisionSequence         int     `json:"divisionSequence"`
		GameTypeID               int     `json:"gameTypeId"`
		GamesPlayed              int     `json:"gamesPlayed"`
		GoalDifferential         int     `json:"goalDifferential"`
		GoalDifferentialPctg     float64 `json:"goalDifferentialPctg"`
		GoalAgainst              int     `json:"goalAgainst"`
		GoalFor                  int     `json:"goalFor"`
		GoalsForPctg             float64 `json:"goalsForPctg"`
		HomeGamesPlayed          int     `json:"homeGamesPlayed"`
		HomeGoalDifferential     int     `json:"homeGoalDifferential"`
		HomeGoalsAgainst         int     `json:"homeGoalsAgainst"`
		HomeGoalsFor             int     `json:"homeGoalsFor"`
		HomeLosses               int     `json:"homeLosses"`
		HomeOtLosses             int     `json:"homeOtLosses"`
		HomePoints               int     `json:"homePoints"`
		HomeRegulationPlusOtWins int     `json:"homeRegulationPlusOtWins"`
		HomeRegulationWins       int     `json:"homeRegulationWins"`
		HomeTies                 int     `json:"homeTies"`
		HomeWins                 int     `json:"homeWins"`
		L10GamesPlayed           int     `json:"l10GamesPlayed"`
		L10GoalDifferential      int     `json:"l10GoalDifferential"`
		L10GoalsAgainst          int     `json:"l10GoalsAgainst"`
		L10GoalsFor              int     `json:"l10GoalsFor"`
		L10Losses                int     `json:"l10Losses"`
		L10OtLosses              int     `json:"l10OtLosses"`
		L10Points                int     `json:"l10Points"`
		L10RegulationPlusOtWins  int     `json:"l10RegulationPlusOtWins"`
		L10RegulationWins        int     `json:"l10RegulationWins"`
		L10Ties                  int     `json:"l10Ties"`
		L10Wins                  int     `json:"l10Wins"`
		LeagueHomeSequence       int     `json:"leagueHomeSequence"`
		LeagueL10Sequence        int     `json:"leagueL10Sequence"`
		LeagueRoadSequence       int     `json:"leagueRoadSequence"`
		LeagueSequence           int     `json:"leagueSequence"`
		Losses                   int     `json:"losses"`
		OtLosses                 int     `json:"otLosses"`
		PlaceName                struct {
			Default string `json:"default"`
		} `json:"placeName,omitempty"`
		PointPctg                float64 `json:"pointPctg"`
		Points                   int     `json:"points"`
		RegulationPlusOtWinPctg  float64 `json:"regulationPlusOtWinPctg"`
		RegulationPlusOtWins     int     `json:"regulationPlusOtWins"`
		RegulationWinPctg        float64 `json:"regulationWinPctg"`
		RegulationWins           int     `json:"regulationWins"`
		RoadGamesPlayed          int     `json:"roadGamesPlayed"`
		RoadGoalDifferential     int     `json:"roadGoalDifferential"`
		RoadGoalsAgainst         int     `json:"roadGoalsAgainst"`
		RoadGoalsFor             int     `json:"roadGoalsFor"`
		RoadLosses               int     `json:"roadLosses"`
		RoadOtLosses             int     `json:"roadOtLosses"`
		RoadPoints               int     `json:"roadPoints"`
		RoadRegulationPlusOtWins int     `json:"roadRegulationPlusOtWins"`
		RoadRegulationWins       int     `json:"roadRegulationWins"`
		RoadTies                 int     `json:"roadTies"`
		RoadWins                 int     `json:"roadWins"`
		SeasonID                 int     `json:"seasonId"`
		ShootoutLosses           int     `json:"shootoutLosses"`
		ShootoutWins             int     `json:"shootoutWins"`
		StreakCode               string  `json:"streakCode"`
		StreakCount              int     `json:"streakCount"`
		TeamName                 struct {
			Default string `json:"default"`
			Fr      string `json:"fr"`
		} `json:"teamName"`
		TeamCommonName struct {
			Default string `json:"default"`
		} `json:"teamCommonName,omitempty"`
		TeamAbbrev struct {
			Default string `json:"default"`
		} `json:"teamAbbrev"`
		TeamLogo         string  `json:"teamLogo"`
		Ties             int     `json:"ties"`
		WaiversSequence  int     `json:"waiversSequence"`
		WildcardSequence int     `json:"wildcardSequence"`
		WinPctg          float64 `json:"winPctg"`
		Wins             int     `json:"wins"`
	} `json:"standings"`
}
