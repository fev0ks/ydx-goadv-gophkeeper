package resources

import (
	"fmt"

	"ydx-goadv-gophkeeper/pkg/model/enum"
)

type BankCard struct {
	Number   string `json:"number,omitempty"`
	ExpireAt string `json:"expireAt"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
}

func NewBankCard(number string, expireAt string, name string, surname string) *BankCard {
	return &BankCard{Number: number, ExpireAt: expireAt, Name: name, Surname: surname}
}

func (b *BankCard) Format(description string) string {
	return fmt.Sprintf("number: %v\nexpireAt: %v\nname: %v\nsurname: %v\ndescription: %v",
		b.Number,
		b.ExpireAt,
		b.Name,
		b.Surname,
		description,
	)
}

func (p *BankCard) Type() enum.ResourceType {
	return enum.BankCard
}
