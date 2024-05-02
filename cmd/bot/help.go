package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func helpResponse(s *discordgo.Session, m *discordgo.MessageCreate, botMentionToken string) {
	var response response
	response.Text = fmt.Sprintf("`%s` - piereģistrēt šodienas Vārduļa rezultātu. Katrai spēlei drīkst piereģistrēt tikai vienu rezultātu.\n"+
		"`%s` - apskatīties savus iepriekšējos Vārduļa rezultātus\n"+
		"`%s` - apskatīties savu un draudziņu rezultātu tabulu\n"+
		"`%s` - apskatīties pagājušās nedēļas rekordistus\n",
		fmt.Sprintf("%s %s 15/12/2021 5/6 <emoji blocks>", botMentionToken, cmdWordle),
		fmt.Sprintf("%s %s", botMentionToken, cmdHistory),
		fmt.Sprintf("%s %s", botMentionToken, cmdScoreboard),
		fmt.Sprintf("%s %s %s", botMentionToken, cmdScoreboard, cmdPreviousWeek),
	)

	response.Emoji = "💁"

	flushEmojiAndResponseToDiscord(s, m, response)
}
