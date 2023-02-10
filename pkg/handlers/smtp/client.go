
/*
Copyright 2020 VMWare

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
 This code is adapted from https://github.com/prometheus/alertmanager/blob/a75cd02786dfecd25e2469fc4df5d920e6b9c226/notify/email/email.go
*/

package smtp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/mkmik/multierror"
	"github.com/sirupsen/logrus"
)

func sendEmail(conf config.SMTP, msg string) error {
	ctx := context.Background()

	host, port, err := net.SplitHostPort(conf.Smarthost)
	if err != nil {
		return err
	}

	var (
		c       *smtp.Client
		conn    net.Conn
		success = false
	)

	tlsConfig := &tls.Config{}
	if port == "465" {

		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = host
		}

		conn, err = tls.Dial("tcp", conf.Smarthost, tlsConfig)
		if err != nil {
			return fmt.Errorf("establish TLS connection to server: %w", err)
		}
	} else {
		var (
			d   = net.Dialer{}
			err error
		)
		conn, err = d.DialContext(ctx, "tcp", conf.Smarthost)
		if err != nil {
			return fmt.Errorf("establish connection to server: %w", err)
		}
	}
	c, err = smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer func() {
		// Try to clean up after ourselves but don't log anything if something has failed.
		if err := c.Quit(); success && err != nil {
			logrus.Warnf("failed to close SMTP connection: %v", err)