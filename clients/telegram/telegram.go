package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// SendMessage by text and for methodget request
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// doRequest for methodget
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SendMessageButton with included Post method request
func (c *Client) SendMessageButton(chatID int, text string) error {
	var botMessage IncomingMessage
	botMessage.Chat.ID = chatID
	botMessage.Chat_id = chatID

	botMessage.Text = text
	botMessage.ReplyMarkup.InlineKeyboard = [][]InlineKeyboardButton{{{Text: "button1", CallbackData: "0"}}}

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return e.Wrap("can't send message", err)
	}
	// fmt.Println(string(buf))
	// fmt.Println(http.Post("https://"+c.host+"/"+c.basePath+"/sendMessage", "application/json", bytes.NewBuffer(buf)))

	res, err := http.Post("https://"+c.host+"/"+c.basePath+"/sendMessage", "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return err
	}

	buffer := make([]byte, 1000)
	res.Body.Read(buffer)

	for name, values := range res.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}

	fmt.Println(string(buffer))

	return nil
}

// SendMessage by IncomingMessage struct and for methodpost
func (c *Client) SendMessagePost(chatID int, msg *IncomingMessage) error {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.doRequestPost(sendMessageMethod, reqBody)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// doRequest for methodget
func (c *Client) doRequestPost(method string, reqBody []byte) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
