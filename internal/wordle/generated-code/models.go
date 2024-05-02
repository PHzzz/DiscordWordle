// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package wordle

import (
	"time"
)

type Account struct {
	DiscordID string `json:"discord_id"`
	TimeZone  string `json:"time_zone"`
}

type Nickname struct {
	DiscordID string `json:"discord_id"`
	ServerID  string `json:"server_id"`
	Nickname  string `json:"nickname"`
}

type WordleScore struct {
	ID        int64     `json:"id"`
	DiscordID string    `json:"discord_id"`
	GameID    time.Time `json:"game_id"`
	GameType  int32     `json:"game_type"`
	Guesses   int32     `json:"guesses"`
	CreatedAt time.Time `json:"created_at"`
}
