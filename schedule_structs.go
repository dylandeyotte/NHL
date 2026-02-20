package main

import (
	"time"
)

type Schedule struct {
	CalendarURL   string `json:"calendarUrl"`
	ClubTimezone  string `json:"clubTimezone"`
	ClubUTCOffset string `json:"clubUTCOffset"`
	Games         []struct {
		ID       int    `json:"id"`
		Season   int    `json:"season"`
		GameType int    `json:"gameType"`
		GameDate string `json:"gameDate"`
		Venue    struct {
			Default string `json:"default"`
		} `json:"venue"`
		NeutralSite       bool      `json:"neutralSite"`
		StartTimeUTC      time.Time `json:"startTimeUTC"`
		EasternUTCOffset  string    `json:"easternUTCOffset"`
		VenueUTCOffset    string    `json:"venueUTCOffset"`
		VenueTimezone     string    `json:"venueTimezone"`
		GameState         string    `json:"gameState"`
		GameScheduleState string    `json:"gameScheduleState"`
		AwayTeam          struct {
			ID         int `json:"id"`
			CommonName struct {
				Default string `json:"default"`
			} `json:"commonName"`
			PlaceName struct {
				Default string `json:"default"`
			} `json:"placeName"`
			PlaceNameWithPreposition struct {
				Default string `json:"default"`
				Fr      string `json:"fr"`
			} `json:"placeNameWithPreposition"`
			Abbrev         string `json:"abbrev"`
			Logo           string `json:"logo"`
			DarkLogo       string `json:"darkLogo"`
			AwaySplitSquad bool   `json:"awaySplitSquad"`
			Score          int    `json:"score"`
		} `json:"awayTeam,omitempty"`
		HomeTeam struct {
			ID         int `json:"id"`
			CommonName struct {
				Default string `json:"default"`
			} `json:"commonName"`
			PlaceName struct {
				Default string `json:"default"`
			} `json:"placeName"`
			PlaceNameWithPreposition struct {
				Default string `json:"default"`
				Fr      string `json:"fr"`
			} `json:"placeNameWithPreposition"`
			Abbrev         string `json:"abbrev"`
			Logo           string `json:"logo"`
			DarkLogo       string `json:"darkLogo"`
			HomeSplitSquad bool   `json:"homeSplitSquad"`
			HotelLink      string `json:"hotelLink"`
			HotelDesc      string `json:"hotelDesc"`
			Score          int    `json:"score"`
		} `json:"homeTeam,omitempty"`
	} `json:"games"`
}
