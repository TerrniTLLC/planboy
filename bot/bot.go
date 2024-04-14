package telegramnotionbot

import (
	"fmt"
	"log"
	"os"
	"strings"

	botInstance "github.com/crocone/tg-bot"

	"github.com/terrnitllc/planboy/notion"

	"github.com/joho/godotenv"
)

// var botApi *botInstance.BotAPI

type button struct {
	name string
	data string
}

func startMenu() botInstance.InlineKeyboardMarkup {
	states := []button{
		{
			name: "Today's habits",
			data: "get_page",
		},
	}

	buttons := make([][]botInstance.InlineKeyboardButton, len(states))
	for index, state := range states {
		buttons[index] = botInstance.NewInlineKeyboardRow(botInstance.NewInlineKeyboardButtonData(state.name, state.data))
	}

	return botInstance.NewInlineKeyboardMarkup(buttons...)
}

func (bot *TelegramBotStruct) startBot() {
	botApi, err := botInstance.NewBotAPI(bot.Token)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}
	bot.BotApi = botApi

	u := botInstance.NewUpdate(0)
	u.Timeout = 60
	updates := bot.BotApi.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to start listening for updates %v", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			bot.callbacks(update)
		} else if update.Message.IsCommand() {
			bot.commands(update)
		} else {
			// simply message
		}
	}
}

type TelegramBotStruct struct {
	BotApi *botInstance.BotAPI
	Notion *notion.NotionApi
	Token  string
}

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	telegramApiCode, ok := os.LookupEnv("TG_API_BOT_TOKEN")
	if !ok {
		log.Fatalf("TG_API_BOT_TOKEN is not set")
	}

	tgBot := TelegramBotStruct{BotApi: nil, Notion: nil, Token: telegramApiCode}

	// Notion client
	notionApiSecret, ok := os.LookupEnv("NOTION_API_SECRET")
	if !ok {
		log.Fatalf("NOTION_API_SECRET is not set")
	}

	notionDbKey, ok := os.LookupEnv("NOTION_DB_KEY")
	if !ok {
		log.Fatalf("NOTION_DB_KEY is not set")
	}

	notion, err := notion.NewNotionApi(notionApiSecret, notionDbKey)
	if err != nil {
		log.Fatalf("Failed to initialize Notion: %v", err)
	}

	tgBot.Notion = notion
	tgBot.startBot()

	log.Print("Bot has been started !")
}

func (bot *TelegramBotStruct) callbacks(update botInstance.Update) {
	data := update.CallbackQuery.Data
	chatId := update.CallbackQuery.From.ID
	firstName := update.CallbackQuery.From.FirstName
	// lastName := update.CallbackQuery.From.LastName
	notionPage := strings.Split(bot.Notion.MainPage.URL, "")
	var text string

	switch data {
	case "get_page":
		text = fmt.Sprintf("Привет %v %v", firstName, notionPage)
	default:
		text = "Неизвестная команда"
	}
	msg := botInstance.NewMessage(chatId, text)
	bot.sendMessage(msg)
}

func (bot *TelegramBotStruct) commands(update botInstance.Update) {
	command := update.Message.Command()
	switch command {
	case "cmd":
		msg := botInstance.NewMessage(update.Message.Chat.ID, "Planboy commands:")
		msg.ReplyMarkup = startMenu()
		msg.ParseMode = "Markdown"
		bot.sendMessage(msg)
	}
}

func (bot *TelegramBotStruct) sendMessage(msg botInstance.Chattable) {
	if _, err := bot.BotApi.Send(msg); err != nil {
		log.Panicf("Send message error: %v", err)
	}
}
