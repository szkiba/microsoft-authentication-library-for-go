// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

func acquireTokenClientCertificate() {
	config := CreateConfig("confidential_config.json")

	pemData, err := os.ReadFile(config.PemData)
	if err != nil {
		log.Fatal(err)
	}

	// This extracts our public certificates and private key from the PEM file. If it is
	// encrypted, the second argument must be password to decode.
	certs, privateKey, err := confidential.CertFromPEM(pemData, "")
	if err != nil {
		log.Fatal(err)
	}

	// PEM files can have multiple certs. This is usually for certificate chaining where roots
	// sign to leafs. Useful for TLS, not for this use case.
	if len(certs) > 1 {
		log.Fatal("too many certificates in PEM file")
	}

	cred := confidential.NewCredFromCert(certs[0], privateKey)
	if err != nil {
		log.Fatal(err)
	}
	app, err := confidential.New(config.ClientID, cred, confidential.WithAuthority(config.Authority), confidential.WithCache(cacheAccessor))
	if err != nil {
		log.Fatal(err)
	}
	result, err := app.AcquireTokenSilent(context.Background(), config.Scopes)
	if err != nil {
		result, err = app.AcquireTokenByCredential(context.Background(), config.Scopes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Access Token Is " + result.AccessToken)
		return
	}
	fmt.Println("Silently acquired token " + result.AccessToken)
}
