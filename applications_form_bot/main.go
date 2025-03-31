package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	// Button options
	sewageButton   = "Хочу подати заявку щодо пробиття сміттєпроводу/не вивезеного сміття."
	lightButton    = "Хочу подати заявку щодо несправності світильника в місцях загального користування."
	plumbingButton = "Хочу подати заявку щодо проблем із каналізацією, опаленням чи водопостачанням."
	roofButton     = "Хочу подати заявку щодо протікання покрівлі чи герметних швів."
	otherButton    = "Інші запитання, пропозиції, зауваження."
	feedbackButton = "Ваші відгуки)))"

	// Main messages
	welcomeMessage = `Вітаю! На зв'язку менеджер житлових будинків КП НМР "ЖКО"👋🏻
Цей бот допоможе Вам і нам поліпшити комунікацію щодо надання послуг з управління багатоквартирними будинками. 

ВАЖЛИВО❗️
Якщо сталася аварійна ситуація, негайно телефонуйте за номером: 0 800 213 775      

Будь ласка, оберіть опцію, яка Вас цікавить👇`

	requestInfoText = `Для обробки Вашої заявки, будь ласка, надайте наступну інформацію:
- ПІБ
- Вулиця
- № будинку
- № квартири
- Контактний номер телефону
- Опишіть ситуацію детально

Введіть всю інформацію в одному повідомленні.`

	feedbackInfoText = `Залиште, будь ласка, свій відгук та контакти для зворотнього звʼязку у повідомленні`

	otherInfoText = `Опишіть Вашу пропозицію та залишіть контакт для зворотнього звʼязку`

	followUpText = `Якщо маєте пропозиції чи зауваження щодо якості надання послуг, можна також зателефонувати за номером 067 895 34 99, щоб ми могли врахувати та вчасно зреагувати.

Якщо і після звернення у Вас залишились зауваження чи пропозиції, зателефонуйте нашому керівнику за номером 098 464 68 63`
)

var groupID = int64(-1002546642948)

// var groupID = int64(-4672540477)

// User state tracking
var userStates = make(map[int64]string)

// Store information about groups the bot is in
type GroupInfo struct {
	ID       int64
	Title    string
	Username string
	Type     string
	JoinedAt time.Time
}

// Map to store groups the bot has seen
var botGroups = struct {
	sync.RWMutex
	groups map[int64]GroupInfo
}{groups: make(map[int64]GroupInfo)}

// Handles a button click by setting user state and sending appropriate info text
func handleButtonClick(bot *tgbotapi.BotAPI, userID int64, buttonText string) error {
	// Update the state with the selected button
	userStates[userID] = buttonText

	// Create message with appropriate text based on button type
	msg := tgbotapi.NewMessage(userID, "")

	// Use switch instead of if-else for better readability
	switch buttonText {
	case feedbackButton:
		msg.Text = feedbackInfoText
	case otherButton:
		msg.Text = otherInfoText
	default:
		msg.Text = requestInfoText
	}

	// Send the message
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending button response message: %v", err)
	}
	return err
}

// Add a group to the tracking list
func trackGroup(chat *tgbotapi.Chat) {
	if chat.Type == "group" || chat.Type == "supergroup" {
		botGroups.Lock()
		defer botGroups.Unlock()

		// Set the groupID to this chat's ID
		groupID = chat.ID
		log.Printf("Set groupID to: %d", groupID)

		// Only add if not already tracked
		if _, exists := botGroups.groups[chat.ID]; !exists {
			botGroups.groups[chat.ID] = GroupInfo{
				ID:       chat.ID,
				Title:    chat.Title,
				Username: chat.UserName,
				Type:     chat.Type,
				JoinedAt: time.Now(),
			}
			log.Printf("Bot added to new group: %s (ID: %d)", chat.Title, chat.ID)
		}
	}
}

func getMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(sewageButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(lightButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(plumbingButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(roofButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(otherButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(feedbackButton),
		),
	)
	keyboard.ResizeKeyboard = true
	return keyboard
}

func main() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: error opening .env file: %v", err)
		log.Println("Continuing without .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("APPLICATIONS_FORM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// Delete webhook (important!)
	_, err = bot.Request(tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateConfig.AllowedUpdates = []string{"message", "callback_query"}

	updates := bot.GetUpdatesChan(updateConfig)

	// Start HTTP server for Cloud Run
	go func() {
		port := os.Getenv("APPLICATIONS_FORM_BOT_PORT")
		if port == "" {
			port = "8080"
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Applications Form Bot is running!"))
		})
		log.Printf("Starting HTTP server on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Track group if message is from a group
		if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
			trackGroup(update.Message.Chat)
		}

		userID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(userID, "")

		// Handle commands
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				// Clear user state when they start again
				delete(userStates, userID)

				// Different behavior for groups vs private chats
				if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
					// Simplified message for groups
					msg.Text = "Будь ласка очікуйте на повідомлення"
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Error sending group welcome message: %v", err)
					}
				} else {
					// Send text message as fallback
					msg.Text = welcomeMessage
					msg.ReplyMarkup = getMainKeyboard()
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Error sending fallback welcome message: %v", err)
					}
				}
			}
			continue
		}

		// Handle request types or user input based on state
		if state, exists := userStates[userID]; exists {
			// Skip processing user state if in a group chat
			if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
				continue
			}

			// Check if user clicked another button while having a state
			switch update.Message.Text {
			case sewageButton, lightButton, plumbingButton, roofButton, otherButton, feedbackButton:
				// Handle button click using helper method
				if err := handleButtonClick(bot, userID, update.Message.Text); err != nil {
					log.Printf("Error handling button click: %v", err)
				}
				continue
			}

			// User is submitting request details
			// Forward message to admin if username is available
			if groupID != 0 {
				// Create forward message with context
				forwardText := "Нова заявка від користувача:\n"
				forwardText += "Тип заявки: " + state + "\n\n"
				forwardText += "Деталі заявки:\n" + update.Message.Text

				// Information about the sender
				if update.Message.From != nil {
					forwardText += "\n\nВідправник: "
					if update.Message.From.UserName != "" {
						forwardText += "@" + update.Message.From.UserName
						if update.Message.From.FirstName != "" {
							forwardText += ", " + update.Message.From.FirstName
						}
						if update.Message.From.LastName != "" {
							forwardText += " " + update.Message.From.LastName
						}
					}
				}

				// Try to forward to admin by username
				adminMsg := tgbotapi.NewMessage(groupID, forwardText)
				if _, err := bot.Send(adminMsg); err != nil {
					log.Printf("Error forwarding message to group: %v", err)
				}
			}

			// Reset user state and send follow-up message
			delete(userStates, userID)
			msg.Text = followUpText
			msg.ReplyMarkup = getMainKeyboard()
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending follow-up message: %v", err)
			}
			continue
		}

		// Handle button clicks to set user state
		switch update.Message.Text {
		case sewageButton, lightButton, plumbingButton, roofButton, otherButton, feedbackButton:
			// Skip menu button handling and message sending in group chats
			if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
				log.Printf("Ignoring button click in group chat: %s", update.Message.Chat.Title)
				continue
			}

			// Skip if user already has a state (handled above)
			if _, exists := userStates[userID]; exists {
				continue
			}

			// Handle button click using helper method
			if err := handleButtonClick(bot, userID, update.Message.Text); err != nil {
				log.Printf("Error handling button click: %v", err)
			}
			continue
		default:
			// If we don't recognize the message and no state, ask to use the menu
			// Only respond with menu prompt in private chats
			if _, exists := userStates[userID]; !exists && update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup" {
				msg.Text = "Будь ласка, використовуйте кнопки меню для вибору типу заявки."
				msg.ReplyMarkup = getMainKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending menu prompt: %v", err)
				}
			}
		}
	}
}
