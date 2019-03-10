package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	host,
	port,
	user,
	password,
	dbname,
	sslmode)

// type Korovan string {
// 	Chat_ID
// }

// не нужное
// type SearchResults struct {
// 	ready   bool
// 	Query   string
// 	Results []Result
// }

// type Result struct {
// 	Name, Description, URL string
// }

// func (sr *SearchResults) UnmarshalJSON(bs []byte) error {
// 	array := []interface{}{}
// 	if err := json.Unmarshal(bs, &array); err != nil {
// 		return err
// 	}
// 	sr.Query = array[0].(string)
// 	for i := range array[1].([]interface{}) {
// 		sr.Results = append(sr.Results, Result{
// 			array[1].([]interface{})[i].(string),
// 			array[2].([]interface{})[i].(string),
// 			array[3].([]interface{})[i].(string),
// 		})
// 	}
// 	return nil
// }

//Создаем таблицу users в БД при подключении к ней
func createTable() error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Создаем таблицу users
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USERID TEXT, MESSAGE TEXT, TIME TEXT);`); err != nil {
		return err
	}

	return nil
}

//Собираем данные полученные ботом
func collectData(userID int, message string, time string) error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// //Конвертируем срез с ответом в строку
	// answ := strings.Join(answer, ", ")

	//Создаем SQL запрос
	data := `INSERT INTO users(userID, message, time) VALUES($1, $2, $3);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, userID, message, time); err != nil {
		return err
	}

	return nil
}

// func isCW()bool {

// }

const idCW3 = 265204902

func main() {
	bot, err := tgbotapi.NewBotAPI("430629496:AAHDBvxHimRzeURldxAz_4v8pp4bKzoeH8s")
	if err != nil {
		log.Panic(err)
	}
	tolerance := 1
	// buf1 := ""

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	log.Printf("createTable: %v", createTable())

	updates, err := bot.GetUpdatesChan(u)
	reply := ""

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "!sun" {
			bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
			continue
		}

		if strings.Contains(update.Message.Text, "!sun") {
			bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
			continue
		}
		if strings.Contains(update.Message.Text, "Ты задержал") {
			log.Printf("collectData: %v %v %v", update.Message.Time().Format(time.RFC1123Z), update.Message.From.ID, update.Message.Text)
			if err := collectData(
				update.Message.From.ID,
				update.Message.Text,
				update.Message.ReplyToMessage.Time().Format(time.RFC1123Z)); err != nil {
				log.Printf("collectData: %v", err)
				reply = "Схоронил корован для сравнения"
				continue
			}
			continue
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				reply = "Привет. Я Солер из Асторы, воин Света, верный слуга короля жёлтого замка Попс Маэллард."
				continue
			case "sun":
				bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
				continue
			case "info":
				log.Printf("ReplyMsg: %s\nMsg: %s", update.Message.ReplyToMessage.Text, update.Message.Text)
				t := time.Now()
				reply = fmt.Sprintf("@%v \nID: %v \nMessageID: %v \nTimeMsg:  %v \nTimeNow: %v",
					update.Message.ReplyToMessage.From.UserName,
					update.Message.ReplyToMessage.From.ID,
					update.Message.ReplyToMessage.MessageID,
					update.Message.ReplyToMessage.Time().Format(time.RFC1123Z),
					t.Format(time.RFC1123Z))
			case "gentle":
				tolerance = 1
				reply = fmt.Sprint("Вежливый режим включен")
			case "hui":
				tolerance = 0
				reply = fmt.Sprint("Вежливый режим отключен")
			case "bunt":
				reply = fmt.Sprint("*БУНД!!1*")
			case "getupdates":
				log.Printf("%v", updates)
			default:
				log.Printf("Не знаю такой команды %s", update.Message.Text)
				if tolerance == 1 {
					continue
				}
				reply = fmt.Sprint("Хуле доебался")
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "markdown"

		bot.Send(msg)
	}
}
