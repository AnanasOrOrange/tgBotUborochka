package notification

import (
	"bufio"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ExcludeLines(filePath string, needToDelete string) (bool, error) {
	lineFound := false
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, needToDelete) {
			lineFound = true
			continue
		}
		lines = append(lines, text)
	}

	if err := scanner.Err(); err != nil {
		return lineFound, err
	}

	err = ioutil.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	return lineFound, err
}

func AppendLine(filePath string, needToAppend string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprint("\n", needToAppend))
	return err
}

// Общая функция для обработки всех имен
func AddNotificationName(c tele.Context, filePathNotifications string, menuMain *tele.ReplyMarkup) error {
	ID := c.Recipient().Recipient()
	_, err := ExcludeLines(filePathNotifications, ID)
	if err != nil {
		return err
	}

	newLine := fmt.Sprintf("%s %s\n", c.Text(), ID)
	err = AppendLine(filePathNotifications, newLine)
	if err != nil {
		return err
	}

	return c.Send("✅ Сохранено ✅", menuMain)
}

// функция отправки уведомления
func SendNotification(messageTemplate string, firstWeek time.Time, tgNames []string, filePathNotifications string, b *tele.Bot, WhoCleaningThisWeek func(time.Time, []string) string) error {
	whoCleaning := WhoCleaningThisWeek(firstWeek, tgNames)

	file, err := os.Open(filePathNotifications)
	if err != nil {
		return err
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
				return err
			}

			chat, err := b.ChatByID(int64(ID))
			if err != nil {
				return err
			}

			message := fmt.Sprintf(messageTemplate, whoCleaning)
			_, err = b.Send(chat, message)
			if err != nil {
				return err
			}
		}
	}
	return err
}
