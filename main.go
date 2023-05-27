package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"crypto-prices-bot/cmd"
	"crypto-prices-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func envValue(name string, desc string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf(desc, "not found.")
	}

	return value
}

func main() {
	// init db credentials
	port, err := strconv.Atoi(envValue("DB_PORT", "DB port"))
	if err != nil {
		log.Fatalf("Can't read port: %s", err)
	}

	dbCred := database.Credentials{}
	dbCred.Host = envValue("DB_HOST", "db hostname")
	dbCred.Port = port
	dbCred.User = envValue("DB_USER", "db username")
	dbCred.Password = envValue("DB_PASSWORD", "db password")
	dbCred.Database = envValue("DB_NAME", "db name")

	err = database.Connect(dbCred)
	if err != nil {
		log.Fatalf("Can't connect to db: %s", err)
	}

	log.Println("Connected to the database.")

	defer database.Disconnect()

	// init bot
	bot, err := tgbotapi.NewBotAPI(envValue("BOT_TOKEN", "telegram bot token"))
	if err != nil {
		log.Fatalf("can't connect to the bot: %s", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// create handlers
	handlers := make(map[string]cmd.Handler)

	handlers["старт"] = cmd.Start
	handlers["цена"] = cmd.Convert
	handlers["статистика"] = cmd.Stats

	for update := range updates {
		if update.Message != nil {
			payload := update.Message.Text
			arg := strings.Split(payload, " ")

			command := strings.ToLower(arg[0])

			if handlers[command] != nil {
				go handlers[command](update, bot)
			} else {
				go handlers["старт"](update, bot)
			}
		}
	}
}
