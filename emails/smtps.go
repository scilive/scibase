package emails

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"strings"

	"github.com/flosch/pongo2/v4"
	"github.com/jordan-wright/email"
)

type EmailAuthType string

const (
	EmailAuthTypeStartTLS EmailAuthType = "starttls"
	EmailAuthTypeTLS      EmailAuthType = "tls"
	EmailAuthTypePlain    EmailAuthType = "plain"
)

type Email struct {
	From    string
	To      []string
	Cc      []string
	Subject string
	Text    string
	HTML    string
}

type Smtps struct {
	addr     string
	username string
	password string
	from     string
	authType EmailAuthType
}

func NewSMTPMailer(addr, username, password, from string, authType EmailAuthType) *Smtps {
	return &Smtps{
		addr:     addr,
		username: username,
		password: password,
		from:     from,
		authType: authType,
	}
}

func (s *Smtps) Send(m Email) error {
	switch s.authType {
	case EmailAuthTypeStartTLS:
		return s.SendStartTLS(m)
	case EmailAuthTypeTLS:
		return s.SendTLS(m)
	case EmailAuthTypePlain:
		return s.SendPlain(m)
	default:
		return s.SendTLS(m)
	}
}

func (s *Smtps) newEmail(m Email) *email.Email {
	e := email.NewEmail()
	e.From = s.from
	e.To = m.To
	e.Cc = m.Cc
	e.Subject = m.Subject
	e.Text = []byte(m.Text)
	e.HTML = []byte(m.HTML)
	return e
}
func (s *Smtps) SendPlain(m Email) error {
	e := s.newEmail(m)
	return e.Send(s.addr, smtp.PlainAuth("", s.username, s.password, strings.Split(s.addr, ":")[0]))
}

func (s *Smtps) SendStartTLS(m Email) error {
	e := s.newEmail(m)
	host := strings.Split(s.addr, ":")[0]
	auth := LoginAuth(s.from, s.password)
	tlsconfig := &tls.Config{
		ServerName: host,
	}
	return e.SendWithStartTLS(s.addr, auth, tlsconfig)
}

func (s *Smtps) SendTLS(m Email) error {
	e := s.newEmail(m)
	host := strings.Split(s.addr, ":")[0]
	auth := LoginAuth(s.from, s.password)
	tlsconfig := &tls.Config{
		ServerName: host,
	}
	return e.SendWithTLS(s.addr, auth, tlsconfig)
}

func (s *Smtps) SendTpl(to []string, cc []string, subject, tpl string, data map[string]any) error {
	template, err := pongo2.FromFile(tpl)
	if err != nil {
		return err
	}
	content, err := template.Execute(data)
	if err != nil {
		return err
	}
	e := Email{
		From:    s.from,
		To:      to,
		Cc:      cc,
		Subject: subject,
		HTML:    content,
		Text:    content,
	}
	return s.Send(e)
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}
