package telegram

const msgHelp = `@sber_invest_bot - сервис для получения статей с базы знаний sber-invest.kz.

<b>Сервис может:</b>
- получить все статьи;
- получить статьи за последние 7 дней;
- получить статьи по категориям;
- получить статьи по авторам;
- экспортировать статьи в формате pdf, csv;

<b>Главное меню</b>
/help - показывает это сообщение

<b>Статьи</b>
/articles - показывет фильтр статей
`

const msgHello = "<b>Привет!</b> 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Незнакомая команда 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
	msgStatusChanged  = "Статус изменен ✅"
)
