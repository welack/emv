package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	acc  *Account
	pw   textinput.Model
	st   int // 0= login, 1 = mailboxes
	sess *session
	mbox []string
	err  error
}

type logged_in_msg struct {
	sess *session
	mbox []string
	err  error
}

func new_model(acc *Account) model {
	ti := textinput.New()
	ti.Placeholder = "password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'
	ti.Focus()
	return model{acc: acc, pw: ti}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logged_in_msg:
		if msg.err != nil {
			m.err, m.st = msg.err, 0
			return m, nil
		}
		m.sess, m.mbox, m.err, m.st = msg.sess, msg.mbox, nil, 1
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.sess != nil {
				m.sess.imap.Logout()
			}
			return m, tea.Quit
		case "enter":
			if m.st == 0 {
				m.err = nil
				return m, do_login(m.acc, m.pw.Value())
			}
		}
	}
	if m.st == 0 {
		var cmd tea.Cmd
		m.pw, cmd = m.pw.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	if m.st == 0 {
		b.WriteString("emv — login " + m.acc.User + "\n\n")
		b.WriteString(m.pw.View() + "\n\n")
		b.WriteString("(enter: login, esc: quit)\n")
	} else {
		b.WriteString("mailbox:\n\n")
		for _, n := range m.mbox {
			b.WriteString("  " + n + "\n")
		}
		b.WriteString("\n(esc: quit)\n")
	}
	if m.err != nil {
		b.WriteString("\nerror: " + m.err.Error() + "\n")
	}
	return b.String()
}

func do_login(acc *Account, password string) tea.Cmd {
	return func() tea.Msg {
		c, err := dial_imap(acc)
		if err != nil {
			return logged_in_msg{err: err}
		}
		s := &session{imap: c, acc: acc}
		if err := s.login(password); err != nil {
			c.Logout()
			return logged_in_msg{err: err}
		}
		mbox, err := s.mailboxes()
		return logged_in_msg{sess: s, mbox: mbox, err: err}
	}
}
