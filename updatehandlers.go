package telebot

import (
	"encoding/json"
	"reflect"

	"github.com/GoLibs/telegram-bot-api/structs"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type updateHandlers struct {
	bot                *Bot
	messageTypeGroup   reflect.Method
	channelPost        reflect.Method
	inlineQuery        reflect.Method
	chosenInlineResult reflect.Method
	callbackQuery      reflect.Method
	shippingQuery      reflect.Method
	preCheckoutQuery   reflect.Method
	pollAnswer         reflect.Method
	myChatMember       reflect.Method
	chatMember         reflect.Method
}

func updateHandlersFromType(bot *Bot, t reflect.Type) (uh *updateHandlers) {
	uh = &updateHandlers{
		bot: bot,
	}
	uh.messageTypeGroup, _ = t.MethodByName("MessageTypeGroupHandler")
	uh.channelPost, _ = t.MethodByName("ChannelPostHandler")
	uh.inlineQuery, _ = t.MethodByName("InlineQueryHandler")
	uh.chosenInlineResult, _ = t.MethodByName("ChosenInlineResultHandler")
	uh.callbackQuery, _ = t.MethodByName("CallbackQueryHandler")
	uh.shippingQuery, _ = t.MethodByName("ShippingQueryHandler")
	uh.preCheckoutQuery, _ = t.MethodByName("PreCheckoutQueryHandler")
	uh.pollAnswer, _ = t.MethodByName("PollAnswerHandler")
	uh.myChatMember, _ = t.MethodByName("MyChatMemberHandler")
	uh.chatMember, _ = t.MethodByName("ChatMemberHandler")
	return
}

func (uh *updateHandlers) allowedUpdates() (allowedUpdates []string) {
	allowedUpdates = append(allowedUpdates, "message", "edited_message")

	if uh.channelPost.Name != "" {
		allowedUpdates = append(allowedUpdates, "channel_post", "edited_channel_post")
	}
	if uh.inlineQuery.Name != "" {
		allowedUpdates = append(allowedUpdates, "inline_query")
	}
	if uh.chosenInlineResult.Name != "" {
		allowedUpdates = append(allowedUpdates, "chosen_inline_result")
	}
	if uh.callbackQuery.Name != "" {
		allowedUpdates = append(allowedUpdates, "callback_query")
	}
	if uh.shippingQuery.Name != "" {
		allowedUpdates = append(allowedUpdates, "shipping_query")
	}
	if uh.preCheckoutQuery.Name != "" {
		allowedUpdates = append(allowedUpdates, "pre_checkout_query")
	}
	if uh.chatMember.Name != "" {
		allowedUpdates = append(allowedUpdates, "chat_member")
	}
	if uh.pollAnswer.Name != "" {
		allowedUpdates = append(allowedUpdates, "poll_answer")
	}
	if uh.myChatMember.Name != "" {
		allowedUpdates = append(allowedUpdates, "my_chat_member")
	}
	return
}

func (uh *updateHandlers) handleProcessUpdateError(update *go_telegram_bot_api.Update, message string) {
	j, _ := json.Marshal(update)
	log.Errorf("Error processing update: %s. Update: %s", message, j)
}

func (uh *updateHandlers) processMessageTypeGroup(app reflect.Value, update *go_telegram_bot_api.Update) {
	if uh.messageTypeGroup.Name == "" {
		j, _ := json.Marshal(update)
		log.Errorf("Handler for Group Message Type was not found!\n%s", string(j))
		return
	}
	uh.messageTypeGroup.Func.Call(uh.bot.appWithUpdate(app, update, &update.Message.Chat.Id))
}

func (uh *updateHandlers) processChannelPost(app reflect.Value, update *go_telegram_bot_api.Update) {
	if uh.channelPost.Name == "" {
		j, _ := json.Marshal(update)
		log.Errorf("Handler for ChannelPost was not found!\n%s", string(j))
		return
	}

	uh.channelPost.Func.Call(uh.bot.appWithUpdate(app, update, &update.ChannelPost.Chat.Id))
}

func (uh *updateHandlers) processMyChatMember(app reflect.Value, update *go_telegram_bot_api.Update) {
	if uh.myChatMember.Name == "" {
		j, _ := json.Marshal(update)
		log.Errorf("Handler for MyChatMember was not found!\n%s", string(j))
		return
	}
	uh.myChatMember.Func.Call(uh.bot.appWithUpdate(app, update, &update.MyChatMember.Chat.Id))
}

func (uh *updateHandlers) processChatMember(app reflect.Value, update *go_telegram_bot_api.Update) {
	if uh.chatMember.Name == "" {
		j, _ := json.Marshal(update)
		log.Errorf("Handler for ChatMember was not found!\n%s", string(j))
		return
	}

	uh.chatMember.Func.Call(uh.bot.appWithUpdate(app, update, &update.ChatMember.Chat.Id))
}

func (uh *updateHandlers) processCallbackQuery(app reflect.Value, update *go_telegram_bot_api.Update) {
	if uh.callbackQuery.Name == "" {
		j, _ := json.Marshal(update)
		log.Errorf("Handler for CallbackQuery was not found!\n%s", string(j))
		return
	}

	uh.callbackQuery.Func.Call(uh.bot.appWithUpdate(app, update, &update.CallbackQuery.Message.Chat.Id))
}

func (uh *updateHandlers) processPollAnswer(app reflect.Value, update *go_telegram_bot_api.Update) {
	if uh.pollAnswer.Name == "" {
		j, _ := json.Marshal(update)
		log.Errorf("Handler for PollAnswer was not found!\n%s", string(j))
		return
	}
	uh.pollAnswer.Func.Call(uh.bot.appWithUpdate(app, update, &update.PollAnswer.User.Id))
}

func (uh *updateHandlers) processUpdate(app reflect.Value, update *go_telegram_bot_api.Update) {
	var message *structs.Message
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	}

	var channelPost *structs.Message
	if update.ChannelPost != nil {
		channelPost = update.ChannelPost
	} else if update.EditedChannelPost != nil {
		channelPost = update.EditedChannelPost
	}
	if message != nil && (message.Chat.Type == "group" || message.Chat.Type == "supergroup") {
		uh.processMessageTypeGroup(app, update)
	} else if channelPost != nil {
		uh.processChannelPost(app, update)
	} else if update.MyChatMember != nil {
		uh.processMyChatMember(app, update)
	} else if update.ChatMember != nil {
		uh.processChatMember(app, update)
	} else if update.CallbackQuery != nil {
		uh.processCallbackQuery(app, update)
	} else if update.PollAnswer != nil {
		uh.processPollAnswer(app, update)
	} else {
		uh.handleProcessUpdateError(update, "message_type_not_supported")
	}
}
