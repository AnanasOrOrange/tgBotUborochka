package main

import (
	"bufio"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"tgBotUborochka/notification"
	"tgBotUborochka/token"

	tele "gopkg.in/telebot.v3"
)

// пути откуда считывать токен бота и куда сохранять имена с ID для отправки уведомлений
const (
	filePathNotifications = "D:\\GoLangProjects\\tgBotUborochka\\saves\\UsersNotifications.txt"
	filePathToken         = "D:\\GoLangProjects\\tgBotUborochka\\saves\\Token.txt"
)

func check(err error, s string) {
	if err != nil {
		log.Println(s)
		log.Fatal(err)
	}
}

func main() {

	// Universal markup builders.
	// главное
	menuMain := &tele.ReplyMarkup{ResizeKeyboard: true}
	// меню для уведомлений
	menuNotification := &tele.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.
	// кнопки для главного меню
	//btnHelp := menuMain.Text("ℹ Help")
	btnWhoClean := menuMain.Text("🧹 Кто убирается на этой неделе? 🧹")
	btnNotification := menuMain.Text("🔔 Включить уведомления 🔔")
	btnNotificationOff := menuMain.Text("🔕 Отключить уведомления 🔕")

	// кнопки для выбора того, кто ты
	// чтобы можно было отправлять уведомления
	btn0 := menuNotification.Text("Суров Алексей")
	btn1 := menuNotification.Text("Дудорова Ева")
	btn2 := menuNotification.Text("Дураев Дмитрий")
	btn3 := menuNotification.Text("Подмарев Никита")
	btn4 := menuNotification.Text("Дудоров Никита")
	btn5 := menuNotification.Text("Топазов Сергей")
	btn6 := menuNotification.Text("Столяров Дмитрий")

	menuMain.Reply(
		//menu.Row(btnHelp),
		menuMain.Row(btnWhoClean),
		menuMain.Row(btnNotification),
		menuMain.Row(btnNotificationOff),
	)

	menuNotification.Reply(
		menuMain.Row(btn0),
		menuMain.Row(btn1),
		menuMain.Row(btn2),
		menuMain.Row(btn3),
		menuMain.Row(btn4),
		menuMain.Row(btn5),
		menuMain.Row(btn6),
	)

	firstWeek := time.Date(2024, 12, 16, 0, 1, 0, 0, time.Local)

	// получаем токен бота
	tgToken, err := token.GetToken(filePathToken)
	check(err, "Не считался токен")

	pref := tele.Settings{
		Token:  tgToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	check(err, "Не запустился бот")

	// КНОПОЧКИ

	// On reply button pressed (message)
	//b.Handle(&btnHelp, func(c tele.Context) error {
	//	return c.Edit("Here is some help: ...")
	//})

	b.Handle(&btnWhoClean, func(c tele.Context) error {
		resp := WhoCleaningThisWeek(firstWeek)
		return c.Send(resp)
	})
	b.Handle(&btnNotification, func(c tele.Context) error {
		return c.Send("Выберете себя из списка. Чтобы вы получали уведомления.", menuNotification)
	})
	b.Handle(&btnNotificationOff, func(c tele.Context) error {
		stringID := c.Recipient().Recipient()
		deleted, err := notification.ExcludeLines(filePathNotifications, stringID)
		if err != nil {
			return err
		}
		if deleted {
			return c.Send("Уведомления отключены")
		}
		return c.Send("У вас нет уведомлений")
	})

	// КОМАНДЫ

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Открыто главное меню", menuMain)
	})
	b.Handle("меню", func(c tele.Context) error {
		return c.Send("Открыто главное меню", menuMain)
	})
	b.Handle("кто", func(c tele.Context) error {
		resp := WhoCleaningThisWeek(firstWeek)
		return c.Send(resp)
	})

	b.Handle("Суров Алексей", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})
	b.Handle("Дудорова Ева", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})
	b.Handle("Дураев Дмитрий", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})
	b.Handle("Подмарев Никита", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})
	b.Handle("Дудоров Никита", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})
	b.Handle("Топазов Сергей", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})
	b.Handle("Столяров Дмитрий", func(c tele.Context) error {
		ID := c.Recipient().Recipient()
		_, err := notification.ExcludeLines(filePathNotifications, ID)
		if err != nil {
			return err
		}

		stingID := c.Recipient().Recipient()
		newLine := fmt.Sprintf("%s %s\n", c.Text(), stingID)

		err = notification.AppendLine(filePathNotifications, newLine)
		if err != nil {
			return err
		}
		return c.Send("✅ Сохранено ✅", menuMain)
	})

	//для отправки уведомлений в определенное время
	c := cron.New()
	// запус cron для отправки уведомлений
	c.Start()
	defer c.Stop()

	// отправка уведомлений
	//_, err = c.AddFunc("0 8 * * 1", func() {
	_, err = c.AddFunc("0 18 * * 1", func() {
		whoCleaning := WhoCleaningThisWeek(firstWeek)

		file, err := os.Open(filePathNotifications)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, whoCleaning) {
				stringID := strings.TrimFunc(line, func(r rune) bool {
					return !unicode.IsNumber(r)
				})
				ID, err := strconv.Atoi(stringID)
				if err != nil {
					log.Println(err)
				}
				chat, err := b.ChatByID(int64(ID))
				if err != nil {
					log.Println(err)
				}
				message := fmt.Sprintf(
					"📣 Привет, %s.\nНаступила твоя неделя уборки на кухне.\n",
					whoCleaning,
				)

				_, err = b.Send(chat, message)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	})
	check(err, "Не запустились уведомление о начале недели")

	_, err = c.AddFunc("0 18 * * 2,4,0", func() {
		whoCleaning := WhoCleaningThisWeek(firstWeek)

		file, err := os.Open(filePathNotifications)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, whoCleaning) {
				stringID := strings.TrimFunc(line, func(r rune) bool {
					return !unicode.IsNumber(r)
				})
				ID, err := strconv.Atoi(stringID)
				if err != nil {
					log.Println(err)
				}
				chat, err := b.ChatByID(int64(ID))
				if err != nil {
					log.Println(err)
				}

				message := fmt.Sprintf(
					"🕒 Привет, %s.\nНе забудь убраться сегодня на кухне.\nХорошего тебе вечера 😉.\n",
					whoCleaning,
				)

				_, err = b.Send(chat, message)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	})
	check(err, "Не запустились уведомление о дежурстве")

	// запускаем ботика
	b.Start()
	defer b.Stop()
}

func WhoCleaningThisWeek(startTime time.Time) string {
	names := []string{
		"Суров Алексей",
		"Дудорова Ева",
		"Дураев Дмитрий",
		"Подмарев Никита",
		"Дудоров Никита",
		"Топазов Сергей",
		//"Столяров Дмитрий",
	}
	dur := time.Since(startTime)
	numCleaner := int(dur.Hours()/24/7) % len(names)

	return names[numCleaner]
}
