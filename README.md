# go-telebot

## Go Telegram Bot Framework 
Make and develop bots for telegram in a flash!

## Install
`go get -u github.com/aliforever/go-telebot`

## Getting Started
1. Create your `main.go` file in the project directory and import `_ github.com/aliforever/go-telebot` with underscore as alias (unused import) like:
```go
package main

import (
	_ "github.com/aliforever/go-telebot" 
)

func main() {
}
```
2. Get your bot token by sending `/newbot` to [BotFather](https://t.me/botfather) on Telegram.
3. Run this command in your project folder:
`go run main.go --new=BOT_TOKEN_HERE`
(Replace your token with BOT_TOKEN_HERE)
4. Enjoy developing your bot!

## Package Conventions
- The concept of `State` is used for private chats. [Read Below this section]
- Updates are divided to different handlers based on update types and each handler has its own file inside app package (eg: appprivate.go, appcallback.go)
  - app.go is your structure. You can embed `tgbotapi.TelegramBot` there to have easier access to bot api methods and I recommend to keep it that way unless you prefer another way.
  - appprivate.go file has 2 methods (`Welcome` & `Hello`) on the fly for welcoming users and responding hello to their hello! `Welcome` method should always be there for the bot to work. (This will change so you can use your own default method)

- (Removing any of these methods will stop the bot from receiving the respected update types) 🔽🔽:
  - appcallback.go contains `CallbackQueryHandler` method to handle callback queries.
  - appchannel.go contains `ChannelPostHandler` method to handle updates of channels. 
  - appchatmember.go contains `ChatMemberHandler` for handling updates regarding a chat member's change. 
  - appmychatmember.go containts `MyChatMemberHandler` handles updates related to the bot's changes.
  - appgroup.go containts `MessageTypeGroupHandler` which handles updates received from a group/supergroup.
  - apppollanswer.go containts `PollAnswerHandler` handling the updates of a user's choice in a non-anonymous poll created by the bot.

## `State` Concept
The key and magic component of the framework is this section! 

To develop an advanced Telegram bot that is offering users a wide range of services that need users' interactions, there are only a few methods you can use for your bot's structure. I've been developing bots for more than 6 years now and this method has been my special recipe that I'm gonna share with you!

For private chats your users are going to communicate to your bot through menus, most menus are made of button and a text asking users to pick a button to process their response or move them to another menu. These menus describe the state of a user and each menu is going to have its own method inside your struct (default menus are put inside `appprivate.go` file). So if your user's state is `Welcome` it means that they are in `Welcome` menu of your bot and the method `app.App{}.Welcome` is going to handle their requests. 

Your `State` handler methods should have this signature:

```go
func (app App) StateName(update *tgbotapi.Update, isSwitched bool) (newState string) {}
```
  - The framework is going to pass the updates of private chats based on their states to their respected handlers and inside the handlers you'll have access to the update using the first argument: `update`.
  - The second argument `isSwitched` is for when you're sending a user to another menu without processing their last sent text. It's for the next menu to be aware that there are no new updates from the user and the handler should ignore checking for them.
  - There's going to be an output called `newState` for your handler in case you want send the user to another menu with `isSwitched` set to true by the framework. You can remove the output in case you don't want to switch their menu or just let it be without changing it's empty value. (`newState = ""`)

## StateStorage
The framework has a built-in state storage storing the states of users inside a map and of course it's not recommended because it's temorary.

You can set your own storage and this is done by implementing the `StateStorage` interface:
```go
type UserStateStorage interface {
	UserState(userId int64) (string, error)
	SetUserState(userId int64, state string) error
}
```

And passing to to the framework by using the 
`bot.SetUserStateStorage(yourTypeImplementingTheInterface)`. 

## Logging
The framework is using `logrus` to log internal errors.

## Why reflection?
The framework is using `reflect` package in the background to match states to their handlers as well as handling different types of updates. If you were to use `maps` you would have ended up with duplicate and more code. I think the trade-off seems logical in this case and using `reflect` is a necessary evil here.

## Contibution
To help building telegram bots easier and more fun using `go-telebot`, you can open an issue or contact me on Telegram [@ali_error](https://t.me/ali_error)

## Bots made by telebot
[Taskilibot](https://t.me/taskilibot): A bot to manage your tasks
[Linkulubot](https://t.me/linkulubot): A bot to convert telegram videos to direct download links
