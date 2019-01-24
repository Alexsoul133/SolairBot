package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("430629496:AAHDBvxHimRzeURldxAz_4v8pp4bKzoeH8s")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	reply := ""

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		// if update.ChannelPost.Text == "!sun" {
		// 	bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
		// 	continue
		// }
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "!sun" {
			bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я Солер из Асторы, воин Света, верный слуга короля жёлтого замка Попс Маэллард."
			continue
		case "sun":
			bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
			continue
		case "info":
			reply = fmt.Sprintf("@%v \nID: %v \nMessageID: ^v\nTimeMsg:  %v \nTimeNow: %v",
				update.Message.ReplyToMessage.From.UserName,
				update.Message.ReplyToMessage.From.ID,

				update.Message.ReplyToMessage.Time().Format(time.RFC1123Z),
				time.Now().Format(time.RFC1123Z))
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
