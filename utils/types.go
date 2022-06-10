package utils

import (
	"main/config"

	"github.com/6uf/apiGO"
	"github.com/Tnze/go-mc/bot"
	"github.com/marcusolsson/tui-go"
)

var (
	LoginInit bool
	C         config.Config
	F         LoginData
)

const (
	LoginType       = "login"
	CloseType       = "close"
	RejoinType      = "relog"
	ChatMessage     = "sendmsg"
	AddedAccDetails = "added"
	PreLoadAccount  = "preload"
)

type Session struct {
	Conn   *bot.Client
	Server string

	Context FuncReq
}

type FuncReq struct {
	History *tui.Box
	Context Value
	Acc     apiGO.Info
	Ui      tui.UI
}

type Sender struct {
	Sender string
	Output string
}

type Value struct {
	Type    string
	Content string
	Error   error
	Login   bool
	Body    []byte
}

type LoginData struct {
	Email, Password, Server, SaveName string
}
