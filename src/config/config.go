package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"utils"

	ini "gopkg.in/ini.v1"
)

type Config struct {
	NotifyType []string // 多类型,逗号分隔
	WebHook    WebHook
	MailServer MailServer
	MailUser   MailUser
	Slack      Slack
	Events     []string
	Shell      Shell
}

type Shell struct {
	Command string
}

type WebHook struct {
	Url string
}

type Slack struct {
	WebHookUrl string
	Channel    string
}

// 邮件服务器
type MailServer struct {
	User     string
	Password string
	Host     string
	Port     int
}

// 接收邮件的用户
type MailUser struct {
	Email []string
}

func ParseConfig() *Config {
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/supervisor-event-listener.ini", "config file")
	flag.Parse()
	configFile = strings.TrimSpace(configFile)
	if configFile == "" {
		Exit("请指定配置文件路径")
	}
	file, err := ini.Load(configFile)
	if err != nil {
		Exit("读取配置文件失败#" + err.Error())
	}

	config := &Config{}

	section := file.Section("default")
	notifyType := section.Key("notify_type").Strings(",")
	for i := range notifyType {
		notifyType[i] = strings.TrimSpace(notifyType[i])
		if !utils.InStringSlice([]string{"mail", "slack", "webhook", "shell"}, notifyType[i]) {
			Exit("不支持的通知类型-" + notifyType[i])
		}

		switch notifyType[i] {
		case "mail":
			config.MailServer = parseMailServer(section)
			config.MailUser = parseMailUser(section)
		case "slack":
			config.Slack = parseSlack(section)
		case "webhook":
			config.WebHook = parseWebHook(section)
		case "shell":
			config.Shell = parseShell(section)
		}

	}

	config.NotifyType = notifyType

	events := section.Key("events").Strings(",")
	for i := range events {
		events[i] = strings.TrimSpace(events[i])
	}

	if len(events) == 0 {
		Exit("监听事件未配置")
	}

	config.Events = events
	return config
}

func parseShell(section *ini.Section) Shell {
	s := Shell{}
	s.Command = section.Key("shell.command").String()

	return s
}

func parseMailServer(section *ini.Section) MailServer {
	user := section.Key("mail.server.user").String()
	user = strings.TrimSpace(user)
	password := section.Key("mail.server.password").String()
	password = strings.TrimSpace(password)
	host := section.Key("mail.server.host").String()
	host = strings.TrimSpace(host)
	port, portErr := section.Key("mail.server.port").Int()
	if user == "" || password == "" || host == "" || portErr != nil {
		Exit("邮件服务器配置错误")
	}

	mailServer := MailServer{}
	mailServer.User = user
	mailServer.Password = password
	mailServer.Host = host
	mailServer.Port = port

	return mailServer
}

func parseMailUser(section *ini.Section) MailUser {
	user := section.Key("mail.user").String()
	user = strings.TrimSpace(user)
	if user == "" {
		Exit("邮件收件人配置错误")
	}
	mailUser := MailUser{}
	mailUser.Email = strings.Split(user, ",")

	return mailUser
}

func parseSlack(section *ini.Section) Slack {
	webHookUrl := section.Key("slack.webhook_url").String()
	webHookUrl = strings.TrimSpace(webHookUrl)
	channel := section.Key("slack.channel").String()
	channel = strings.TrimSpace(channel)
	if webHookUrl == "" || channel == "" {
		Exit("Slack配置错误")
	}

	slack := Slack{}
	slack.WebHookUrl = webHookUrl
	slack.Channel = channel

	return slack
}

func parseWebHook(section *ini.Section) WebHook {
	url := section.Key("webhook_url").String()
	url = strings.TrimSpace(url)
	if url == "" {
		Exit("WebHookUrl配置错误")
	}
	webHook := WebHook{}
	webHook.Url = url

	return webHook
}

func Exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
