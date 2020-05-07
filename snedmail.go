package main

import (
	"flag"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)


var (
	to, from, pass, host, port, subject, body, format string
)

func init() {
	flag.StringVar(&to, "to", "", "send to user example: xx@xx.com;yy@yy.com")                     // 发送给谁，支持多个，以分号隔开
	flag.StringVar(&from, "from", "", "from user")                                                 // 发件人
	flag.StringVar(&pass, "pass", "", "user pass")                                                 // 发件人密码
	flag.StringVar(&host, "host", "", "mail server address")                                       // smtp服务地址
	flag.StringVar(&port, "port", "25", "mail server port")                                        // smtp端口号
	flag.StringVar(&subject, "subject", "Test", "mail's subject")                                  // 邮件主题
	flag.StringVar(&body, "body", "sendmail ok", "mail's body")                                    // 邮件正文
	flag.StringVar(&format, "fm", "text/html", "mail's body format example: text/html text/plain") // 邮件格式 text/html text/plain
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

//需要使用Login作为参数
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}

func SendEmail(subject, body string) error {
	// smtp.SendMail 参数to接收的是一个数组，格式化命令行传进来的收件人
	send_to := strings.Split(to, ";")
	content_type := fmt.Sprintf("Content-Type: %s; charset=UTF-8",format)
	msg := []byte("To: " + to + " \r\nFrom: " + from + " >\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	smtpServer := fmt.Sprintf("%s:%s",host,port)
	auth := LoginAuth(from, pass)
	err := smtp.SendMail(smtpServer, auth, from, send_to, msg)
	return err
}

func main() {
	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if err := SendEmail(subject, body); err != nil {
		fmt.Printf("send mail err:%s\n", err.Error())
	}
}
