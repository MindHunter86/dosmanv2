package main

import "dosmanv2/system/config"
import "github.com/go-telegram-bot-api/telegram-bot-api"


type tgrmService struct {
	tbot *tgbotapi.BotAPI
	config *config.SysConfig
}


func (m *tgrmService) configure(cfg *config.SysConfig) (*tgrmService, error) {

	return m,nil
}

func (m *tgrmService) authenticate() error {
	var e error

	m.tbot, e = tgbotapi.NewBotAPI(m.config.Telegram.Token); if e != nil { return nil }

	return nil
}

func (m *TelegramModule) telegramAuthorization() error {
	var e error

	if m.tbot,e = tgbotapi.NewBotAPI(m.mods.Config.Telegram.Token); e != nil { return e }

	// load TelegramModule submodules:
	if m.tbotCustomer,e = new(tgrmCustomer).configure(m.tbot); e != nil { return e }

	tgUpdatesConfig := tgbotapi.NewUpdate(0)
	tgUpdatesConfig.Timeout = m.mods.Config.Telegram.Timeout
	if m.tbotUpdates,e = m.tbot.GetUpdatesChan(tgUpdatesConfig); e != nil { return e }

	return nil
}
