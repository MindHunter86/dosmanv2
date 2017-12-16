package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api"


type tgrmCustomer struct {
	tgb *tgbotapi.BotAPI
}

func (m *tgrmCustomer) configure(tgb *tgbotapi.BotAPI) (*tgrmCustomer,error) {
	m.tgb = tgb
	return m,nil
}

func (m *tgrmCustomer) requestContact(chatid int64) {
	msgAgreement := tgbotapi.NewMessage(chatid, "Для получения уведомлений о входе на сервера от меня, тебе необходимо зарегистрироваться. При нажатии на кнопку, телеграм расшарит твой телефон для меня. После я тебя запишу тебя в свою базу и буду отправлять уведомления по команде в мою апишку.")
	msgAgreement.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonContact("Заргеистрировать в системе"),
	})
	m.tgb.Send(msgAgreement)
}

func (m *tgrmCustomer) registerContact(chatid int64, fromid int, contact *tgbotapi.Contact) error {
	if contact.UserID != fromid {
		m.tgb.Send(tgbotapi.NewMessage(
			chatid, "Номер телефона, который ты мне дал, не принадлежит тебе! Просто нажми на кнопку ДА и не страдай фигнёй!"))
		return nil
	}

//	if e := m.updateCustomer(contact); e != nil { return e }

	msg := tgbotapi.NewMessage(chatid, "Ты был зарегистрирован в системе, поздравляю! Теперь все уведомления о заходе на сервера по SSH будут приходить в этот чатик. Если есть вопросы спрашивай у Бажина.")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	if _,e := m.tgb.Send(msg); e != nil { return e }

	return nil
}
