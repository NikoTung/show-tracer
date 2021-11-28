package main

type Config struct {
	Rss   string `json:"rss"`
	Since string `json:"since"`
	Dt    string `json:"dt"`
	Api   string `json:"api"`
}

type Aria2 struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Id      string   `json:"id"`
	Params  []string `json:"params"`
}

func NewAria2(id string, links string) *Aria2 {

	return &Aria2{
		Jsonrpc: "2.0",
		Method:  "aria2.addUri",
		Id:      id,
		Params:  []string{links},
	}
}
