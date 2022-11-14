package telegram

const msgHelp = `@sber_invest_bot - service to get articles from knowledgebase on sber-invest.kz.

Service is able to:
- get all articles;
- get articles by category;
- get articles by author;
- export articles in pdf, csv;

*General*
/help - shows this message

*Articles*
/articles - shows article filter
`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
)
