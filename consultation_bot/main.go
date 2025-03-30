package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	aboutButton        = "1. Що таке перерахунок плати за послугу з вивезення побутових відходів?"
	instructionButton  = "2. Покрокова інструкція у випадку перебування за кордоном"
	consultationButton = "3. Особиста консультація з юрисконсультом"
	otherCasesButton   = "4. Я перебуваю не за кордоном"
	welcomeMessage     = `Вітаю! На зв'язку юрисконсульт КП НМР "ЖКО"👋🏻
Цей бот допоможе Вам розібратися з питаннями, які виникають при написанні заяви про перерахунок плати за послугу з вивезення побутових відходів у зв'язку з перебуванням за кордоном📃
Будь ласка, уточніть, яке питання Вас цікавить?👇`
	welcomeImagePath     = "imgs/welcome.JPEG"
	aboutImagePath       = "imgs/about.JPEG"
	instructionImagePath = "imgs/instruction.JPEG" // Add your instruction image here
	aboutText            = `1. Відповідно до Закону України «Про житлово-комунальні послуги» (пп. 6 п. 1 ст. 7), споживач звільняється від оплати комунальних послуг у разі їх невикористання під час тимчасової відсутності в житловому приміщенні понад 30 днів. Для цього необхідно надати документальне підтвердження. 
📌Тобто, якщо ви не проживаєте у своєму помешканні більше ніж 30 днів, ви маєте право не сплачувати за вивезення сміття. Тимчасова відсутність - це до 6-и місяців`
	instructionText = `Спершу, важливі нюанси:
📍Заяву необхідно оформляти окремо для кожної особи.
📍У разі потреби отримання даних про перетин кордону неповнолітньою особою, заяву подає один із батьків.
📍Заяву пишемо вручну на аркуші А4.
📍Для підтвердження Вашої особистості необхідно надати ксерокопію паспорту - українського або закордонного; перевага надається українському.
📍КП НМР "ЖКО" здійснює перерахунок оплати виключно за період, що не перевищує 6 місяців.`
	otherCasesTextFirst = `Окрім перебування за кордоном, існують й інші підстави для перерахунку. Ви можете ознайомитися з ними, переглянувши Положення:

Відповідно до Рішення Виконавчого  комітету НМР Хмельницької області від 25.07.2024 року №207/2024 «Про внесення змін до рішення виконавчого комітету Нетішинської міської ради від 27 лютого 2020 року Лº 99/2020 «Про затвердження Положення про перерахунок плати за послугу з поводження з побутовими відходами за період тимчасової відсутності споживача та/або членів його сімʼї»»:

1. При тимчасовій відсутності квартиронаймача, власника (далі
Споживач) та/або членів його сімʼї за місцем постійного проживання безперервно більше, ніж 30 календарних днів, підприємство - виконавець послуг згідно з письмовою заявою Споживача та підтверджуючих документів про його тимчасову відсутність та/або членів його сімʼї не нараховує плату на послуги з поводження з побутовими відходами за термін їх відсутності.
2. Підставою для звільнення від нарахування плати за послугу з поводження з побутовими відходами с тимчасова відсутність (не проживання) Споживача
членів його сімʼї в квартирі безперервно 30 календарних днів.
3. Для звільнення від нарахування плати за послугу поводження побутовими відходами Споживач або члени його сімʼї повинні надати заяву та документально підтвердити факт тимчасової відсутності за місцем проживання, а саме:
📍особи, відсутні за місцем рестрації та проживають за іншою адресою надають:
- довідку/акт про підтвердження перебування Споживача та/або членів сімʼї у іншій місцевості завірену Управителем, ОСББ, житлово-експлуатаційною організацією, ЖБК, на території якого розташовано будинок, або органом місцевого самоврядування із вказівкою періоду перебування;
- довідку, яка підтверджує проживання в гуртожитках, готелях та ін. із вказівкою періоду проживання;
- документ з тимчасовою пропискою в іншій місцевості;
📍особи, які мешкають за межами України, подають витяг із Бази даних про перетинання
державного кордону України (приклад заяви та інформація розміщені на сайті Державної прикордонної служби - https://dpsu.gov.ua/ua/Zayava-stosovno-peretinannya-osoboyu-derzhavnogo-kordonu);
📍особи, які знаходяться на стаціонарному лікуванні:
- довідки з місця стаціонарного лікування, засвідченої підписом лікаря та печаткою лікарняної установи, із зазначенням періоду лікування; 
📍особи, які перебувають на навчанні, надають довідку з місця навчання, засвідчену підписом керівника з печаткою установи, із зазначенням періоду навчання, населеного пункту, в якому особа проходила навчання, та форми навчання (денна (стаціонар), у випадку навчання за кордоном надається довідка або інший документ про таке навчання з офіційним перекладом на державну мову. Вказані довідки поновлюються щорічно на початку навчального року.
📍особи, які проходять службу в Збройних силах України надають довідку з місця проходження служби, засвідченої підписом керівника та печаткою установи, із зазначенням періоду служби;
📍особи, які відбувають покарання в місцях позбавлення волі, за наявності таких документів:
- довідки з місця позбавлення волі, засвідченої підписом керівника та печаткою установи, із зазначенням періоду позбавлення волі;
- з моменту звільнення з місць позбавлення волі - копію відповідної довідки про звільнення із зазначенням терміну перебування, засвідченої підписом керівника та печаткою установи;
📍особи, працевлаштовані в іншому населеному пункті належним чином, подають оформлену довідку з місця працевлаштування, або трудового договору, із зазначенням населеного пункту (адреси) роботодавця і періоду працевлаштування, графіком роботи, за умови, що місце роботи розташоване на відстані не менше 80 км від місця ресстрації.
📍для осіб, визначених у встановленому порядку безвісно відсутніми або місце перебування яких невідомо, копії рішення суду або інший відповідний
документ.
`
	otherCasesTextSecond = `❗️Усі документи повинні бути складені державною мовою та подаються споживачем особисто або на електронну пошту виконавця послуги - jko_netish@ukr.net. 
Документи, складені іноземною мовою - подаються з офіційним перекладом на українську мову.

4. Вище перелічені документи Споживач та/або члени його сімʼї повинні подати до підприємств-виконавців послуги протягом 14 днів з момситу повернення до постійного місця проживання.
5. Заява від Споживача або членів його сімʼї для перерахунку плати за послуги з поводження з побутовими відходами за попередній період приймається на термін вказаний в заяві, але не більше як на шість календарних місяців, що передують даті звернення.
6. Споживач та/або члени його сімʼї напередодні вибуття з місця
проживання, надають письмову заяву про не нарахування плати за послугу на імʼя керівника підприємства - виконавця послуг. У заяві вказується термін, на який не проводиться нарахування, але не більше шести календарних місяців з дня подачі заяви.
По закінченню вказаного терміну Споживач протягом 14 календарних днів
З дня повернення надас до Підприємства підтверджуючі документи про тимчасову відсутність з підтвердженням Терміну відсутності.
7. Право на звільнення від нарахувань плати за послугу з поводження з побутовими відходами на період більше шести місяців мають споживачі у випадку:
- призову до Збройних сил 
- оголошення в розшук
- затримання органом внутрішніх справ;
- навчання на денній формі в учбових закладах, які знаходяться в іншій місцевості;
- перебування в лікувальних закладах;
- визнання особи безвісно відсутньою або такою, місце перебування якої не відомо.
У вказаних випадках нарахування припиняється на період, вказаний у
наданих Споживачем документах.
8. Якщо документи, які підтверджують відсутність Споживача та/або членів його сімʼї, надані не в повному обсязі або не оформлені належним чином (наявність штампу, дати ресстрації, підписів відповідальних осіб та ін.), плата за послугу з поводження з побутовими відходами нараховується в повному обсязі.
9. Поновлення нарахування плати за послуги відбувається автоматично. В разі дострокового повернення Споживача та/або членів його сімʼї до місця його постійного проживання, Споживач повідомляє про це Виконавця послуги протягом 14 календарних днів з моменту повернення.
У разі подання такої інформації після 14-денного терміну з моменту повернення споживача до місця проживання - перерахунок минулих періодів не здійснюється.
10. У разі неврегульованої ситуації, рішення може бути прийнято індивідуально керівником КП НМР «ЖКО».
11. Інформація про відсутність споживачів за місцем проживання розголошенню не підлягає.`
	instructionStepOneText = `Перший крок👇🏻
Пишете заяву до міграційної служби. На аркуші формату А4 у верхньому правому куті зазначається: "Головний центр обробки спеціальної інформації Державної прикордонної служби". Нижче вписуються ваші дані: прізвище, ім'я, по батькові, громадянство, дата народження, номер паспорта (який додається до заяви), орган, що видав паспорт, а також повна адреса (область, район, місто, вулиця, номер будинку та квартири).`
	instructionStepTwoText = `Другий крок👇🏻
У центрі заяви зазначаємо слово "Заява". Далі пишемо: "Прошу надати інформацію про перетинання мною державного кордону України за період з _________ (тут пишемо, коли востаннє перетинали кордон України) по _________."`
	instructionStepThreeText = `Третій крок👇🏻
Далі зазначаємо: "Відповідь прошу надіслати на електронну адресу _____ або на поштову адресу _____." (можете вказати поруч зі своєю ел.пошту так і ел.пошту КП НМР "ЖКО" (lawyerjko@gmail.com), аби відповідь прийшла безпосередньо до нас)`
	instructionStepFourText = `Четвертий крок👇🏻
Далі пишемо: "Додаток на ___ арк.", де вказуємо кількість аркушів ксерокопії паспорта, що додаються до заяви.`
	instructionStepFiveText = `П'ятий крок👇🏻
Сфотографуйте та підпишіть заяву, разом із ксерокопіями паспорта надішліть на електронну адресу Державної прикордонної служби України: zvernennia@dpsu.gov.ua.`
	instructionStepSixText = `Шостий крок 👇🏻
Готово! Тепер мешканцю квартири, за адресою, де ви тимчасово не проживаєте, необхідно відвідати КП НМР "ЖКО" для подання додаткової заяви з метою проведення перерахунку безпосередньо для квартири.  
Якщо це вже виконано – чудово! Очікуємо відповіді від ДПС!

‼️Якщо споживач мешкає у багатоквартирному будинку, в якому створено ОСББ та обрано модель договору колективний споживач, заява до КП НМР "ЖКО" про зміну кількості споживачів у будинку подається головою ОСББ. 
Тому у разі зміни кількості споживачів у Вашому помешканні, підтверджуючі документи необхідно подати голові Вашого ОСББ.

Якщо у Вас залишились ще якісь запитання, оберіть опцію "Особиста консультація з юрисконсультом" або "Закінчити консультацію".`
)

