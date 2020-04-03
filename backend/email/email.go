package email

import (
	"bytes"
	"path/filepath"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

type Email interface {
	Send(msg *gomail.Message) error
	GetTemplate(name string, data interface{}) (string, error)
	LoadTemplates(dir string) error
	GetAddress() string
}

type email struct {
	address string
	dialer  *gomail.Dialer
	tpl     *template.Template
}

var DefaultInstance = &email{}

func (e *email) Send(msg *gomail.Message) error {
	return e.dialer.DialAndSend(msg)
}

func (e *email) GetTemplate(name string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	if err := e.tpl.ExecuteTemplate(&tpl, name, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (e *email) LoadTemplates(dir string) error {
	tpl, err := template.ParseGlob(filepath.Join(dir, "*.gohtml"))
	if err != nil {
		return err
	}
	e.tpl = tpl
	return nil
}

func (e *email) GetAddress() string {
	return e.address
}

func New(host string, port int, username, password, address string) Email {
	return &email{dialer: gomail.NewDialer(host, port, username, password), address: address}
}

func NewDialer(host string, port int, username, password string) {
	DefaultInstance.dialer = gomail.NewDialer(host, port, username, password)
	if DefaultInstance.address == "" {
		DefaultInstance.address = username
	}
}

func Send(msg *gomail.Message) error {
	return DefaultInstance.Send(msg)
}

func GetTemplate(template string, data interface{}) (string, error) {
	return DefaultInstance.GetTemplate(template, data)
}

func LoadTemplates(dir string) error {
	return DefaultInstance.LoadTemplates(dir)
}

func GetAddress() string {
	return DefaultInstance.GetAddress()
}
