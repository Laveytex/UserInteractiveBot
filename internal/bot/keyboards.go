package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func UserKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
			tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
			tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
		),
	)
}

func AdminKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить текст", "add_text"),
			tgbotapi.NewInlineKeyboardButtonData("Добавить фото", "add_photo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Превью поста", "preview_post"),
		),
	)
}
