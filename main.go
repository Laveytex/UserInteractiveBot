package main

/*// Добавляем ваш chat_id в authorizedUsers для тестов, если список пуст
if len(authorizedUsers) == 0 && !isAuthorizedInList(1603504912, authorizedAdmins) {
	authorizedUsers = append(authorizedUsers, 1603504912)
	if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
		log.Printf("Failed to save authorized users: %v", err)
	} else {
		log.Printf("Added chat_id 1603504912 to authorized users")
	}
}*/

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"
)

const (
	AuthorizedUsersFile  = "authorized_users.json"
	AuthorizedAdminsFile = "authorized_admins.json"
	NicknamesFile        = "nicknames.json"
	PostsFile            = "posts.json"
	SecretCode           = "your_secret_code"
	AdminCode            = "your_admin_code"
	MaxCaptionLength     = 1024 // Максимальная длина подписи в Telegram
	MaxMediaCount        = 10
)

type UserNickname struct {
	ChatID   int64  `json:"chat_id"`
	Nickname string `json:"nickname"`
}

// Post хранит информацию о посте
type Post struct {
	ID          string    `json:"id"`
	AdminChatID int64     `json:"admin_chat_id"`
	Text        string    `json:"text"`
	PhotoIDs    []string  `json:"photo_ids"`
	PublishTime time.Time `json:"publish_time"`
	Published   bool      `json:"published"`
}

// CreatingPost хранит временное состояние создаваемого поста
type CreatingPost struct {
	AdminChatID int64
	Text        string
	PhotoIDs    []string
	State       AdminState
}

// AdminState определяет состояние админа при создании поста
type AdminState string

const (
	StateIdle        AdminState = "idle"
	StateAddingText  AdminState = "adding_text"
	StateAddingPhoto AdminState = "adding_photo"
	StateSettingTime AdminState = "setting_time"
)

/***************************** РАБОТА С АВТОРИЗАЦИЕЙ *****************************/

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

/***************************** РАБОТА С НИКАМИ *****************************/

