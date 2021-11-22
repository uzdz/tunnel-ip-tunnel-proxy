package gui

import (
	"fmt"
	"io"
	"ip-tunnel-proxy/pkg/config"
	"log"
	"os"
	"time"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

func Start() {
	var mw *walk.MainWindow
	var number *walk.Label
	var sStatus *walk.Label
	var port *walk.Label
	var version *walk.Label

	animal := Animal{Number: config.Number}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(3))
			if config.ProxyRun {
				sStatus.SetText(fmt.Sprintln(" 当前代理服务器状态：运行中..."))
			} else {
				sStatus.SetText(fmt.Sprintln(" 当前代理服务器状态：关闭"))
			}

			port.SetText(fmt.Sprintf(" 当前代理端口号：%d", config.Port))
		}
	}()

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "代理终端",
		MinSize:  Size{320, 240},
		Size:     Size{400, 600},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			PushButton{
				Text: "启动代理",
				OnClicked: func() {
					if config.ProxyRun {
						log.Println("当前服务器正在启动中，请先关闭服务器...")
					} else {
						go func() {
							config.CommandChan <- config.StartCommand
						}()
					}
				},
			},
			PushButton{
				Text: "关闭代理",
				OnClicked: func() {
					go func() {
						config.CommandChan <- config.ShutdownCommand
					}()
				},
			},
			PushButton{
				Text: "编辑代理服务器信息",
				OnClicked: func() {
					if cmd, err := RunAnimalDialog(mw, &animal); err != nil {
						log.Print(err)
					} else if cmd == walk.DlgCmdOK {
						number.SetText(fmt.Sprintf(" 设备唯一号：%s", animal.Number))
						config.Number = animal.Number
					}
				},
			},
			Label{
				AssignTo: &version,
				Text:     fmt.Sprintf(" 客户端版本号：%s", config.Version),
			},
			Label{
				AssignTo: &port,
				Text:     fmt.Sprintf(" 当前代理端口号：%s", ""),
			},
			Label{
				AssignTo: &sStatus,
				Text:     fmt.Sprintf(" 当前代理服务器状态：%s", "关闭"),
			},
			Label{
				AssignTo: &number,
				Text:     fmt.Sprintf(" 设备唯一号：%s", animal.Number),
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lv, err := NewLogView(mw)
	if err != nil {
		log.Fatal(err)
	}

	writers := []io.Writer{
		lv,
		os.Stderr,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)

	mw.Run()
}

type Animal struct {
	Number string
}

func RunAnimalDialog(owner walk.Form, animal *Animal) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         "编辑代理服务器信息",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "animal",
			DataSource:     animal,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{

					Label{
						Text: "设备唯一号:",
					},
					LineEdit{
						Text: Bind("Number"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
