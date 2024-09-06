package bot

import (
	"encoding/json"
	"goBot/goUnits/logger/logger"
	"os"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

var Bot *tgbotapi.BotAPI
var Err error

type Config struct {
	Token    string `json:"token"`
	Loglevel int    `json:"loglevel"`
}

func CheckConfigFile() bool {
	cfginfo, err := os.Stat("config.json")
	if err != nil {
		// 如果文件不存在或无法读取，创建新的 config.json 文件
		logger.Error("Error in reading config! \n Reason:%s", err)
		logger.Debug("config info:%v \n Trying to create config", cfginfo)
		time.Sleep(2 * time.Second)
		// 创建新的配置
		newconfig := Config{
			Token:    "YOUR_BOT_TOKEN",
			Loglevel: 1,
		}

		// 将配置转换为 JSON 格式
		data, err := json.MarshalIndent(newconfig, "", "  ")
		if err != nil {
			logger.Error("Error marshaling config: %s", err)
			return false
		}

		// 写入文件
		err = os.WriteFile("config.json", data, 0755)
		if err != nil {
			logger.Error("Error writing config file: %s", err)
			return false
		}
		logger.Info("Successfully created config.json\n Please edit your config")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	} else {
		// 文件存在
		return true
	}
	return false
}
func GetToken(file string) (token string) {
	if CheckConfigFile() {
		configFile, err := os.Open(file)
		if err != nil {
			logger.Error("Error in reading config! \n Reason:%s", err)
		}
		defer configFile.Close()
		var config Config
		decoder := json.NewDecoder(configFile)
		if err := decoder.Decode(&config); err != nil {
			logger.Error("Error decoding config file: %s", err)
			return ""
		}

		if err != nil {
			logger.Error("%s", err)
		}
		token = config.Token
		return token
	} else {
		logger.Error("Fail to read config, please check it")
		return
	}

}
func InitBot(file string) {
	token := GetToken(file)
	Bot, Err = tgbotapi.NewBotAPI(token)
	if Err != nil {
		logger.Error("Failed to create Telegram bot: %v", Err)
	}
	logger.Info("Authorized on account %s", Bot.Self.UserName)
}

func VerifiedUser(uid, gid int64, gname string) bool {
	USerconfig := tgbotapi.ChatConfigWithUser{
		ChatID: gid,
		UserID: uid,
	}
	chatMenberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: USerconfig,
	}

	getChatMenber, err := Bot.GetChatMember(chatMenberConfig)

	if getChatMenber.Status == "creator" || getChatMenber.Status == "administrator" {
		return true
	}

	if err != nil {

		logger.Error("Get chat error: %s \n ChatId: %d \n UserId : %d \n gname : %s \n", err, gid, uid, gname)
		//fmt.Printf("chatmenber :%s", getChatMenber)
	}
	return false
}

func GetLogLevel(file string) (LogLevel int) {

	if CheckConfigFile() {
		configFile, err := os.Open(file)
		if err != nil {
			logger.Error("Error in reading config! \n Reason:%s", err)
		}
		defer configFile.Close()
		var config Config
		decoder := json.NewDecoder(configFile)
		if err := decoder.Decode(&config); err != nil {
			logger.Error("Error decoding config file: %s", err)

		}

		if err != nil {
			logger.Error("%s", err)
		}
		LogLevel = config.Loglevel
		return LogLevel
	} else {
		logger.Error("Fail to read config, please check it")
	}
	return 1
}
