package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int                `json:"update_id"`
	Message       *IncomingMessage   `json:"message"`
	ChannelPost   *IncomingMessage   `json:"channel_post"`
	CallbackQuery *CallbackQuery     `json:"callback_query"`
	MyChatMember  *ChatMemberUpdated `json:"my_chat_member"`
	EditedMessage *IncomingMessage   `json:"edited_message"`
}

type IncomingMessage struct {
	Chat_id     int                  `json:"chat_id"`
	Text        string               `json:"text"`
	From        User                 `json:"from"`
	SenderChat  Chat                 `json:"sender_chat"`
	Chat        Chat                 `json:"chat"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type SendMessage struct {
	ChatID      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ParseMode   string               `json:"parse_mode"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type CallbackQuery struct {
	ID      string          `json:"id"`
	From    User            `json:"from"`
	Message IncomingMessage `json:"message"`
	// ChatInstance string          `json:"chat_instance"`
	Data string `json:"data"`
}

// This object represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	Chat Chat `json:"chat"`
	From User `json:"from"`
	// Date          int        `json:"date"`
	OldChatMember ChatMember `json:"old_chat_member"`
	NewChatMember ChatMember `json:"new_chat_member"`
}

type ChatMember struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}
type User struct {
	Username string `json:"username"`
}

type Chat struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Date     int    `json:"date"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"`
	CallbackData string `json:"callback_data"`
}