var consultationText = `Для отримання особистої консультації зверніться до юрисконсультанта - @bulliia, та опишіть їй своє питання.`

func getMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(aboutButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(instructionButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(consultationButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(otherCasesButton),
		),
	)
	keyboard.ResizeKeyboard = true
	return keyboard
}

func getInstructionKeyboard(showNextStep bool, step string) tgbotapi.InlineKeyboardMarkup {
	if !showNextStep {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Перейти до покрокової інструкції", "step_one"),
			),
		)
	}

	var buttons []tgbotapi.InlineKeyboardButton

	// Add "Previous" button for all steps except first
	if step != "step_one" {
		log.Printf("step: %s", step)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Попередній крок", getPreviousStep(step)))
	}

	// Add "Next" button for all steps except last
	if step != "step_six" {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Наступний крок", getNextStep(step)))
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
}

func getPreviousStep(currentStep string) string {
	switch currentStep {
	case "step_two":
		return "step_one"
	case "step_three":
		return "step_two"
	case "step_four":
		return "step_three"
	case "step_five":
		return "step_four"
	case "step_six":
		return "step_five"
	default:
		return "step_one"
	}
}

func getNextStep(currentStep string) string {
	switch currentStep {
	case "step_one":
		return "step_two"
	case "step_two":
		return "step_three"
	case "step_three":
		return "step_four"
	case "step_four":
		return "step_five"
	case "step_five":
		return "step_six"
	default:
		return "step_six"
	}
}

