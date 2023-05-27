package cmd

import (
	"crypto-prices-bot/req"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

type Handler func(update tgbotapi.Update, bot *tgbotapi.BotAPI)

func Start(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	keys := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Цена"),
			tgbotapi.NewKeyboardButton("Статистика"),
		),
	)

	r := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите 'Цена <Тикер1/Тикер2> или 'Статистика'. Например 'Цена ETH/USDT'")

	r.ReplyMarkup = keys

	bot.Send(r)
}

func Convert(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	r := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	r.ReplyToMessageID = update.Message.MessageID

	args := strings.Split(update.Message.Text, " ")
	if len(args) < 2 {
		r.Text = "Введите запрос в формате 'Цена ETH/BTC'"

		bot.Send(r)
		return
	}

	pair := args[1]
	pair = strings.Replace(pair, "/", "", -1)

	data, err := req.Fetch(pair)
	if err != nil {
		log.Println(err)

		r.Text = "Ошибка! Проверьте правильность тикера"

		bot.Send(r)
		return
	}

	r.Text = fmt.Sprintf("🐙 Kraken: %.2f", data.Result)

	bot.Send(r)

	go func() {
		err = req.RecordStats(update.Message.From.ID, pair)
		if err != nil {
			log.Println(err)
		}
	}()
}

func Stats(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	res := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	res.ReplyToMessageID = update.Message.MessageID

	stats, err := req.ShowStats(update.Message.From.ID)
	if err != nil {
		log.Println(err)

		res.Text = "Не удалось получить статистику. Возможно, вы не отправляли запросов"

		bot.Send(res)
		return
	}

	res.Text = fmt.Sprintf("⏱ Первый запрос: %s\n📈 Кол-во запросов: %d\n🔗 Последняя пара: %s",
		stats.FirstReq.UTC().Format(time.UnixDate), stats.ReqAmount, stats.LastPair)

	bot.Send(res)
}