// saveNicknames сохраняет список никнеймов в JSON-файл
func saveNicknames(filePath string, list []UserNickname) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// loadNicknames загружает список никнеймов из JSON-файла
func loadNicknames(filePath string) ([]UserNickname, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyList := []UserNickname{}
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

	var list []UserNickname
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// getNickname возвращает никнейм пользователя или пустую строку
func getNickname(chatID int64, nicknames []UserNickname) string {
	for _, user := range nicknames {
		if user.ChatID == chatID {
			return user.Nickname
		}
	}
	return ""
}

// setNickname устанавливает или обновляет никнейм пользователя
func setNickname(chatID int64, nickname string, nicknames []UserNickname) []UserNickname {
	for i, user := range nicknames {
		if user.ChatID == chatID {
			nicknames[i].Nickname = nickname
			return nicknames
		}
	}
	return append(nicknames, UserNickname{ChatID: chatID, Nickname: nickname})
}

/***************************** РАБОТА С ПАНЕЛЬЮ АДМИНА *****************************/

// loadPosts загружает список постов из JSON-файла
func loadPosts(filePath string) ([]Post, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyList := []Post{}
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

	var list []Post
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// savePosts сохраняет список постов в JSON-файл
func savePosts(filePath string, list []Post) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func generateUUID() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		randBytes(4), randBytes(2), randBytes(2), randBytes(2), randBytes(6))
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

/**********************************************************/

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

	// Загружаем списки авторизованных пользователей, администраторов, никнеймов и постов
	authorizedUsers, err := loadAuthorizedList(AuthorizedUsersFile)
	if err != nil {
		log.Fatalf("Failed to load authorized users: %v", err)
	}
	authorizedAdmins, err := loadAuthorizedList(AuthorizedAdminsFile)
	if err != nil {
		log.Fatalf("Failed to load authorized admins: %v", err)
	}
	nicknames, err := loadNicknames(NicknamesFile)
	if err != nil {
		log.Fatalf("Failed to load nicknames: %v", err)
	}
	posts, err := loadPosts(PostsFile)
	if err != nil {
		log.Fatalf("Failed to load posts: %v", err)
	}

	// Логируем списки для отладки
	log.Printf("Authorized users: %v", authorizedUsers)
	log.Printf("Authorized admins: %v", authorizedAdmins)

	// Хранилище текущих создаваемых постов
	creatingPosts := make(map[int64]*CreatingPost)

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		chatID := int64(0)
		var msg tgbotapi.MessageConfig

		// Обработка callback-запросов от кнопок
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
			msg = tgbotapi.NewMessage(chatID, "")
			data := update.CallbackQuery.Data

			// Проверяем, является ли пользователь админом
			if isAuthorizedInList(chatID, authorizedAdmins) {
				// Получаем или создаем текущий пост
				creatingPost, exists := creatingPosts[chatID]
				if !exists {
					creatingPost = &CreatingPost{AdminChatID: chatID, State: StateIdle}
					creatingPosts[chatID] = creatingPost
				}

				switch data {
				case "add_text":
					creatingPost.State = StateAddingText
					msg.Text = "Введите текст поста:"
					bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
				case "add_photo":
					if len(creatingPost.PhotoIDs) >= MaxMediaCount {
						msg.Text = fmt.Sprintf("Достигнуто максимальное количество фотографий (%d).", MaxMediaCount)
						bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					} else {
						creatingPost.State = StateAddingPhoto
						msg.Text = "Отправьте фотографию для поста:"
						bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					}
				case "preview_post":
					if creatingPost.Text == "" && len(creatingPost.PhotoIDs) == 0 {
						msg.Text = "Пост пуст. Добавьте текст или фотографии."
						bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					} else {
						log.Printf("Preview post: text=%s, photos=%d", creatingPost.Text, len(creatingPost.PhotoIDs))
						caption := creatingPost.Text
						keyboard := tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Опубликовать", "publish_post"),
								tgbotapi.NewInlineKeyboardButtonData("Отменить", "cancel_post"),
							),
						)
						if len(creatingPost.PhotoIDs) > 0 {
							var media []interface{}
							for i, photoID := range creatingPost.PhotoIDs {
								inputMedia := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(photoID))
								if i == 0 {
									inputMedia.Caption = caption // Текст как подпись к первому фото
								}
								media = append(media, inputMedia)
							}
							mediaGroup := tgbotapi.NewMediaGroup(chatID, media)
							if _, err := bot.SendMediaGroup(mediaGroup); err != nil {
								log.Printf("Failed to send media group: %v", err)
							}
							// Отправляем отдельное сообщение с кнопками
							msg.Text = "Выберите действие:"
							msg.ReplyMarkup = keyboard
						} else {
							// Если нет фотографий, отправляем текст поста с кнопками
							msg.Text = caption
							msg.ReplyMarkup = keyboard
						}
						if _, err := bot.Send(msg); err != nil {
							log.Printf("Failed to send message with buttons: %v", err)
						}
						bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
						// Очищаем msg.Text, чтобы избежать дублирования
						msg.Text = ""
					}
				case "publish_post":
					if creatingPost.Text == "" && len(creatingPost.PhotoIDs) == 0 {
						msg.Text = "Нельзя опубликовать пустой пост."
						bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					} else {
						post := Post{
							ID:          generateUUID(),
							AdminChatID: chatID,
							Text:        creatingPost.Text,
							PhotoIDs:    creatingPost.PhotoIDs,
							Published:   true,
						}
						// Объединяем списки пользователей и админов, избегая дубликатов
						recipients := make(map[int64]bool)
						for _, userID := range authorizedUsers {
							recipients[userID] = true
						}
						for _, adminID := range authorizedAdmins {
							recipients[adminID] = true
						}
						log.Printf("Sending post to recipients: %v", recipients)

						// Немедленная публикация
						for userID := range recipients {
							if len(post.PhotoIDs) > 0 {
								var media []interface{}
								for j, photoID := range post.PhotoIDs {
									inputMedia := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(photoID))
									if j == 0 {
										inputMedia.Caption = post.Text
									}
									media = append(media, inputMedia)
								}
								mediaGroup := tgbotapi.NewMediaGroup(userID, media)
								if _, err := bot.SendMediaGroup(mediaGroup); err != nil {
									log.Printf("Failed to send media group to user %d: %v", userID, err)
								} else {
									log.Printf("Sent media group to user %d", userID)
								}
							} else {
								userMsg := tgbotapi.NewMessage(userID, post.Text)
								if _, err := bot.Send(userMsg); err != nil {
									log.Printf("Failed to send message to user %d: %v", userID, err)
								} else {
									log.Printf("Sent message to user %d: %s", userID, post.Text)
								}
							}
						}
						log.Printf("Saving post: %+v", post)
						posts = append(posts, post)
						if err := savePosts(PostsFile, posts); err != nil {
							log.Printf("Failed to save posts: %v", err)
							msg.Text = "Ошибка при сохранении поста. Попробуйте позже."
							bot.Send(msg)
						} else {
							// Уведомление админа с клавиатурой
							adminMsg := tgbotapi.NewMessage(post.AdminChatID, "Пост успешно опубликован! Начните создание нового поста:")
							adminMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
									tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
								),
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
								),
							)
							if _, err := bot.Send(adminMsg); err != nil {
								log.Printf("Failed to send notification to admin %d: %v", post.AdminChatID, err)
							} else {
								log.Printf("Notified admin %d about publication", post.AdminChatID)
							}
						}
						delete(creatingPosts, chatID)
						bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					}
				case "cancel_post":
					delete(creatingPosts, chatID)
					msg.Text = "Создание поста отменено. Начните создание нового поста:"
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
							tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
						),
					)
					msg.ReplyMarkup = keyboard
					bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
				}
			} else if isAuthorizedInList(chatID, authorizedUsers) {
				switch data {
				case "create_nick", "change_nick":
					msg.Text = "Введите ваш новый никнейм:"
					bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
				case "my_nick":
					nick := getNickname(chatID, nicknames)
					if nick == "" {
						msg.Text = "У вас еще нет никнейма. Создайте его с помощью кнопки 'Создать ник'."
					} else {
						msg.Text = fmt.Sprintf("Ваш текущий никнейм: %s", nick)
					}
					bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
				}
			}
			// Отправляем msg только если оно явно задано
			if msg.Text != "" {
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
				} else {
					log.Printf("Sent message: %s", msg.Text)
				}
			}
			continue
		}

		chatID = update.Message.Chat.ID
		msg = tgbotapi.NewMessage(chatID, "")

		// Обработка команды /start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			args := update.Message.CommandArguments()
			if args == SecretCode {
				if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
					authorizedUsers = append(authorizedUsers, chatID)
					if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
						log.Printf("Failed to save authorized users: %v", err)
						msg.Text = "Ошибка при сохранении авторизации. Попробуйте позже."
					} else {
						msg.Text = "Вы подключились с секретным кодом! Теперь вы можете использовать бота как пользователь."
						keyboard := tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
								tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
								tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
							),
						)
						msg.ReplyMarkup = keyboard
					}
				} else {
					msg.Text = "Вы уже авторизованы."
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
						),
					)
					msg.ReplyMarkup = keyboard
				}
			} else if args == AdminCode {
				if !isAuthorizedInList(chatID, authorizedAdmins) {
					if isAuthorizedInList(chatID, authorizedUsers) {
						authorizedUsers = removeFromList(chatID, authorizedUsers)
						if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
							log.Printf("Failed to save authorized users: %v", err)
							msg.Text = "Ошибка при обновлении списка пользователей. Попробуйте позже."
						} else {
							authorizedAdmins = append(authorizedAdmins, chatID)
							if err := saveAuthorizedList(AuthorizedAdminsFile, authorizedAdmins); err != nil {
								log.Printf("Failed to save authorized admins: %v", err)
								msg.Text = "Ошибка при сохранении авторизации администратора. Попробуйте позже."
							} else {
								msg.Text = "Вы повышены до администратора! Теперь вы можете использовать бота как администратор."
								keyboard := tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
										tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
									),
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
									),
								)
								msg.ReplyMarkup = keyboard
							}
						}
					} else {
						authorizedAdmins = append(authorizedAdmins, chatID)
						if err := saveAuthorizedList(AuthorizedAdminsFile, authorizedAdmins); err != nil {
							log.Printf("Failed to save authorized admins: %v", err)
							msg.Text = "Ошибка при сохранении авторизации администратора. Попробуйте позже."
						} else {
							msg.Text = "Вы подключились с кодом администратора! Теперь вы можете использовать бота как администратор."
							keyboard := tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
									tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
								),
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
								),
							)
							msg.ReplyMarkup = keyboard
						}
					}
				} else {
					msg.Text = "Вы уже администратор."
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
							tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
						),
					)
					msg.ReplyMarkup = keyboard
				}
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			} else if msg.Text != "" {
				log.Printf("Sent message: %s", msg.Text)
			}
			continue
		}

		chatID = update.Message.Chat.ID
		msg = tgbotapi.NewMessage(chatID, "")

		// Проверяем, авторизован ли пользователь
		if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
			continue
		}

		// Обработка сообщений
		if isAuthorizedInList(chatID, authorizedUsers) && !isAuthorizedInList(chatID, authorizedAdmins) {
			// Обработка установки никнейма для обычных пользователей
			if update.Message.Text != "" {
				nickname := update.Message.Text
				nicknames = setNickname(chatID, nickname, nicknames)
				if err := saveNicknames(NicknamesFile, nicknames); err != nil {
					log.Printf("Failed to save nicknames: %v", err)
					msg.Text = "Ошибка при сохранении никнейма. Попробуйте позже."
				} else {
					msg.Text = fmt.Sprintf("Ваш никнейм установлен: %s", nickname)
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
						),
					)
					msg.ReplyMarkup = keyboard
				}
			}
		} else if isAuthorizedInList(chatID, authorizedAdmins) {
			creatingPost, exists := creatingPosts[chatID]
			if !exists {
				creatingPost = &CreatingPost{AdminChatID: chatID, State: StateIdle}
				creatingPosts[chatID] = creatingPost
			}

			switch creatingPost.State {
			case StateAddingText:
				if update.Message.Text != "" {
					if len(update.Message.Text) > MaxCaptionLength {
						msg.Text = fmt.Sprintf("Текст слишком длинный. Максимум %d символов.", MaxCaptionLength)
					} else {
						creatingPost.Text = update.Message.Text
						creatingPost.State = StateIdle
						msg.Text = "Текст поста сохранен. Продолжайте редактировать пост."
						keyboard := tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
								tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
							),
						)
						msg.ReplyMarkup = keyboard
					}
				}
			case StateAddingPhoto:
				if update.Message.Photo != nil && len(update.Message.Photo) > 0 {
					photo := update.Message.Photo[len(update.Message.Photo)-1]
					creatingPost.PhotoIDs = append(creatingPost.PhotoIDs, photo.FileID)
					creatingPost.State = StateIdle
					msg.Text = "Фотография добавлена. Можете добавить еще или продолжить."
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
							tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
						),
					)
					msg.ReplyMarkup = keyboard
				}
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		} else if msg.Text != "" {
			log.Printf("Sent message: %s", msg.Text)
		}
	}
}
