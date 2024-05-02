package main

import (
	wordle "DiscordWordle/internal/wordle/generated-code"
	"bytes"
	"context"

	// "database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func persistScore(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a wordle.Account, gameId time.Time, guesses int) {
	response, scoreObj := buildScoreObjFromInput(a, gameId, guesses)

	scoreParams := wordle.CreateScoreParams{
		DiscordID: a.DiscordID,
		GameID:    scoreObj.GameID,
		Guesses:   scoreObj.Guesses,
	}

	q := wordle.New(db)
	_, err := q.CreateScore(ctx, scoreParams)

	if err != nil {
		log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to persist score")
		response.Emoji = "⛔"
		/* serverHasDisabledQuips, _ := q.CheckIfServerHasDisabledQuips(ctx, m.GuildID)
		if len(serverHasDisabledQuips) == 0 {
			response.Text = "You already created a score for this game, try updating it if it's wrong"
		} else { */
		response.Text = ""
		/* } */
	} else {
		response = scoreColorfulResponse(guesses, ctx, m)
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

/* func enableQuips(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session) {
	var response response

	q := wordle.New(db)
	err := q.EnableQuipsForServer(ctx, m.GuildID)

	if err != nil {
		log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to enable quips")
		response.Text = "Error enabling quips"
		response.Emoji = "💣"
	}

	response.Text = "Prepare to laugh to death at these mad jokes"
	response.Emoji = "💭"

	flushEmojiAndResponseToDiscord(s, m, response)
}

func disableQuips(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session) {
	var response response

	q := wordle.New(db)
	err := q.DisableQuipsForServer(ctx, m.GuildID)

	if err != nil {
		log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to disable quips")
		response.Text = "Error disabling quips"
		response.Emoji = "💣"
	}

	response.Text = "" //No response, only emoji
	response.Emoji = "😶"

	flushEmojiAndResponseToDiscord(s, m, response)
}

func listQuips(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session) {
	var response response

	q := wordle.New(db)
	quips, _ := q.GetQuipsByServerId(ctx, sql.NullString{String: m.GuildID, Valid: true})

	for len(quips) > 0 {
		var buf bytes.Buffer
		var responseSize int
		responseSize = 0

		w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tGuesses\tQuip\t")

		for _, v := range quips {
			quipLine := fmt.Sprintf("%d\t%d\t%s\t", v.ID, v.ScoreValue, v.Quip)
			responseSize = responseSize + len(quipLine)
			if responseSize > 500 {
				break
			}
			_, _ = fmt.Fprintln(w, quipLine)
			quips = quips[1:]
		}
		_ = w.Flush()

		response.Text = fmt.Sprintf("```\n%s\n```", buf.String())
		flushEmojiAndResponseToDiscord(s, m, response)
	}
}

func deleteQuip(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, quipId int) {
	var response response

	q := wordle.New(db)
	p := wordle.DeleteQuipByIdAndServerIdParams{
		ID:                 int64(quipId),
		InsideJokeServerID: sql.NullString{String: m.GuildID, Valid: true},
	}
	err := q.DeleteQuipByIdAndServerId(ctx, p)

	if err != nil {
		log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to delete quip")
		response.Emoji = "⁉️"
		response.Text = "Failed to delete quip"
		flushEmojiAndResponseToDiscord(s, m, response)
		return
	}

	response.Emoji = "✌️"
	response.Text = "Quip has been deleted"
	flushEmojiAndResponseToDiscord(s, m, response)
}

func persistQuip(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, account wordle.Account, score int, quip string) {
	var nicknames []wordle.Nickname
	if m.GuildID == "" {
		q := wordle.New(db)
		nicknames, _ = q.GetNicknamesByDiscordId(ctx, account.DiscordID)
	} else {
		nicknames = append(nicknames, wordle.Nickname{
			DiscordID: account.DiscordID,
			ServerID:  m.GuildID,
			Nickname:  m.Member.Nick,
		})
	}

	var response response
	for _, nick := range nicknames {
		quipParams := wordle.CreateQuipForScoreParams{
			ScoreValue:         int32(score),
			Quip:               quip,
			InsideJoke:         true,
			InsideJokeServerID: sql.NullString{String: nick.ServerID, Valid: true},
			CreatedByAccount:   nick.DiscordID,
		}

		q := wordle.New(db)
		_, err := q.CreateQuipForScore(ctx, quipParams)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to create quip")
			response.Emoji = "⁉️"
			response.Text = "Them words are not right"
			flushEmojiAndResponseToDiscord(s, m, response)
			return
		}
	}

	response.Emoji = "🤣"
	flushEmojiAndResponseToDiscord(s, m, response)
} */

func getHistory(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a wordle.Account) {

	historyByAccountParams := wordle.GetScoreHistoryByAccountParams{
		DiscordID: a.DiscordID,
		ServerID:  m.GuildID,
	}

	q := wordle.New(db)
	scores, err := q.GetScoreHistoryByAccount(ctx, historyByAccountParams)

	var response response

	if err != nil {
		response.Emoji = "⁉️"
		response.Text = "Neatrodu nevienu rezultātu"
	} else {
		response.Emoji = "👍"
		response.Text = fmt.Sprintf("Atradu %d rezultātus, priekšniek!", len(scores))
		for _, v := range scores {
			response.Text += fmt.Sprintf("\n Spēle: %s - %d/6", v.GameID.Format("02/01/2006"), v.Guesses)
		}
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

func findMinAndMax(a []time.Time) (min time.Time, max time.Time) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value.Before(min) {
			min = value
		}
		if value.After(max) {
			max = value
		}
	}
	return min, max
}

func getScoreboard(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session) {
	q := wordle.New(db)
	scores, err := q.GetScoresByServerId(ctx, m.GuildID)
	expectedGames, _ := q.GetExpectedWeekGames(ctx, m.GuildID)

	var response response

	if err != nil {
		response.Emoji = "⁉️"
		response.Text = "Neatradu nevienu rezultātu"
	} else {
		response.Emoji = "🔢"

		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

		var maxNumOfGames int
		maxNumOfGames = 0
		_, _ = fmt.Fprintln(w, "Vieta\tVārds\tMinējumi\tPunkti\t")
		for _, v := range scores {

			displayGameGuesses := dashDisplayForMissingScores(expectedGames, v)

			if int(v.GamesCount) > maxNumOfGames {
				maxNumOfGames = int(v.GamesCount)
			}
			_, _ = fmt.Fprintf(w, "%d.\t%s\t%s\t%d\t\n", v.Position, v.Nickname, displayGameGuesses, v.Total)
		}

		var lwBuf bytes.Buffer
		var lwMinDate time.Time
		var lwMaxDate time.Time
		lw := tabwriter.NewWriter(&lwBuf, 0, 0, 3, ' ', 0)
		if maxNumOfGames == 1 {
			lastWeekScores, _ := q.GetScoresByServerIdPreviousWeek(ctx, m.GuildID)
			lastWeekExpectedGames, _ := q.GetExpectedPreviousWeekGames(ctx, m.GuildID)
			lwMinDate, lwMaxDate = findMinAndMax(lastWeekExpectedGames)
			_, _ = fmt.Fprintln(lw, "Vieta\tVārds\tMinējumi\tPunkti\t")
			for _, lwv := range lastWeekScores {
				displayGameGuesses := dashDisplayForMissingScores(lastWeekExpectedGames, wordle.GetScoresByServerIdRow(lwv))
				_, _ = fmt.Fprintf(lw, "%d.\t%s\t%s\t%d\t\n", lwv.Position, lwv.Nickname, displayGameGuesses, lwv.Total)
			}
			_ = lw.Flush()
		}

		_ = w.Flush()
		if len(lwBuf.String()) > 0 {
			response.Text = fmt.Sprintf("**Šonedēļ:**\n```\n%s\n```\n**Pagājušonedēļ (%s - %s):**\n```\n%s\n```", buf.String(), lwMinDate.Format("02/01/2006"), lwMaxDate.Format("02/01/2006"), lwBuf.String())
		} else {
			response.Text = fmt.Sprintf("```\n%s\n```", buf.String())
		}
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

func dashDisplayForMissingScores(expectedGames []time.Time, v wordle.GetScoresByServerIdRow) []string {
	var displayGameGuesses []string
	for _, g := range expectedGames {
		var nestedGameGuessesMap []map[string]int
		cleanGameGuesses := make(map[string]int)
		_ = json.Unmarshal(v.GameGuesses, &nestedGameGuessesMap)
		for _, gameGuess := range nestedGameGuessesMap {
			for stringGameId, guesses := range gameGuess {
				gameId, _ := time.Parse("2006-01-02", stringGameId)
				// TODO: check why doesn't work without reformatting time to string
				cleanGameGuesses[gameId.Format("2006-01-02")] = guesses
			}
		}
		if val, ok := cleanGameGuesses[g.Format("2006-01-02")]; ok {
			displayGameGuesses = append(displayGameGuesses, strconv.Itoa(val))
		} else {
			displayGameGuesses = append(displayGameGuesses, "-")
		}
	}
	return displayGameGuesses
}

func getPreviousScoreboard(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session) {
	q := wordle.New(db)
	scores, err := q.GetScoresByServerIdPreviousWeek(ctx, m.GuildID)
	lastWeekExpectedGames, _ := q.GetExpectedPreviousWeekGames(ctx, m.GuildID)
	lwMinDate, lwMaxDate := findMinAndMax(lastWeekExpectedGames)
	var response response
	var mentionList string
	var medal string

	if err != nil {
		response.Emoji = "⁉️"
		response.Text = "Neatradu nevienu rezultātu"
	} else {
		response.Emoji = "🔢"

		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

		_, _ = fmt.Fprintln(w, "Vieta\tVārds\tMinējumi\tPunkti\t")
		for _, v := range scores {
			displayGameGuesses := dashDisplayForMissingScores(lastWeekExpectedGames, wordle.GetScoresByServerIdRow(v))
			_, _ = fmt.Fprintf(w, "%d.\t%s\t%s\t%d\t\n", v.Position, v.Nickname, displayGameGuesses, v.Total)
			switch v.Position {
			case 1:
				medal = "🥇"
			case 2:
				medal = "🥈"
			case 3:
				medal = "🥉"
			default:
				medal = ""
			}
			mentionList = mentionList + medal + "<@!" + v.DiscordID + "> "
		}

		_ = w.Flush()

		response.Text = fmt.Sprintf("**Pagājušonedēļ (%s - %s):**\n```\n%s\n```%s", lwMinDate.Format("02/01/2006"), lwMaxDate.Format("02/01/2006"), buf.String(), mentionList)
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

/* func updateExistingScore(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a wordle.Account, gameId int, guesses int) {
	response, wordlecoreObj := buildScoreObjFromInput(a, gameId, guesses)

	priceParams := wordle.UpdateScoreParams{
		DiscordID: a.DiscordID,
		GameID:    wordlecoreObj.GameID,
		Guesses:   wordlecoreObj.Guesses,
	}

	q := wordle.New(db)
	_, err := q.UpdateScore(ctx, priceParams)

	if err != nil {
		response.Emoji = "⁉️"
		response.Text = "I didn't find an existing price."
	} else {
		response = scoreColorfulResponse(guesses, ctx, m)
	}

	flushEmojiAndResponseToDiscord(s, m, response)
} */

func buildScoreObjFromInput(a wordle.Account, gameId time.Time, guesses int) (response, wordle.WordleScore) {
	var response response

	scoreThing := wordle.WordleScore{
		DiscordID: a.DiscordID,
		GameID:    time.Time(gameId),
		Guesses:   int32(guesses),
	}

	return response, scoreThing
}

func scoreColorfulResponse(guesses int, ctx context.Context, m *discordgo.MessageCreate) response {
	var response response
	/* q := wordle.New(db)
	serverHasDisabledQuips, _ := q.CheckIfServerHasDisabledQuips(ctx, m.GuildID)
	if len(serverHasDisabledQuips) == 0 {
		response = selectResponseText(guesses, ctx, m, response)
	} */
	response = selectResponseEmoji(guesses, response)
	return response
}

/* func selectResponseText(guesses int, ctx context.Context, m *discordgo.MessageCreate, response response) response {
	if guesses >= 1 && guesses <= 6 || guesses == noSolutionGuesses {
		responseParams := wordle.GetQuipByScoreParams{
			ScoreValue:         int32(guesses),
			InsideJokeServerID: sql.NullString{String: m.GuildID, Valid: true},
		}

		q := wordle.New(db)
		r, _ := q.GetQuipByScore(ctx, responseParams)
		_ = q.IncrementQuip(ctx, r.ID)
		response.Text = r.Quip
	} else if guesses == 69 {
		response.Text = "nice."
	} else {
		response.Text = "Is that even a real number? Did you fail to guess it?"
	}

	return response
} */

func selectResponseEmoji(guesses int, response response) response {
	if guesses == 69 {
		response.Emoji = "♋️"
	} else if guesses == noSolutionGuesses {
		response.Emoji = "0️⃣" // "<:mau5:665049437688692739"
	} else if guesses == 1 {
		response.Emoji = "1️⃣"
	} else if guesses == 2 {
		response.Emoji = "2️⃣"
	} else if guesses == 3 {
		response.Emoji = "3️⃣"
	} else if guesses == 4 {
		response.Emoji = "4️⃣"
	} else if guesses == 5 {
		response.Emoji = "5️⃣"
	} else if guesses == 6 {
		response.Emoji = "6️⃣"
	} else {
		response.Emoji = "❌"
	}

	return response
}
