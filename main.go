package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

var host = "localhost"
var port = 5432
var user = os.Getenv("USER")
var password = 1234
var dbname = "postgres"
var sslmode = ""

var dbInfo = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
	host,
	port,
	user,
	password,
	dbname,
	sslmode)

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

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
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USER_ID INT, USERNAME TEXT, CHAT_ID INT, MESSAGE TEXT, ANSWER TEXT);`); err != nil {
		return err
	}

	return nil
}

//Собираем данные полученные ботом
func collectData(userid int, username string, chatid int, message string, answer []string) error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Конвертируем срез с ответом в строку
	answ := strings.Join(answer, ", ")

	//Создаем SQL запрос
	data := `INSERT INTO users(user_id,username, chat_id, message, answer) VALUES($1, $2, $3, $4, $5);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, userid, `@`+username, chatid, message, answ); err != nil {
		return err
	}

	return nil
}

// func isCW()bool {

// }

const idCW3 = 265204902

var sun = []string{"CAADAgADOgAD5R-VAnqF-5FEu7a2Ag",
	"CAADBAADCAIAAkb3JwABgT3OjSZrl3gC",
	"CAADBAADHwIAAkb3JwACZMCHHkNgmgI",
	"CAADBAADDwIAAkb3JwABnscPo7pyJQYC",
	"CAADBAAD-QEAAkb3JwABrrmzBNftMI8C"}

func main() {
	bot, err := tgbotapi.NewBotAPI("430629496:AAHDBvxHimRzeURldxAz_4v8pp4bKzoeH8s")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	log.Printf("createTable: %v", createTable())

	updates, err := bot.GetUpdatesChan(u)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		reply := ""
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if reflect.TypeOf(update.Message.Text).Kind() != reflect.String && update.Message.Text == "" {
			log.Println("Не текст. Игнор")
			continue
		}
		if update.Message.Text == "!sun" {
			bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAADAgADOgAD5R-VAnqF-5FEu7a2Ag"))
			continue
		}

		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я Солер из Асторы, воин Света, верный слуга короля жёлтого замка Попс Маэллард."
			//continue
		case "sun":
			rand.Seed(time.Now().Unix())
			bot.Send(tgbotapi.NewStickerShare(update.Message.Chat.ID, sun[rand.Intn(len(sun))]))
			continue
		case "info":
			if update.Message.ReplyToMessage == nil {
				log.Printf("%v Не является ответом на сообщение. Игнор", update.Message.ReplyToMessage)
				continue
			}
			t := time.Now()
			reply = fmt.Sprintf("Nickname: %v \nID: %v \nMessageID: %v \nTimeMsg:  %v \nTimeNow: %v \nDateFrwd: %v \nIsBot? %v \n%v",
				update.Message.ReplyToMessage.From.UserName,
				update.Message.ReplyToMessage.From.ID,
				update.Message.ReplyToMessage.MessageID,
				update.Message.ReplyToMessage.Time().Format(time.RFC1123Z),
				t.Format(time.RFC1123Z),
				update.Message.ReplyToMessage.ForwardDate,
				update.Message.ReplyToMessage.Entities)
		case "reg":
			if update.Message.ReplyToMessage == nil {
				log.Printf("%v Не является ответом на сообщение. Игнор", update.Message.ReplyToMessage)
				continue
			}
			err = collectData(update.Message.ReplyToMessage.From.ID, update.Message.ReplyToMessage.From.UserName, update.Message.From.ID, update.Message.Text, []string{update.Message.ReplyToMessage.Text})
			if err != nil {
				log.Println("collect data: err")
			}
			continue
		case "fwrdid":
			reply = fmt.Sprintln(update.Message.ReplyToMessage.ForwardFrom.ID)
		default:
			log.Printf("Не знаю такой команды %s", update.Message.Text)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
