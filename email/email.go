package email

import (
    "crypto/tls"
    "fmt"
    "github.com/moonlightming/simple-docker-inside-webhook/conf"
    "log"
    "net/smtp"
    "strings"
)

var (
    // 配置读取
    config = conf.NewConfig()
    // 发送对象
    sendTo = strings.Split(config.Email.SendTo, ";")
    // 邮件格式
    contentType = "Content-Type: text/plain; charset=UTF-8"
    // 邮件服务器Host
    smtpHost = config.Email.SmtpHost + config.Email.SmtpPort
    // SMTP认证对象
    auth = smtp.PlainAuth(
        "",
        config.Email.UserEmail,
        config.Email.Password,
        config.Email.SmtpHost,
    )
)

func emailContent(message string) []byte {
    var contentBody = fmt.Sprintf(`
	此次构建已响应，可以直接登录网站查看效果：

	https://blog.moonlightming.top

	信息如下：

	%s

	`, message)

    subject := "Rebuild Message"
    return []byte("To: " + strings.Join(sendTo, ",") + "\r\nFrom: <" + config.Email.UserEmail + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + contentBody)
}

func sendToEmail(message string) error {

    if config.Email.Open == false {
        return nil
    }

    if err := smtp.SendMail(smtpHost, auth, config.Email.UserEmail, sendTo, emailContent(message)); err != nil {
        return err
    }

    return nil
}

func dial() (*smtp.Client, error) {

    conn, err := tls.Dial(
        "tcp",
        config.Email.SmtpHost+config.Email.SmtpPort,
        nil)

    if err != nil {
        return nil, err
    }

    return smtp.NewClient(conn, config.Email.SmtpHost)

}

func sendToEmailTLS(message string) error {

    if config.Email.Open == false {
        return nil
    }

    smtpClient, err := dial()
    if err != nil {
        return err
    }
    defer smtpClient.Close()

    if ok, _ := smtpClient.Extension("AUTH"); ok {
        if err := smtpClient.Auth(auth); err != nil {
            log.Println(err)
            return err
        }
    }

    if err := smtpClient.Mail(config.Email.UserEmail); err != nil {
        return err
    }

    for _, sendToAddr := range sendTo {
        if err := smtpClient.Rcpt(sendToAddr); err != nil {
            return err
        }
    }

    writer, err := smtpClient.Data()
    defer writer.Close()
    if _, err := writer.Write(emailContent(message)); err != nil {
        return err
    }
    return nil

}

func SendMail(message string) error {

    if config.Email.SmtpPort == ":25" {
        return sendToEmail(message)
    }
    return sendToEmailTLS(message)

}
