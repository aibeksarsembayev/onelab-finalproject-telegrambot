package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
}

type IncomingMessage struct {
	Chat_id     int                  `json:"chat_id"`
	Text        string               `json:"text"`
	From        From                 `json:"from"`
	Chat        Chat                 `json:"chat"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type CallbackQuery struct {
	ID string `json:"id"`
	// From         From            `json:"from"`
	Message IncomingMessage `json:"message"`
	// ChatInstance string          `json:"chat_instance"`
	Data string `json:"data"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"`
	CallbackData string `json:"callback_data"`
}
