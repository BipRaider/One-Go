package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to BipGo.pw" // Мема песьма
)

// текст что будет появляться если сообщений было оправлено несколько
const welcomeText = `Hi there!
Welcome to BipGo.pw ! We really hope you enjov using our application

Best,
Bip
`

// текст что будет в нутри письма
const welcomeHTML = `Hi there! <br/>
<br/>
Welcom to
<a href="https://bipus.Bipgo.pw "> BipGo.pw</a>! We really  hope you enjoy using our applocation!<br/>
<br/>
Best ,<br/>
Bip
`
const resetTextTmpl = `Hi there!

It appears that you have requested a password reset. If this was you, please follow the link below to update your password:

%s

If you are asked for a token, please use the following value:

%s

If you didn't request a password reset you can safely ignore this email and your account will not be changed.

Best,
LensLocked Support
`

const resetHTMLTmpl = `Hi there!<br/>
<br/>
It appears that you have requested a password reset. If this was you, please follow the link below to update your password:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this email and your account will not be changed.<br/>
<br/>
Best,<br/>
LensLocked Support<br/>
`

// Нашь идефикатор пользователя , ключи к серевру отправителя и от акаунта
func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey) //
		c.mg = mg
	}
}

// 1 Имя отправителя ,2 Email c которого отправили
// собирает пиьсмо от кого отправляется
func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		from: "bipusgo@gmail.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcom(toName, toEmail string) error {
	//1 от кого,2 Имя писма ,3 содержание , 4 кому отправлять
	message := mailgun.NewMessage(buildEmail(toName, toEmail), welcomeSubject, welcomeText, "bipusgo@gmail.com")
	//организует связывание HTML-представления вашего
	//сообщения в дополнение к вашему текстовому сообщению.
	message.SetHtml(welcomeHTML)
	_, id, err := c.mg.Send(message)
	fmt.Println("ID=", id)
	return err
	//https://documentation.mailgun.com/en/latest/user_manual.html#sending-via-api
}

const (
	resetSubject = "Instructions for resetting your password."
	resetBaseURL = "https://your_adres"
)

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
	_, _, err := c.mg.Send(message)
	return err
}

// Функция собирающая  строку   КОМУ отправлять письмо.
func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)

}
