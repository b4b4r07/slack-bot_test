package quote

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

const QuoteAPI = "http://api.forismatic.com/api/1.0/?method=getQuote&lang=en&format=json"

type Quote struct {
	Text       string `json:"quoteText"`
	Author     string `json:"quoteAuthor"`
	SenderName string `json:"senderName"`
	SenderLink string `json:"senderLink"`
	QuoteLink  string `json:"quoteLink"`
}

func Get() *Quote {
	req := gorequest.New()
	_, body, _ := req.Post(QuoteAPI).End()

	quote := &Quote{}
	json.Unmarshal([]byte(body), &quote)

	return quote
}
