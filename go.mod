module DiscordWordle

go 1.14

// Comment below is needed for heroku-buildpack-go https://github.com/heroku/heroku-buildpack-go/issues/301

// +heroku goVersion go1.14

require (
	github.com/bwmarrin/discordgo v0.25.0
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/lib/pq v1.10.5
	github.com/rs/zerolog v1.26.1
	github.com/rubenv/sql-migrate v1.0.0
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f // indirect
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
)