func getStepText(step string) string {
	switch step {
	case "step_one":
		return instructionStepOneText
	case "step_two":
		return instructionStepTwoText
	case "step_three":
		return instructionStepThreeText
	case "step_four":
		return instructionStepFourText
	case "step_five":
		return instructionStepFiveText
	case "step_six":
		return instructionStepSixText
	default:
		return ""
	}
}

func main() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error opening .env file: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("CONSULTATION_BOT_TOKEN"))
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
		port := os.Getenv("CONSULTATION_BOT_PORT")
		if port == "" {
			port = "8080"
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Bot is running!"))
		})
		log.Printf("Starting HTTP server on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	for update := range updates {

		// Handle callback queries first
		if update.CallbackQuery != nil {
			log.Printf("Processing callback query: %s", update.CallbackQuery.Data)

			// Answer the callback query first
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Printf("Error answering callback: %v", err)
			}

			stepText := getStepText(update.CallbackQuery.Data)
			if stepText != "" {
				// Delete previous message
				deleteMsg := tgbotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Printf("Error deleting previous message: %v", err)
				}

				// Send new message
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, stepText)
				markup := getInstructionKeyboard(true, update.CallbackQuery.Data)
				msg.ReplyMarkup = &markup
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending step message: %v", err)
				}
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Show keyboard for /start command
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			// Send welcome image with caption
			photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(welcomeImagePath))
			photo.Caption = welcomeMessage
			photo.ReplyMarkup = getMainKeyboard()
			if _, err := bot.Send(photo); err != nil {
				log.Printf("Error sending welcome message with image: %v", err)
			}
			continue
		}

		// Handle button clicks
		switch update.Message.Text {
		case aboutButton:
			photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(aboutImagePath))
			photo.Caption = aboutText
			photo.ReplyMarkup = getMainKeyboard()
			if _, err := bot.Send(photo); err != nil {
				log.Printf("Error sending about message with image: %v", err)
			}
			continue
		case instructionButton:
			photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(instructionImagePath))
			photo.Caption = instructionText
			photo.ReplyMarkup = getInstructionKeyboard(false, "step_one")
			if _, err := bot.Send(photo); err != nil {
				log.Printf("Error sending instruction message with image: %v", err)
			}
			continue
		case consultationButton:
			msg.Text = consultationText
		case otherCasesButton:
			// Send text message first
			msg.Text = otherCasesTextFirst
			msg.ReplyMarkup = getMainKeyboard()
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending other cases message: %v", err)
			}

			msg.Text = otherCasesTextSecond
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending other cases message: %v", err)
			}
			continue
		default:
			msg.Text = "Будь ласка, використовуйте кнопки меню для навігації."
		}

		msg.ReplyMarkup = getMainKeyboard()
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
