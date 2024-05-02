package main

import (
	wordle "DiscordWordle/internal/wordle/generated-code"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// Variables used for command line parameters
var (
	Token       string
	DatabaseUrl string
)

type response struct {
	Text  string
	Emoji string
}

var db *sql.DB

const cmdHistory = "history"

// const cmdUpdate = "update"
const cmdScoreboard = "scoreboard"
const cmdPreviousWeek = "previous"

/* const cmdQuip = "quip"
const cmdQuipEnable = "enable"
const cmdQuipDisable = "disable"
const cmdQuipList = "list"
const cmdQuipDelete = "delete"
const cmdTimeZone = "timezone" */
const cmdWordle = "Vārdulis"
const noSolutionResult = "X"
const hardModeIndicator = "*"
const noSolutionGuesses = 7

/* var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
	}
) */

func init() {
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Fatal().Msg("DISCORD_TOKEN must be set")
	}

	DatabaseUrl = os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		log.Fatal().Msg("DATABASE_URL must be set")
	}

	dbConnection, err := sql.Open("postgres", DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot connect to database: %s", DatabaseUrl)
	}

	db = dbConnection
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating Discord session")
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Add slash command handler
	/* dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	}) */

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("error opening connection to discord over websocket")
	}

	// Add slash commands
	/* registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create '%v' command", v.Name)
		}
		log.Info().Msgf("Created '%v' command", v.Name)
		registeredCommands[i] = cmd
	} */

	// Wait here until CTRL-C or other term signal is received.
	log.Info().Msg("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt /*, os.Kill*/)
	<-sc

	// Cleanly close down the Discord session.
	defer dg.Close()

	// Remove own slash commands
	/* for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot delete '%v' command", v.Name)
		} else {
			log.Info().Msgf("Deleted '%v' command", v.Name)
		}
	} */
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var r response

	// Ignore all messages created by the bot itself
	botName := s.State.User.Username
	if m.Author.ID == s.State.User.ID {
		return
	}

	tokenizedContent, err := m.ContentWithMoreMentionsReplaced(s)
	if err != nil {
		log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to replace mentions")
		return
	}

	botMentionToken := fmt.Sprintf("@%s", botName)
	wordleScoreDetected, _ := mentionlessWordleScoreDetection(tokenizedContent)
	if strings.HasPrefix(tokenizedContent, botMentionToken) || wordleScoreDetected {
		input := strings.TrimSpace(strings.Replace(tokenizedContent, botMentionToken, "", 1))
		q := wordle.New(db)
		ctx := context.Background()

		r.Emoji = "❌"
		existingAccount, err := q.CountAccountsByDiscordId(ctx, m.Author.ID)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to count accounts by Discord ID")
			r.Text = "Ak, nē! Viss zudis! Punktulis salūzis!"
			r.Emoji = "🔥"
			flushEmojiAndResponseToDiscord(s, m, r)
			return
		}

		existingNickname, err := q.CountNicknameByDiscordIdAndServerId(ctx, wordle.CountNicknameByDiscordIdAndServerIdParams{
			DiscordID: m.Author.ID,
			ServerID:  m.GuildID,
		})
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to count nicknames by Discord ID and server ID")
			r.Text = "Ak, nē! Viss zudis! Punktulis salūzis!"
			r.Emoji = "🔥"
			flushEmojiAndResponseToDiscord(s, m, r)
			return
		}

		var account wordle.Account
		if m.Message.GuildID != "" {
			account = getOrCreateAccount(ctx, s, m, existingAccount, existingNickname, q)
		} else {
			account = wordle.Account{
				DiscordID: m.Message.Author.ID,
				TimeZone:  "Europe/Riga",
			}
		}

		routeMessageToAction(ctx, s, m, input, account, q, botMentionToken)
	}
}

func routeMessageToAction(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, input string, account wordle.Account, q *wordle.Queries, botMentionToken string) {
	var r response

	if strings.Contains(input, cmdWordle) {
		gameId, guesses, err := extractGameGuesses(input)
		if err != nil {
			log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdWordle).Err(err).Msg("Error parsing guess count")
			r.Emoji = "❌"
			r.Text = ""
			flushEmojiAndResponseToDiscord(s, m, r)
		} else {
			log.Info().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdWordle).Int("guesses", guesses).Time("game_id", gameId).Msg("Found a Wordle")
			persistScore(ctx, m, s, account, gameId, guesses)
		}

		/* } else if strings.HasPrefix(input, cmdUpdate) {
		gameId, guesses, err := extractGameGuesses(input)
		if err != nil {
			log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdUpdate).Err(err).Msg("Error parsing guess count")
		}
		log.Info().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdUpdate).Int("guesses", guesses).Int("game_id", gameId).Msg("Updated a Wordle")
		updateExistingScore(ctx, m, s, account, gameId, guesses) */
	} else if strings.HasPrefix(input, cmdHistory) {
		getHistory(ctx, m, s, account)
		/* } else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipEnable) {
			enableQuips(ctx, m, s)
		} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipDisable) {
			disableQuips(ctx, m, s)
		} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipList) {
			listQuips(ctx, m, s)
		} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipDelete) {
			quipId, err := extractQuipId(input)
			if err != nil {
				log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdQuip).Err(err).Msg("Error parsing quipId for deletion")
			} else {
				deleteQuip(ctx, m, s, quipId)
			}
		} else if strings.HasPrefix(input, cmdQuip) {
			score, quip, err := extractScoreQuip(input)
			if err != nil {
				log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdQuip).Err(err).Msg("Error parsing quip")
			} else {
				persistQuip(ctx, m, s, account, score, quip)
			} */
	} else if strings.HasPrefix(input, cmdScoreboard+" "+cmdPreviousWeek) {
		getPreviousScoreboard(ctx, m, s)
	} else if strings.HasPrefix(input, cmdScoreboard) {
		getScoreboard(ctx, m, s)
		/* } else if strings.HasPrefix(input, cmdTimeZone) {
		updateAccountTimeZone(ctx, input, cmdTimeZone, s, m, q, account) */
	} else if strings.HasPrefix(input, "help") {
		helpResponse(s, m, botMentionToken)
	} else {
		log.Info().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", "").Msg("Failed to match command")
		r.Text = ""
		r.Emoji = "🤷"
		flushEmojiAndResponseToDiscord(s, m, r)
	}
}

