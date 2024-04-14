package telegramnotionbot

import (
	"fmt"
	"log"
	"os"

	botInstance "github.com/crocone/tg-bot"
	"github.com/terrnitllc/planboy/notion"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
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

type TelegramBotStruct struct {
	BotApi *botInstance.BotAPI
	Notion *notion.NotionApi
	Token  string
}

func (bot *TelegramBotStruct) startBot() {
	bot, err := botInstance.NewBotAPI(bot.Token)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	u := botInstance.NewUpdate(0)
	u.Timeout = 60
	updates := bot.BotAPI.GetUpdatesChan(u)
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

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	telegramApiCode, ok := os.LookupEnv("TG_API_BOT_TOKEN")
	if !ok {
		log.Fatalf("TG_API_BOT_TOKEN is not set")
	}

	tgBot := TelegramBotStruct{BotAPI: nil, Notion: nil, Token: telegramApiCode}

	// Notion client
	notionApiSecret, ok := os.LookupEnv("NOTION_API_SECRET")
	if !ok {
		log.Fatalf("NOTION_API_SECRET is not set")
	}

	notionDbKey, ok := os.LookupEnv("NOTION_DB_KEY")
	if !ok {
		log.Fatalf("NOTION_DB_KEY is not set")
	}

	notionClient, ok := notionapi.NewClient(notionApiSecret, notionDbKey)
	if !ok {
		log.Fatalf("Notion client creation failed")
	}

	tgBot.Notion = notionClient
	tgBot.startBot()

	log.Print("Bot has been started !")
}

func callbacks(update botInstance.Update) {
	data := update.CallbackQuery.Data
	chatId := update.CallbackQuery.From.ID
	firstName := update.CallbackQuery.From.FirstName
	lastName := update.CallbackQuery.From.LastName
	var text string
	switch data {
	case "get_page":
		text = fmt.Sprintf("Привет %v %v", firstName, lastName)
	default:
		text = "Неизвестная команда"
	}
	msg := botInstance.NewMessage(chatId, text)
	sendMessage(msg)
}

func commands(update botInstance.Update) {
	command := update.Message.Command()
	switch command {
	case "commands":
		msg := botInstance.NewMessage(update.Message.Chat.ID, "Planboy commands:")
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
