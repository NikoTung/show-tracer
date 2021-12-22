package main

type Config struct {
	Rss            string `json:"rss"`
	Since          string `json:"since"`
	Dt             string `json:"dt"`
	Api            string `json:"api"`
	Secret         string `json:"secret"`
	TelegramToken  string `json:"telegram_token"`
	TelegramChatId string `json:"telegram_chat_id"`
}

type Aria2 struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      string `json:"id"`
	Params  Params `json:"params"`
}

type Params []interface{}
type Token struct {
	Token string `json:"token"`
}

func NewAria2(id string, links string, secret string) *Aria2 {

	y := make(Params, 2)
	y[0] = "token:" + secret
	y[1] = []string{links}
	return &Aria2{
		Jsonrpc: "2.0",
		Method:  "aria2.addUri",
		Id:      id,
		Params:  y,
	}
}
