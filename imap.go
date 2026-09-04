package main

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/emersion/go-imap/v2/imapclient"
)

type session struct {
	imap *imapclient.Client
	acc  *Account
}

func dial_imap(acc *Account) (*imapclient.Client, error) {
	addr := acc.IMAPHost + ":" + strconv.Itoa(acc.IMAPPort)
	tls_cfg := &tls.Config{ServerName: acc.IMAPHost}

	if acc.TLSSkipVerify {
		tls_cfg.InsecureSkipVerify = true
	}

	var c *imapclient.Client
	var err error
	switch acc.IMAPTLSMode {
	case "implicit":
		c, err = imapclient.DialTLS(addr, &imapclient.Options{TLSConfig: tls_cfg})
	case "starttls", "":
		c, err = imapclient.DialStartTLS(addr, &imapclient.Options{TLSConfig: tls_cfg})
	default:
		return nil, fmt.Errorf("unknown imap_tls_mode %q", acc.IMAPTLSMode)
	}
	if err != nil {
		return nil, fmt.Errorf("imap dial %s: %w", addr, err)
	}

	if err := c.WaitGreeting(); err != nil {
		c.Close()
		return nil, fmt.Errorf("imap greeting: %w", err)
	}
	return c, nil
}

func (s *session) login(password string) error {
	if err := s.imap.Login(s.acc.User, password).Wait(); err != nil {
		return fmt.Errorf("imap login: %w", err)
	}
	return nil
}

func (s *session) mailboxes() ([]string, error) {
	items, err := s.imap.List("", "*", nil).Collect()
	if err != nil {
		return nil, fmt.Errorf("imap list: %w", err)
	}
	names := make([]string, 0, len(items))
	for _, it := range items {
		names = append(names, it.Mailbox)
	}
	return names, nil
}
