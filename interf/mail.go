package interf

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"main/tools"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Config lecture seule
var smtpHost string
var smtpPort string
var smtpIdent string
var smtpPassword string
var tlsConfig *tls.Config

// TODO Certificat propre éventuel à configurer

// Fait une seule fois à l'init de l'application
func SetupSmtp(host, port, ident, pass string) {

	// Configuré en yml
	smtpHost = host
	smtpPort = port
	smtpIdent = ident
	smtpPassword = pass

	// Charger le certificat racine (CA) dans un pool
	pool := x509.NewCertPool()
	caCert, err := os.ReadFile("/etc/pki/myca/rootCA.pem")
	if err != nil {
		tools.Fatal("Certificat racine non trouvé", err)
	}
	if !pool.AppendCertsFromPEM(caCert) {
		tools.Fatal("impossible d'ajouter le CA racine au pool")
	}
	caCert, err = os.ReadFile("/etc/pki/myca/intermediate.pem")
	if err != nil {
		tools.Fatal("Certificat intermédiaire non trouvé", err)
	}
	if !pool.AppendCertsFromPEM(caCert) {
		tools.Fatal("impossible d'ajouter le CA intermédiaire au pool")
	}
	// TLS config
	tlsConfig = &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
		MinVersion:         tls.VersionTLS12,
	}
}

func EncodeSubject(subject string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(subject))
	return fmt.Sprintf("=?UTF-8?B?%s?=", b64)
}

func MySendMail(from string, dest string, data []byte) error {

	// Connect using TLS
	rawConn, err := net.DialTimeout("tcp", smtpHost+":"+smtpPort, 10*time.Second)
	if err != nil {
		return err
	}

	// borne en durée TOUTES les lectures/écritures suivantes
	//if err := rawConn.SetDeadline(time.Now().Add(60 * time.Second)); err != nil {
	//	rawConn.Close()
	//	return err
	//}

	conn := rawConn
	if strings.HasSuffix(smtpPort, "465") {
		// SMTPS Implicite : port 465  ou 2465 ou 8465  ...
		tlsConn := tls.Client(rawConn, tlsConfig)
		if err = tlsConn.Handshake(); err != nil {
			rawConn.Close()
			return err
		}
		conn = tlsConn
	}
	// Check connection
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		conn.Close()
		return err
	}
	defer c.Close()

	// Vérifier que le serveur supporte STARTTLS
	if !strings.HasSuffix(smtpPort, "465") {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			err = c.StartTLS(tlsConfig)
			if err != nil {
				return err
			}
		}
	}
	// Quels mécanismes : PLAIN  ( LOGIN, CRAM, XOAUTH2 : non géré ! )
	if ok, mechs := c.Extension("AUTH"); ok {
		fmt.Println("AUTH annoncés:", mechs)
		if !strings.Contains(mechs, "PLAIN") {
			tools.Fatal("pas d'authentification PLAIN")
		}
	} else {
		tools.Fatal("pas d'extension AUTH")
	}

	// Auth SMTP : classic  Plain on StartTLS
	auth := smtp.PlainAuth("", smtpIdent, smtpPassword, smtpHost)
	if err = c.Auth(auth); err != nil {
		return err
	}

	// Send email
	if err = c.Mail(from); err != nil {
		return err
	}
	if err = c.Rcpt(dest); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}

	// Done
	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}
