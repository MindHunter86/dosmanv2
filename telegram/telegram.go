package telegram

import "mailru/rooster22/config"

// import "github.com/rs/zerolog"
import "golang.org/x/net/context"
import "github.com/go-telegram-bot-api/telegram-bot-api"


type TelegramBot struct {
	*tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel

	ctxPipeDone <-chan struct{}
}
func (self *TelegramBot) ConfigureAndConnect(ctx context.Context) (*TelegramBot, error) {
	var e error
	var appConfig *config.AppConfig

	// import config from context:
	self.ctxPipeDone = ctx.Done()
	appConfig = ctx.Value(config.CTX_APP_CONFIG).(*config.AppConfig)

	// initialize new telegram bot connention:
	if self.BotAPI, e = tgbotapi.NewBotAPI(appConfig.Telegram.Token); e != nil { return nil,e }

	// telegram bot event loop configuration:
	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = 60

	if self.updates, e = self.GetUpdatesChan(updatesConfig); e != nil { return nil,e }
	go self.getUpdates()

	return self,nil
}
func (self *TelegramBot) getUpdates() {
	var event tgbotapi.Update

	for event = range self.updates {
		if event.Message == nil { continue }
		if event.Message.IsCommand() == false { continue }
		if event.Message.Command() != "start" { continue }

		// check event.Message.From.ID in database:
		// TODO

		// show msg with agreement:
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, event.Message.Text)
		msg.ReplyToMessageID = event.Message.MessageID
		self.Send(msg)

		msgAgreement := tgbotapi.NewMessage(event.Message.Chat.ID, "Do you agree terms?")
		msgAgreement.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButtonContact("Yes"),
			tgbotapi.NewKeyboardButton("No"),
		})
		self.Send(msgAgreement)
	}
}



//
//func main() {
//	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
//	log.Debug().Msg("Mecroservice has been initialized!")
//
//	tbot, e := tgbotapi.NewBotAPI(CONFIG_TELEGRAM_TOKEN)
//	if e != nil { log.Error().Str("NewBotAPI function has been failed!", e.Error()) }
//
//	tbot.Debug = CONFIG_MAIN_DEBUG
//	log.Info().Str("Authorized on account:", tbot.Self.UserName)
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	tbotUpdates, e := tbot.GetUpdatesChan(u)
//	if e != nil { log.Error().Str("GetUpdatesChan has been failed!", e.Error()) }
//
//	tbotsync := new(sync.WaitGroup)
//	go func(wg *sync.WaitGroup) {
//		wg.Add(1)
//		for up := range tbotUpdates {
//			if up.Message == nil { continue }
//
//			log.Info().Str("From:", up.Message.From.UserName).Str("Text:", up.Message.Text).Msg("NEW MESSAGE!")
//
//			msg := tgbotapi.NewMessage(up.Message.Chat.ID, up.Message.Text)
//			msg.ReplyToMessageID = up.Message.MessageID
//			tbot.Send(msg)
//
//			var requestForm *tlgrRequestForm = new(tlgrRequestForm)
//			requestForm.message = tgbotapi.NewMessage(up.Message.Chat.ID, "Do you agree terms?")
//			requestForm.button_accept = tgbotapi.NewKeyboardButtonContact("Yes")
//			requestForm.button_decline = tgbotapi.NewKeyboardButton("No")
//			requestForm.message.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{requestForm.button_accept, requestForm.button_decline})
//			tbot.Send(requestForm.message)
//		}
//		wg.Done()
//	}(tbotsync)
//
//	time.Sleep(1 * time.Second)
//
//	tbotsync.Wait()
//}
