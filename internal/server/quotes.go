package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed quotes.json
var quotesJson []byte

type Quote struct {
	Text   string `json:"q"`
	Author string `json:"a"`
}

func (q Quote) String() string {
	return fmt.Sprintf("%s -- %s", q.Text, q.Author)
}

func loadQuotes(raw []byte) ([]Quote, error) {
	var quotes []Quote
	err := json.Unmarshal(raw, &quotes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal quotes: %w", err)
	}
	return quotes, nil
}
