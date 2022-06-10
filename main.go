package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"main/utils"

	"github.com/6uf/apiGO"
	"github.com/Lambels/cronjob"
	"github.com/Tnze/go-mc/chat"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/marcusolsson/tui-go"
)

func init() {
	utils.C.LoadState()
}

func main() {
	Data := &utils.Session{}

	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	historyBox.SetTitle(" Genocide ")

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(
		tui.Expanding,
		tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(
		tui.Expanding,
		tui.Maximum)
	chats := tui.NewVBox(
		historyBox,
		inputBox)
	chats.SetSizePolicy(
		tui.Expanding,
		tui.Expanding)

	ui, err := tui.New(tui.NewHBox(chats))
	if err != nil {
		log.Fatal(err)
	}

	utils.AppendToChatBox(history, fmt.Sprintf("<%s> Welcome to the chat!", "Genocide"))
	ui.Repaint()

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	input.OnSubmit(func(e *tui.Entry) {
		if Msg := chat.Text(e.Text()).String(); Msg != "" {
			if Context := utils.F.CheckAndUpdate(e.Text(), e, history); !reflect.DeepEqual(Context, nil) {
				switch Context.Type {
				case utils.CloseType:
					Data.Conn.Close()
					Data = &utils.Session{}
				case utils.RejoinType:
					Data.Conn.Close()
					Auth := Data.Conn.Auth
					Data = &utils.Session{
						Context: utils.FuncReq{
							History: history,
							Context: Context,
							Acc: apiGO.Info{
								Bearer: Auth.AsTk,
								Info: apiGO.UserINFO{
									ID:   Auth.UUID,
									Name: Auth.Name,
								},
							},
							Ui: ui,
						},
					}
					cron := cronjob.New()
					cron.AddFunc(
						Data.SetUpBasicConn,
						cronjob.In(cron.Now(), 1*time.Second),
					)
					cron.Run()
				case utils.ChatMessage:
					utils.AppendToChatBox(history, fmt.Sprintf("<%s> %v", Data.Conn.Name, Msg))
					ui.Repaint()
					if err := Data.Conn.Conn.WritePacket(pk.Marshal(
						0x03,
						pk.String(Msg),
					)); err != nil {
						utils.AppendToChatBox(history, fmt.Sprintf("<Genocide> Error: %v", err))
					}
				case utils.LoginType:
					for _, data := range apiGO.Auth([]string{fmt.Sprintf("%v:%v", utils.F.Email, utils.F.Password)}).Details {
						if data.AccountType == "Microsoft" && data.Bearer != "" && data.Error == "" {
							Data.Context = utils.FuncReq{
								History: history,
								Context: Context,
								Acc:     data,
								Ui:      ui,
							}
							utils.C.UploadAccount(utils.F.SaveName, data.Bearer, data.Info.ID, data.Info.Name)
							utils.AppendToChatBox(history, fmt.Sprintf("<Genocide> Succesfully authed %v", data.Email))
							Data.SetUpBasicConn()
							break
						} else {
							utils.AppendToChatBox(history, fmt.Sprintf("<Genocide> Unable to auth %v - %v", data.Email, data.Error))
						}
					}
				case utils.PreLoadAccount:
					if ConfigData := utils.C.GetValueFromConfig(Context.Content); ConfigData.Bearer != "" {
						Data.Context = utils.FuncReq{
							History: history,
							Context: Context,
							Acc: apiGO.Info{
								Bearer: ConfigData.Bearer,
								Info:   ConfigData.Info,
							},
							Ui: ui,
						}
						utils.AppendToChatBox(history, fmt.Sprintf("<Genocide> Succesfully loaded %v", ConfigData.Info.Name))
						if err := Data.SetUpBasicConn(); err != nil {
							utils.AppendToChatBox(history, fmt.Sprintf("<Genocide> Error: %v", err))
						}
					}
				}
			}
		}
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