/* func extractScoreQuip(input string) (int, string, error) {
	var dataExp = regexp.MustCompile(`(?P<score>\d+)\s(?P<quip>.+)`)

	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return 0, "", err
	}

	score, _ := strconv.Atoi(result["score"])
	return score, result["quip"], nil
}

func extractQuipId(input string) (int, error) {
	var dataExp = regexp.MustCompile(`(?P<quipId>\d+)`)

	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return 0, err
	}

	quipId, _ := strconv.Atoi(result["quipId"])
	return quipId, nil
} */

func extractGameGuesses(input string) (time.Time, int, error) {
	var dataExp = regexp.MustCompile(fmt.Sprintf(`(?P<game_id>\d{2}/\d{2}/\d{4})\s(?P<guesses>\d+|%s)`, noSolutionResult))
	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return time.Time{}, 0, err
	}
	gameId, _ := time.Parse("02/01/2006", result["game_id"])
	// allow only today's scores
	today := time.Now()
	if !(gameId.Year() == today.Year() && gameId.YearDay() == today.YearDay()) {
		errorMessage := fmt.Sprintf("%s is not today", gameId.String())
		return time.Time{}, 0, errors.New(errorMessage)
	}
	var guesses int
	/* if strings.ToUpper(result["guesses"]) == noSolutionResult {
		guesses = noSolutionGuesses
	} else { */
	guesses, _ = strconv.Atoi(result["guesses"])
	/* } */
	// check if emojis match guesses number
	var emojiExp = regexp.MustCompile(`(?m)^[⬛🟨🟩]{5}$`)
	match := emojiExp.FindAllString(input, -1)
	if len(match) == 0 {
		return time.Time{}, 0, errors.New("emoji lines not found")
	} else if len(match) != guesses {
		errorMessage := fmt.Sprintf("%d guesses didn't match %d emoji lines", guesses, len(match))
		return time.Time{}, 0, errors.New(errorMessage)
	}
	// check if last line is 5 greens if winning
	var winLine = strings.Repeat("🟩", 5)
	if match[len(match)-1] != winLine {
		if guesses == 6 {
			guesses = noSolutionGuesses
			// check if winning line is not higher if lost
			for _, v := range match {
				if v == winLine {
					return time.Time{}, 0, errors.New("6 guesses but winning emoji line not at the end")
				}
			}
		} else {
			errorMessage := fmt.Sprintf("%d guesses but didn't find winning emoji line at the end", guesses)
			return time.Time{}, 0, errors.New(errorMessage)
		}
	}
	// check for multiple winning lines
	var winExp = regexp.MustCompile(`(?m)^🟩{5}$`)
	match = winExp.FindAllString(input, -1)
	if len(match) > 1 {
		errorMessage := fmt.Sprintf("Found %d winning emoji lines", len(match))
		return time.Time{}, 0, errors.New(errorMessage)
	}

	return gameId, guesses, nil
}

func mentionlessWordleScoreDetection(input string) (bool, error) {
	var dataExp = regexp.MustCompile(fmt.Sprintf(`Vārdulis (?P<game_id>\d{2}/\d{2}/\d{4})\s(?P<guesses>\d+|%s)/6[\%s]?\n`, noSolutionResult, hardModeIndicator))
	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

func matchGroupsToStringMap(input string, dataExp *regexp.Regexp) (map[string]string, error) {
	match := dataExp.FindStringSubmatch(input)
	result := make(map[string]string)
	if len(match) == 0 {
		errorMessage := fmt.Sprintf("%s didn't match %s", input, dataExp)
		return result, errors.New(errorMessage)

	}
	for i, name := range dataExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, nil
}

func flushEmojiAndResponseToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, r response) {
	reactToMessage(s, m, r.Emoji)
	respondAsNewMessage(s, m, r.Text)
}

func respondAsNewMessage(s *discordgo.Session, m *discordgo.MessageCreate, response string) {
	if response != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Str("response", response).Msg("Error responding")
		}
	}
}

func reactToMessage(s *discordgo.Session, m *discordgo.MessageCreate, reactionEmoji string) {
	if reactionEmoji != "" {
		err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, reactionEmoji)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Str("reaction", reactionEmoji).Msg("Error reacting")
		}
	}
}
