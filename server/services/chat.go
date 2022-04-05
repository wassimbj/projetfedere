package services

import (
	"context"
	"pfserver/db"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type C struct{}

func Chat() C {
	return C{}
}

type CreateData struct {
	SentFrom int
	SentTo   int
	Msg      string
}

func (C) Create(ctx context.Context, data CreateData) error {

	_, err := db.Conn().Exec(ctx, `
	INSERT INTO messages(sent_from, sent_to, msg) VALUES($1,$2,$3) 
	`, data.SentFrom, data.SentTo, data.Msg)

	return err
}

type Message struct {
	Id        int       `json:"id"`
	SentFrom  int       `json:"sentFrom"`
	SentTo    int       `json:"sentTo"`
	Msg       string    `json:"msg"`
	CreatedAt time.Time `json:"createdAt"`
}

// get the messages of the two peers
func (C) Get(ctx context.Context, otherPeerId, myId int) ([]*Message, error) {
	var data []*Message
	err := pgxscan.Select(ctx, db.Conn(), &data, `
			SELECT id, sent_from, sent_to, msg, created_at
			FROM messages
			WHERE sent_from = $1 OR sent_to = $1
			ORDER BY created_at ASC
		`, otherPeerId)
	return data, err

}
