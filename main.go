package main

import (
	"fmt"
	"log"
	"os"

	botInstance "github.com/crocone/tg-bot"
	"github.com/joho/godotenv"
)

var bot *botInstance.BotAPI

type button struct {
	name string
	data string
}

func startMenu() botInstance.InlineKeyboardMarkup {
	states := []button{
		{
			name: "Привет",
			data: "hi",
		},
		{
			name: "Пока",
			data: "buy",
		},
	}

	buttons := make([][]botInstance.InlineKeyboardButton, len(states))
	for index, state := range states {
		buttons[index] = botInstance.NewInlineKeyboardRow(botInstance.NewInlineKeyboardButtonData(state.name, state.data))
	}

	return botInstance.NewInlineKeyboardMarkup(buttons...)
}

func main() {
	log.Print("Bot has been started !")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	bot, err = botInstance.NewBotAPI(os.Getenv("TG_API_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	u := botInstance.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to start listening for updates %v", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callbacks(update)
		} else if update.Message.IsCommand() {
			commands(update)
		} else {
			// simply message
		}
	}
}

func callbacks(update botInstance.Update) {
	data := update.CallbackQuery.Data
	chatId := update.CallbackQuery.From.ID
	firstName := update.CallbackQuery.From.FirstName
	lastName := update.CallbackQuery.From.LastName
	var text string
	switch data {
	case "hi":
		text = fmt.Sprintf("Привет %v %v", firstName, lastName)
	case "buy":
		text = fmt.Sprintf("Пока %v %v", firstName, lastName)
	default:
		text = "Неизвестная команда"
	}
	msg := botInstance.NewMessage(chatId, text)
	sendMessage(msg)
}

func commands(update botInstance.Update) {
	command := update.Message.Command()
	switch command {
	case "start":
		msg := botInstance.NewMessage(update.Message.Chat.ID, "Выберите действие")
		msg.ReplyMarkup = startMenu()
		msg.ParseMode = "Markdown"
		sendMessage(msg)
	}
}

func sendMessage(msg botInstance.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Panicf("Send message error: %v", err)
	}
}
