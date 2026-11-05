package tools

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"golang.org/x/oauth2"
	"resty.dev/v3"
)

// Trust cert for http Context
func CtxLoadTrustedCertFiles(ctx context.Context, filepath string, insecure bool) context.Context {

	// Pas de root CA certificates
	if filepath == "" {
		return ctx
	}
	// Load https server certificate
	caCert, err := os.ReadFile(CfgFilename(filepath))
	CheckError(err, "Fichier certificat https non chargé : %v", err)

	// Create certificates store
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Config http client with this store
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: insecure,
		},
	}
	httphelper.DefaultHTTPClient.Transport = customTransport
	customClient := &http.Client{Transport: customTransport}

	// Return this client context
	return context.WithValue(ctx, oauth2.HTTPClient, customClient)
}

// Trust cert for Resty client
func RestyLoadTrustedCertFiles(client *resty.Client, filepath string, insecure bool) {

	// Pas de root CA certificates
	if filepath == "" {
		return
	}
	caCert, err := os.ReadFile(CfgFilename(filepath))
	CheckError(err, "Fichier certificat https non chargé : %v", err)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client.SetTLSClientConfig(&tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: insecure,
	})
}
