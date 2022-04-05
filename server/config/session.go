package config

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v8"
)

type Session struct {
	Req *http.Request
	Res http.ResponseWriter
}

var store *redisstore.RedisStore

func init() {
	client, err := Redis().Client()

	if err != nil {
		log.Fatal("Redis connection error: ", err)
	}

	var rsErr error
	store, rsErr = redisstore.NewRedisStore(context.Background(), client)

	if rsErr != nil {
		log.Fatal("init session Error: ", rsErr)
	}
	// session options
	sessMaxAge := (24 * 60 * 60) * 30 // 30 days
	if IsProd() {
		sessMaxAge = 5 * 60 * 60 // 5 hours
	}
	store.Options(sessions.Options{
		MaxAge:   sessMaxAge,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteNoneMode, // for ngrok Set-Cookie
		// Secure:   false,                 // for ngrok Set-Cookie
	})
}

func NewSession(req *http.Request, res http.ResponseWriter) *Session {
	return &Session{
		Req: req,
		Res: res,
	}
}

func (s *Session) Get(name string) (*sessions.Session, error) {
	return store.Get(s.Req, name)
}

type SessData map[interface{}]interface{}

func (s *Session) Save(name string, data SessData) error {
	session, _ := store.Get(s.Req, name)

	session.Values = data

	err := session.Save(s.Req, s.Res)

	return err
}

type UpdateData struct {
	Key interface{}
	Val interface{}
}

func (s *Session) Update(name string, data SessData) error {
	session, _ := store.Get(s.Req, name)

	for k, v := range data {
		session.Values[k] = v
	}

	err := session.Save(s.Req, s.Res)

	return err
}

func (s *Session) Del(name string) error {
	sess, _ := s.Get(name)

	// del session
	sess.Options.MaxAge = -1

	return sess.Save(s.Req, s.Res)
}

func (s *Session) GetUser() (map[interface{}]interface{}, error) {
	sess, err := s.Get("user")

	return sess.Values, err
}
