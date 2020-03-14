package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

		sounds := GetSounds(SoundFilterOptions{
			Quote: update.InlineQuery.Query,
		}, 1)

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
