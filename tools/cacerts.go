package tools

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"golang.org/x/oauth2"
)

func LoadTrustedCertFiles(ctx context.Context) context.Context {

	// Load https server certificate
	caCert, err := os.ReadFile(Cfg.Openidc.Client.RootCA)
	if err != nil {
		Panic("Fichier certificat https non chargé : %v", err)
	}

	// Create certificates store
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Config http client with this store
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool, /* InsecureSkipVerify: true, */
		},
	}
	httphelper.DefaultHTTPClient.Transport = customTransport
	customClient := &http.Client{Transport: customTransport}

	// Return this client context
	return context.WithValue(ctx, oauth2.HTTPClient, customClient)
}
