package telegram

import (
	"reflect"

	"mh00appserver/modules"

	"github.com/rs/zerolog"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)


type TelegramModule struct {
	modName string
	tbot *tgbotapi.BotAPI
	tbotUpdates tgbotapi.UpdatesChannel
	tbotCustomer *tgrmCustomer

	mods *modules.Modules
	logger zerolog.Logger
}


// TelegramModule API:
func (m *TelegramModule) Configure(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	m.mods = mods
	m.modName = reflect.TypeOf(m).Elem().Name()
	m.logger = m.mods.Logger.With().Str("MODULE", m.modName).Logger()


	return m,nil
}

func (m *TelegramModule) Bootstrap() error {
	if e := m.telegramAuthorization(); e != nil { return e }

	m.logger.Debug().Msg(m.modName+" has been bottstrapped!")
LOOP:
	for {
		select{
		case <-m.mods.DonePipe:
			break LOOP
		case updt := <-m.tbotUpdates:
			if updt.Message == nil { continue }

			if updt.Message.IsCommand() {
				if e := m.commandRouter(updt.Message); e != nil {
					m.logger.Warn().Err(e).Msg("Some errors in telegram command router!")
				}
			} else if updt.Message.Contact != nil {
				if e := m.tbotCustomer.registerContact(updt.Message.Chat.ID, updt.Message.From.ID, updt.Message.Contact); e != nil {
					m.logger.Warn().Err(e).Msg("Some errors in telegram customer registration!")
				}
			} else {
				m.logger.Debug().Msg(updt.Message.Text)
				continue // TODO: add warning message?
			}
		}
	}

	return nil
}


// TelegramModule internal methods:
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

func (m *TelegramModule) commandRouter(msg *tgbotapi.Message) error {
	switch msg.Command() {
	case "start":
		m.tbotCustomer.requestContact(msg.Chat.ID)
	default:
		m.logger.Warn().Str("customer_id", msg.From.UserName).Msg("Unknown command received from customer - " + msg.Text)
	}

	return nil
}
