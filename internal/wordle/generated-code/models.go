// Code generated by sqlc. DO NOT EDIT.

package wordle

import (
	"database/sql"
	"time"
)

type Account struct {
	DiscordID string `json:"discord_id"`
	TimeZone  string `json:"time_zone"`
}

type DisableQuip struct {
	ServerID string `json:"server_id"`
}

type Nickname struct {
	DiscordID string `json:"discord_id"`
	ServerID  string `json:"server_id"`
	Nickname  string `json:"nickname"`
}

type Quip struct {
	ID                 int64          `json:"id"`
	ScoreValue         int32          `json:"score_value"`
	Quip               string         `json:"quip"`
	InsideJoke         bool           `json:"inside_joke"`
	InsideJokeServerID sql.NullString `json:"inside_joke_server_id"`
	CreatedByAccount   string         `json:"created_by_account"`
	CreatedAt          time.Time      `json:"created_at"`
	Uses               int32          `json:"uses"`
}

type WordleScore struct {
	ID        int64     `json:"id"`
	DiscordID string    `json:"discord_id"`
	GameID    int32     `json:"game_id"`
	Guesses   int32     `json:"guesses"`
	CreatedAt time.Time `json:"created_at"`
}
