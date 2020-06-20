package main

import (
	"./lib"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	TOKEN  = os.Getenv("TOKEN")
	db, ng = lib.SetupDB()
	help = `help Message
登録:  メッセージを登録します
削除:  登録メッセージから削除します
登録一覧  登録メッセージ一覧を表示します
restart_db  データベースを再読込します
ngadd:  NGワードに追加します
ngdel:  NGワードから削除します
realface  AIが生成した三次元の顔を送ります
oji:  おじさん
	`
)

func get_format_data() string {
	list := ""
	for key, value := range db.Msgs {
		list += key + ":" + value + "\n"
	}
	return list
}
func update_db() {
	for {
		time.Sleep(time.Second * 300)
		db, ng = lib.SetupDB()
	}
}

func main() {
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("error:start\n", err)
	}

	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		log.Fatal("error:wss\n", err)
	}
	go update_db()
	log.Println("Ready!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func stringInMap(s string, e map[string]string) bool {
	for k := range e {
		if k == s {
			return true
		}
	}
	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	switch m.Content {
	case "help":
                s.ChannelMessageSend(m.ChannelID, help)
	case "登録一覧":
		list := get_format_data()
		s.ChannelMessageSend(m.ChannelID, list)
		return

	case "restart_db":
		db, ng = lib.SetupDB()
		s.ChannelMessageSend(m.ChannelID, "更新しました")
		return

	case "realface":
		lib.Realface("disgord.jpeg")
		i, _ := os.Open("./img/disgord.jpeg")
		s.ChannelFileSend(m.ChannelID, "disgord.jpeg", i)
		return
        }
	if strings.Contains(m.Content, "help:") {
		_name := strings.Split(m.Content, ":")[1]
		msg := strings.Split(help, "\n")
		helpMsg := ""
		for _, s := range msg {
			if strings.Contains(s, _name) {
				helpMsg += s + "\n"
			}
		}
		if helpMsg == "" {
			s.ChannelMessageSend(m.ChannelID, "見つからなかったよ><")
			return
		}
		s.ChannelMessageSend(m.ChannelID, helpMsg)
        } else if strings.Contains(m.Content, "ngadd:") {
		if strings.Count(m.Content, ":") == 1 {
			s.ChannelMessageSend(m.ChannelID, "ちゃんと送ってね")
			return
		}
		_name := strings.SplitN(m.Content, ":", 3)
		if stringInMap(_name[1], db.Msgs) {
			s.ChannelMessageSend(m.ChannelID, "その言葉は既に反応リストに含まれているので登録できません")
			return
		}
		ng.Add_msg(_name[1], _name[2], "ngword")
		s.ChannelMessageSend(m.ChannelID, "登録成功しました")
	} else if strings.Contains(m.Content, "ngdel:") {
		_name := strings.SplitN(m.Content, ":", 2)
		if !stringInMap(_name[1], ng.Msgs) {
			ng.Delete_msg(_name[1], "ngword", "word")
			s.ChannelMessageSend(m.ChannelID, "削除に成功しました")
		} else {
			s.ChannelMessageSend(m.ChannelID, "登録されていません")
		}
	} else if strings.Contains(m.Content, "登録:") {
		if strings.Count(m.Content, ":") == 1 {
			s.ChannelMessageSend(m.ChannelID, "ちゃんと送ってね")
			return
		}
		_name := strings.SplitN(m.Content, ":", 3)
		if stringInMap(_name[1], ng.Msgs) {
			s.ChannelMessageSend(m.ChannelID, "その言葉はngワードに含まれているので登録できません")
			return
		}
		db.Add_msg(_name[1], _name[2], "msg")
		s.ChannelMessageSend(m.ChannelID, "["+_name[1]+"]"+"と言った時に["+_name[2]+"]と返って来るようにしました")

	} else if strings.Contains(m.Content, "削除:") {
		_name := strings.SplitN(m.Content, ":", 2)
		if stringInMap(_name[1], db.Msgs) {
			db.Delete_msg(_name[1], "msg", "come")
			s.ChannelMessageSend(m.ChannelID, "["+_name[1]+"]"+"と言った時に何も返ってこないようにしました")
		} else {
			s.ChannelMessageSend(m.ChannelID, "["+_name[1]+"]という言葉は登録されていません")
		}

	} else if strings.Contains(m.Content, "oji:") {
		_name := strings.SplitN(m.Content, ":", 2)
		oji, _ := lib.Ojichat(_name[1])
		s.ChannelMessageSend(m.ChannelID, oji)

	} else if stringInMap(m.Content, ng.Msgs) {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, "NGワードが検出されたので削除しました\n理由: "+ng.Msgs[m.Content])

	} else if stringInMap(m.Content, db.Msgs) {
		s.ChannelMessageSend(m.ChannelID, db.Msgs[m.Content])
	}
}
