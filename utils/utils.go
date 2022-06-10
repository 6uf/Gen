package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"github.com/marcusolsson/tui-go"
)

var Method = regexp.MustCompile(`^<(.+)> (.*)$`)

func AppendToChatBox(V *tui.Box, Text string) {
	V.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format("15:04")+" "+Text),
		tui.NewSpacer(),
	))
}

func GetDataForText(Name string) (Data Sender) {
	if Context := Method.FindAllStringSubmatch(Name, -1); len(Context) == 1 && len(Context[0]) > 1 {
		return Sender{
			Sender: Context[0][1],
			Output: Context[0][2],
		}
	}
	return Sender{
		Sender: "Server",
		Output: Name,
	}
}

func (Data *Session) SetUpBasicConn() error {
	history := Data.Context.History
	Data.Conn = bot.NewClient()

	Data.Server = F.Server
	P := basic.NewPlayer(Data.Conn, basic.DefaultSettings)
	Data.Conn.Auth = bot.Auth{
		AsTk: Data.Context.Acc.Bearer,
		Name: Data.Context.Acc.Info.Name,
		UUID: Data.Context.Acc.Info.ID,
	}
	basic.EventsListener{
		ChatMsg: func(c chat.Message, pos byte, uuid uuid.UUID) error {
			if Profile := GetDataForText(c.ClearString()); !(Profile.Sender == Data.Conn.Name) {
				AppendToChatBox(history, fmt.Sprintf("<%v> %v", Profile.Sender, Profile.Output))
				Data.Context.Ui.Repaint()
			}
			return nil
		},
		Disconnect: func(reason chat.Message) error {
			AppendToChatBox(history, fmt.Sprintf("<%s> You have disconnected for: %v", "Genocide", reason.ClearString()))
			Data.Context.Ui.Repaint()
			return nil
		},
		Death: func() error {
			return P.Respawn()
		},
	}.Attach(Data.Conn)

	if err := Data.Conn.JoinServer(Data.Server); err != nil {
		return err
	}
	go Data.Conn.HandleGame()
	return nil
}

func (s *LoginData) CheckAndUpdate(Text string, e *tui.Entry, history *tui.Box) Value {
	e.SetText("")
	switch true {
	case strings.Contains(Text, "-close"):
		return Value{
			Type: CloseType,
		}
	case strings.Contains(Text, "-load"):
		return Value{
			Type:    PreLoadAccount,
			Content: strings.Trim(strings.Split(Text, "-load")[1], " "),
		}
	case strings.Contains(Text, "-help"):
		AppendToChatBox(history, "<Help> -server <blockmania.com> (New value for the server address you wanna join.)")
		AppendToChatBox(history, "<Help> -new account <-email, -password, -savename, -login> (Sub commands to save acc info)")
		AppendToChatBox(history, "<Help> -load <name> (This loads a acc from the config)")
		return Value{Type: AddedAccDetails}
	case strings.Contains(Text, "-new account"):
		AppendToChatBox(history, fmt.Sprintf("<%s> Please type your email, password then server for your next 3 messages!", "Genocide"))
		AppendToChatBox(history, fmt.Sprintf("<%s> -email test@gmail.com", "Genocide"))
		AppendToChatBox(history, fmt.Sprintf("<%s> -password password123", "Genocide"))
		AppendToChatBox(history, fmt.Sprintf("<%s> -savename newacc1", "Genocide"))
		AppendToChatBox(history, fmt.Sprintf("<%s> -login", "Genocide"))
		LoginInit = true
		return Value{Type: AddedAccDetails}
	case strings.Contains(Text, "-server"):
		Old := s.Server
		if Old == "" {
			Old = "None"
		}
		s.Server = strings.Trim(strings.Split(Text, "-server")[1], " ")
		AppendToChatBox(history, fmt.Sprintf("<Genocide> Changed server %v to %v", Old, s.Server))
		return Value{Type: AddedAccDetails}
	case strings.Contains(Text, "-email") && LoginInit:
		s.Email = strings.Trim(strings.Split(Text, "-email")[1], " ")
		return Value{Type: AddedAccDetails}
	case strings.Contains(Text, "-password") && LoginInit:
		s.Password = strings.Trim(strings.Split(Text, "-password")[1], " ")
		return Value{Type: AddedAccDetails}
	case strings.Contains(Text, "-savename") && LoginInit:
		s.SaveName = strings.Trim(strings.Split(Text, "-savename")[1], " ")
		return Value{Type: AddedAccDetails}
	case strings.Contains(Text, "-login") && LoginInit:
		LoginInit = false
		return Value{
			Type:  LoginType,
			Login: true,
		}
	}
	return Value{
		Type:    ChatMessage,
		Content: Text,
	}
}
