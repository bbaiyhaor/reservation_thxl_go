package workflow

import (
	"fmt"
	"github.com/scorredoira/email"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"net/smtp"
	"os"
	"strings"
)

func SendEmail(subject string, body string, attached []string, to []string) error {
	appEnv := os.Getenv("RESERVATION_THXL_ENV")
	if !strings.EqualFold(appEnv, "ONLINE") {
		fmt.Printf("Send Email: \"%s\" to %s.\n", subject, strings.Join(to, ","))
		return nil
	}
	m := email.NewMessage(subject, body)
	m.From = utils.MAIL_USERNAME
	m.To = to
	for _, file := range attached {
		m.Attach(file)
	}
	return email.Send(fmt.Sprintf("%s:25", utils.MAIL_SMTP),
		smtp.PlainAuth("", utils.MAIL_USERNAME, utils.MAIL_PASSWORD, utils.MAIL_SMTP), m)
}
