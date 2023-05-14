package model

import "fmt"

type LoginPassword struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewLoginPassword(login string, password string) *LoginPassword {
	return &LoginPassword{Login: login, Password: password}
}

func (p *LoginPassword) Format(description string) string {
	return fmt.Sprintf("\nlogin:%s\npassword:%s\ndescription:%s", p.Login, p.Password, description)
}
