package main

import (
	"log"
	"time"

	"tgBotUborochka/names"
	"tgBotUborochka/notification"
	"tgBotUborochka/token"

	"github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v3"
)

// пути откуда считывать токен бота и куда сохранять имена с ID для отправки уведомлений
const (
	//пути к файлам
	filePathNotifications = "D:\\GoLangProjects\\tgBotUborochka\\saves\\UsersNotifications.txt"
	filePathToken         = "D:\\GoLangProjects\\tgBotUborochka\\saves\\Token.txt"
	filePathNames         = "D:\\GoLangProjects\\tgBotUborochka\\saves\\Names.txt"

	//мин,час,день месяца,месяц,день недели(вс=0,пн=1,вт=2,...,сб=6)
	// * = любое
	// через запятую перечисление

	//время отравки уведомления о начале дежурства
	timeSendNotificationStartDuty = "0 18 * * 1"
	//время отравки уведомления о дежурстве вечером
	timeSendNotificationСleaning = "0 18 * * 2,4,0"
)

func check(err error, s string) {
	if err != nil {
		log.Println(s)
		log.Fatal(err)
	}
}

func WhoCleaningThisWeek(startTime time.Time, tgNames []string) string {
	dur := time.Since(startTime)
	numCleaner := int(dur.Hours()/24/7) % len(tgNames)

	return tgNames[numCleaner]
}

func main() {
	//время отсчета дежурства первого человека из списка
	//год,месяц,день,час,минута,секунда,наносек,
	firstWeek := time.Date(2024, 12, 23, 0, 1, 0, 0, time.Local)

	tgNames, err := names.GetNames(filePathNames)
	check(err, "Не считались имена")

	// получаем токен бота
	tgToken, err := token.GetToken(filePathToken)
	check(err, "Не считался токен")

	// главное
	menuMain := &tele.ReplyMarkup{ResizeKeyboard: true}
	// меню для уведомлений
	menuNotification := &tele.ReplyMarkup{ResizeKeyboard: true}

	// кнопки для главного меню
	btnWhoClean := menuMain.Text("🧹 Кто убирается на этой неделе? 🧹")
	btnNotification := menuMain.Text("🔔 Включить уведомления 🔔")
	btnNotificationOff := menuMain.Text("🔕 Отключить уведомления 🔕")

	// кнопки для выбора того, кто ты
	// чтобы можно было отправлять уведомления

	var rows []tele.Row // Создаем срез для хранения строк кнопок

	for _, name := range tgNames {
		// Создаем новую строку с одной кнопкой
		row := tele.Row{menuNotification.Text(name)}
		rows = append(rows, row) // Добавляем строку в общий срез
	}

	// Теперь вы можете использовать rows для создания меню
	menuNotification.Reply(rows...)

	menuMain.Reply(
		menuMain.Row(btnWhoClean),
		menuMain.Row(btnNotification),
		menuMain.Row(btnNotificationOff),
	)

	pref := tele.Settings{
		Token:  tgToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	check(err, "Не запустился бот")

	// КНОПОЧКИ

	b.Handle(&btnWhoClean, func(c tele.Context) error {
		resp := WhoCleaningThisWeek(firstWeek, tgNames)
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
		resp := WhoCleaningThisWeek(firstWeek, tgNames)
		return c.Send(resp)
	})

	// если написали имя из списка, сохраняем его ID, чтобы включить напоминание
	for _, name := range tgNames {
		b.Handle(name, func(c tele.Context) error {
			return notification.AddNotificationName(c, filePathNotifications, menuMain)
		})
	}

	//для отправки уведомлений в определенное время
	c := cron.New()
	// запус cron для отправки уведомлений
	c.Start()
	defer c.Stop()

	// отправка уведомлений
	_, err = c.AddFunc(timeSendNotificationStartDuty, func() {
		err = notification.SendNotification(
			"📣 Привет, %s.\nНаступила твоя неделя уборки на кухне.\n", // messageTemplate
			firstWeek, tgNames, filePathNotifications, b, WhoCleaningThisWeek,
		)
		check(err, "Не запустились уведомление о начале недели")
	})

	_, err = c.AddFunc(timeSendNotificationСleaning, func() {
		err = notification.SendNotification(
			"🕒 Привет, %s.\nНе забудь убраться сегодня на кухне.\nХорошего тебе вечера 😉.\n", // messageTemplate
			firstWeek, tgNames, filePathNotifications, b, WhoCleaningThisWeek,
		)
		check(err, "Не запустились уведомление о дежурстве")
	})

	// запускаем ботика
	b.Start()
	defer b.Stop()
}