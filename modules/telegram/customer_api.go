package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

type chat struct { id int }
type customer struct {
	tgid  int
	firstname, lastname, phone string
	chatid *chat
}


// updt.Message.Chat.ID, updt.Message.From.ID, updt.Message.Contact
func (m *customer) create(db *dbDriver, mess *tgbotapi.Message) (*customer, error) {

	m.tgid = mess.Contact.UserID
	m.firstname = mess.Contact.FirstName
	m.lastname = mess.Contact.LastName
	m.phone = mess.Contact.PhoneNumber
	m.chatid = &chat{id: mess.Chat.ID}

	return m,nil
}
