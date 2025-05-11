package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

// SecretCode определяет секретный код для команды /start для обычных пользователей
const SecretCode = "secretcode"

// AdminCode определяет секретный код для команды /start для администраторов
const AdminCode = "admincode"

// AuthorizedUsersFile определяет путь к JSON-файлу для хранения авторизованных пользователей
const AuthorizedUsersFile = "authorized_users.json"

// AuthorizedAdminsFile определяет путь к JSON-файлу для хранения авторизованных администраторов
const AuthorizedAdminsFile = "authorized_admins.json"

// loadAuthorizedList загружает список авторизованных ID из JSON-файла
func loadAuthorizedList(filePath string) ([]int64, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyList := []int64{}
		data, err := json.Marshal(emptyList)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, err
		}
		return emptyList, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var list []int64
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// saveAuthorizedList сохраняет список авторизованных ID в JSON-файл
func saveAuthorizedList(filePath string, list []int64) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// isAuthorizedInList проверяет, есть ли chatID в указанном списке
func isAuthorizedInList(chatID int64, list []int64) bool {
	for _, id := range list {
		if id == chatID {
			return true
		}
	}
	return false
}

// isAuthorized проверяет, есть ли chatID в одном из списков
func isAuthorized(chatID int64, lists ...[]int64) bool {
	for _, list := range lists {
		for _, id := range list {
			if id == chatID {
				return true
			}
		}
	}
	return false
}

// removeFromList удаляет chatID из списка
func removeFromList(chatID int64, list []int64) []int64 {
	for i, id := range list {
		if id == chatID {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func main() {
	// Получаем токен из переменной среды
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in environment variables")
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Загружаем списки авторизованных пользователей и администраторов
	authorizedUsers, err := loadAuthorizedList(AuthorizedUsersFile)
	if err != nil {
		log.Fatalf("Failed to load authorized users: %v", err)
	}
	authorizedAdmins, err := loadAuthorizedList(AuthorizedAdminsFile)
	if err != nil {
		log.Fatalf("Failed to load authorized admins: %v", err)
	}

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil { // Игнорируем не-сообщения
			continue
		}

		chatID := update.Message.Chat.ID

		// Обработка команды /start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			args := update.Message.CommandArguments()
			msg := tgbotapi.NewMessage(chatID, "")
			if args == SecretCode {
				// Проверяем, не авторизован ли пользователь уже
				if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
					// Добавляем пользователя в список обычных пользователей
					authorizedUsers = append(authorizedUsers, chatID)
					if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
						log.Printf("Failed to save authorized users: %v", err)
						msg.Text = "Ошибка при сохранении авторизации. Попробуйте позже."
					} else {
						msg.Text = "Вы подключились с секретным кодом! Теперь вы можете использовать бота как пользователь."
					}
				} else {
					msg.Text = "Вы уже авторизованы."
				}
			} else if args == AdminCode {
				// Проверяем, не авторизован ли пользователь как администратор
				if !isAuthorizedInList(chatID, authorizedAdmins) {
					// Если пользователь уже в списке обычных пользователей, переносим его
					if isAuthorizedInList(chatID, authorizedUsers) {
						authorizedUsers = removeFromList(chatID, authorizedUsers)
						if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
							log.Printf("Failed to save authorized users: %v", err)
							msg.Text = "Ошибка при обновлении списка пользователей. Попробуйте позже."
						} else {
							// Добавляем в список администраторов
							authorizedAdmins = append(authorizedAdmins, chatID)
							if err := saveAuthorizedList(AuthorizedAdminsFile, authorizedAdmins); err != nil {
								log.Printf("Failed to save authorized admins: %v", err)
								msg.Text = "Ошибка при сохранении авторизации администратора. Попробуйте позже."
							} else {
								msg.Text = "Вы повышены до администратора! Теперь вы можете использовать бота как администратор."
							}
						}
					} else {
						// Если пользователь не авторизован, просто добавляем в администраторы
						authorizedAdmins = append(authorizedAdmins, chatID)
						if err := saveAuthorizedList(AuthorizedAdminsFile, authorizedAdmins); err != nil {
							log.Printf("Failed to save authorized admins: %v", err)
							msg.Text = "Ошибка при сохранении авторизации администратора. Попробуйте позже."
						} else {
							msg.Text = "Вы подключились с кодом администратора! Теперь вы можете использовать бота как администратор."
						}
					}
				}
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			continue
		}

		// Проверяем, авторизован ли пользователь
		if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
			continue
		}

		// Разделяем механику для администраторов и обычных пользователей
		msg := tgbotapi.NewMessage(chatID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		if isAuthorizedInList(chatID, authorizedAdmins) {
			// Для администраторов отправляем их chatID
			msg.Text = fmt.Sprintf("Ваш ID: %d", chatID)
		} else {
			// Для обычных пользователей эхобот: повторяем текст сообщения
			msg.Text = update.Message.Text
		}
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

/*// Проверяем, авторизован ли пользователь
if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
	msg := tgbotapi.NewMessage(chatID, "Доступ запрещён. Используйте команду /start с правильным секретным кодом.")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
	continue
}*/
