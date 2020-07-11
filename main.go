package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// First is the default
var allowedSubdomains = []string{
	"www",
	"dlc",
	"dlc2",
	"p2music",
	"p1",
	"p1music",
	"tf2",
	"tf2music",
}

func subdomainAllowed(subdomain string) bool {
	for _, a := range allowedSubdomains {
		if a == subdomain {
			return true
		}
	}
	return false
}

func main() {
	tok := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(tok)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery == nil { // We only care about inline queries
			continue
		}

		useSubdomain := allowedSubdomains[0]
		query := update.InlineQuery.Query

		split := strings.Split(update.InlineQuery.Query, ":")
		if len(split) == 2 && subdomainAllowed(strings.TrimSpace(split[0])) {
			useSubdomain = strings.TrimSpace(split[0])
			query = strings.TrimSpace(split[1])
		}

		sounds := GetSounds(SoundFilterOptions{
			Quote: query,
		}, useSubdomain, 1)

		var results []interface{}

		for _, sound := range sounds {
			res := tgbotapi.NewInlineQueryResultAudio(
				sound.ID,
				sound.Link(),
				sound.Who,
			)

			res.Caption = sound.Text
			res.Performer = sound.Text
			results = append(results, res)
		}

		if len(results) > 50 {
			results = results[0:49]
		}

		bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			Results:       results,
			CacheTime:     30,
			IsPersonal:    false,
		})
	}

	// sounds := GetSounds(SoundFilterOptions{
	// 	Quote: "test",
	// }, 1)
	//
	// for _, sound := range sounds {
	// 	fmt.Printf("%s: %s (%s)\n", sound.Who, Ellipsize(sound.Text, 1024), sound.Link())
	// }
}
