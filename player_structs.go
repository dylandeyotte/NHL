package main

type PlayerStats struct {
	PlayerID          int    `json:"playerId"`
	IsActive          bool   `json:"isActive"`
	CurrentTeamID     int    `json:"currentTeamId"`
	CurrentTeamAbbrev string `json:"currentTeamAbbrev"`
	FullTeamName      struct {
		Default string `json:"default"`
		Fr      string `json:"fr"`
	} `json:"fullTeamName"`
	TeamCommonName struct {
		Default string `json:"default"`
	} `json:"teamCommonName"`
	TeamPlaceNameWithPreposition struct {
		Default string `json:"default"`
		Fr      string `json:"fr"`
	} `json:"teamPlaceNameWithPreposition"`
	FirstName struct {
		Default string `json:"default"`
	} `json:"firstName"`
	LastName struct {
		Default string `json:"default"`
	} `json:"lastName"`
	Badges              []any  `json:"badges"`
	TeamLogo            string `json:"teamLogo"`
	SweaterNumber       int    `json:"sweaterNumber"`
	Position            string `json:"position"`
	Headshot            string `json:"headshot"`
	HeroImage           string `json:"heroImage"`
	HeightInInches      int    `json:"heightInInches"`
	HeightInCentimeters int    `json:"heightInCentimeters"`
	WeightInPounds      int    `json:"weightInPounds"`
	WeightInKilograms   int    `json:"weightInKilograms"`
	BirthDate           string `json:"birthDate"`
	BirthCity           struct {
		Default string `json:"default"`
	} `json:"birthCity"`
	BirthStateProvince struct {
		Default string `json:"default"`
	} `json:"birthStateProvince"`
	BirthCountry  string `json:"birthCountry"`
	ShootsCatches string `json:"shootsCatches"`
	DraftDetails  struct {
		Year        int    `json:"year"`
		TeamAbbrev  string `json:"teamAbbrev"`
		Round       int    `json:"round"`
		PickInRound int    `json:"pickInRound"`
		OverallPick int    `json:"overallPick"`
	} `json:"draftDetails"`
	PlayerSlug      string `json:"playerSlug"`
	InTop100AllTime int    `json:"inTop100AllTime"`
	InHHOF          int    `json:"inHHOF"`
	FeaturedStats   struct {
		Season        int `json:"season"`
		RegularSeason struct {
			SubSeason struct {
				Assists           int     `json:"assists"`
				GameWinningGoals  int     `json:"gameWinningGoals"`
				GamesPlayed       int     `json:"gamesPlayed"`
				Goals             int     `json:"goals"`
				OtGoals           int     `json:"otGoals"`
				Pim               int     `json:"pim"`
				PlusMinus         int     `json:"plusMinus"`
				Points            int     `json:"points"`
				PowerPlayGoals    int     `json:"powerPlayGoals"`
				PowerPlayPoints   int     `json:"powerPlayPoints"`
				ShootingPctg      float64 `json:"shootingPctg"`
				ShorthandedGoals  int     `json:"shorthandedGoals"`
				ShorthandedPoints int     `json:"shorthandedPoints"`
				Shots             int     `json:"shots"`
			} `json:"subSeason"`
			Career struct {
				Assists           int     `json:"assists"`
				GameWinningGoals  int     `json:"gameWinningGoals"`
				GamesPlayed       int     `json:"gamesPlayed"`
				Goals             int     `json:"goals"`
				OtGoals           int     `json:"otGoals"`
				Pim               int     `json:"pim"`
				PlusMinus         int     `json:"plusMinus"`
				Points            int     `json:"points"`
				PowerPlayGoals    int     `json:"powerPlayGoals"`
				PowerPlayPoints   int     `json:"powerPlayPoints"`
				ShootingPctg      float64 `json:"shootingPctg"`
				ShorthandedGoals  int     `json:"shorthandedGoals"`
				ShorthandedPoints int     `json:"shorthandedPoints"`
				Shots             int     `json:"shots"`
			} `json:"career"`
		} `json:"regularSeason"`
	} `json:"featuredStats"`
	CareerTotals struct {
		RegularSeason struct {
			Assists            int     `json:"assists"`
			AvgToi             string  `json:"avgToi"`
			FaceoffWinningPctg float64 `json:"faceoffWinningPctg"`
			GameWinningGoals   int     `json:"gameWinningGoals"`
			GamesPlayed        int     `json:"gamesPlayed"`
			Goals              int     `json:"goals"`
			OtGoals            int     `json:"otGoals"`
			Pim                int     `json:"pim"`
			PlusMinus          int     `json:"plusMinus"`
			Points             int     `json:"points"`
			PowerPlayGoals     int     `json:"powerPlayGoals"`
			PowerPlayPoints    int     `json:"powerPlayPoints"`
			ShootingPctg       float64 `json:"shootingPctg"`
			ShorthandedGoals   int     `json:"shorthandedGoals"`
			ShorthandedPoints  int     `json:"shorthandedPoints"`
			Shots              int     `json:"shots"`
		} `json:"regularSeason"`
		Playoffs struct {
			Assists            int     `json:"assists"`
			AvgToi             string  `json:"avgToi"`
			FaceoffWinningPctg float64 `json:"faceoffWinningPctg"`
			GameWinningGoals   int     `json:"gameWinningGoals"`
			GamesPlayed        int     `json:"gamesPlayed"`
			Goals              int     `json:"goals"`
			OtGoals            int     `json:"otGoals"`
			Pim                int     `json:"pim"`
			PlusMinus          int     `json:"plusMinus"`
			Points             int     `json:"points"`
			PowerPlayGoals     int     `json:"powerPlayGoals"`
			PowerPlayPoints    int     `json:"powerPlayPoints"`
			ShootingPctg       float64 `json:"shootingPctg"`
			ShorthandedGoals   int     `json:"shorthandedGoals"`
			ShorthandedPoints  int     `json:"shorthandedPoints"`
			Shots              int     `json:"shots"`
		} `json:"playoffs"`
	} `json:"careerTotals"`
	ShopLink    string `json:"shopLink"`
	TwitterLink string `json:"twitterLink"`
	WatchLink   string `json:"watchLink"`
	Last5Games  []struct {
		Assists          int    `json:"assists"`
		GameDate         string `json:"gameDate"`
		GameID           int    `json:"gameId"`
		GameTypeID       int    `json:"gameTypeId"`
		Goals            int    `json:"goals"`
		HomeRoadFlag     string `json:"homeRoadFlag"`
		OpponentAbbrev   string `json:"opponentAbbrev"`
		Pim              int    `json:"pim"`
		PlusMinus        int    `json:"plusMinus"`
		Points           int    `json:"points"`
		PowerPlayGoals   int    `json:"powerPlayGoals"`
		Shifts           int    `json:"shifts"`
		ShorthandedGoals int    `json:"shorthandedGoals"`
		Shots            int    `json:"shots"`
		TeamAbbrev       string `json:"teamAbbrev"`
		Toi              string `json:"toi"`
	} `json:"last5Games"`
	Awards []struct {
		Trophy struct {
			Default string `json:"default"`
			Fr      string `json:"fr"`
		} `json:"trophy,omitempty"`
		Seasons []struct {
			Assists      int `json:"assists"`
			BlockedShots int `json:"blockedShots"`
			GameTypeID   int `json:"gameTypeId"`
			GamesPlayed  int `json:"gamesPlayed"`
			Goals        int `json:"goals"`
			Hits         int `json:"hits"`
			Pim          int `json:"pim"`
			PlusMinus    int `json:"plusMinus"`
			Points       int `json:"points"`
			SeasonID     int `json:"seasonId"`
		} `json:"seasons"`
	} `json:"awards"`
}

type PlayerSearchResults []struct {
	PlayerID       string `json:"playerId"`
	Name           string `json:"name"`
	PositionCode   string `json:"positionCode"`
	TeamAbbrev     string `json:"teamAbbrev"`
	Height         string `json:"height"`
	WeightInPounds int    `json:"weightInPounds"`
	BirthCountry   string `json:"birthCountry"`
}
