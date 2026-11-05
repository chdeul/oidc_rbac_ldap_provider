package interf

import (
	"fmt"
	"main/tools"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// Prerequis : lancer un serveur smtp mailpit

func TestMail(t *testing.T) {

	svrConfig := tools.SmtpCfg{Host: "localhost", Port: "1025", Ident: "mailpit", Password: "chdeul"}

	SetupSmtp(&svrConfig)

	from := "chdeul@localhost.loc"

	subject := EncodeSubject("TEST: Vous avez demandé un Test de notification par mail")
	body := "<p>Bonjour,</p>\r\n"
	body += "<p>Vous avez demandé un Test de notification par mail<br/>"
	body += "<p>A bientôt !</p>\r\n"

	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", "test@test.com")
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0;\r\n"
	msg += "Content-Transfer-Encoding: quoted-printable\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg += fmt.Sprint(body)

	err := MySendMail(from, "test@test.com", []byte(msg))
	tools.TEcheck(t, err, "Mail non envoyé : erreur")
}
