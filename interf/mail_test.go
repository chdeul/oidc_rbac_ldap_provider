package interf

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMail(t *testing.T) {

	SetupSmtp("localhost", "1025", "mailpit", "chdeul")

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
	if err != nil {
		t.Error("Mail non envoyé : erreur")
	}
}
