package alerts

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EmailNotifier struct {
	enabled       bool
	smtpHost      string
	smtpPort      int
	smtpUsername  string
	smtpPassword  string
	emailFrom     string
	emailTo       string
	subjectPrefix string
	cooldown      time.Duration

	mu       sync.Mutex
	lastSent map[string]time.Time
}

func NewEmailNotifierFromEnv() *EmailNotifier {
	n := &EmailNotifier{
		smtpHost:      strings.TrimSpace(os.Getenv("ALERT_SMTP_HOST")),
		smtpPort:      parsePort(os.Getenv("ALERT_SMTP_PORT"), 587),
		smtpUsername:  strings.TrimSpace(os.Getenv("ALERT_SMTP_USERNAME")),
		smtpPassword:  strings.TrimSpace(os.Getenv("ALERT_SMTP_PASSWORD")),
		emailFrom:     strings.TrimSpace(os.Getenv("ALERT_EMAIL_FROM")),
		emailTo:       strings.TrimSpace(os.Getenv("ALERT_EMAIL_TO")),
		subjectPrefix: strings.TrimSpace(os.Getenv("ALERT_SUBJECT_PREFIX")),
		cooldown:      parseCooldown(os.Getenv("ALERT_COOLDOWN_MINUTES"), 60),
		lastSent:      make(map[string]time.Time),
	}

	if n.subjectPrefix == "" {
		n.subjectPrefix = "[go-tene.life]"
	}
	if n.emailFrom == "" {
		n.emailFrom = n.smtpUsername
	}

	// SMTP host + receiver + sender are mandatory for sending.
	n.enabled = n.smtpHost != "" && n.emailTo != "" && n.emailFrom != ""
	if !n.enabled {
		log.Println("email alerts disabled: set ALERT_SMTP_HOST, ALERT_EMAIL_FROM and ALERT_EMAIL_TO to enable")
	}
	return n
}

func (n *EmailNotifier) Notify(key, subject, body string) {
	if n == nil || !n.enabled {
		return
	}
	if !n.canSend(key) {
		return
	}
	if err := n.send(subject, body); err != nil {
		log.Printf("email alert send failed: %v", err)
		return
	}
	log.Printf("email alert sent: key=%s to=%s", key, n.emailTo)
}

func (n *EmailNotifier) canSend(key string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now()
	last, ok := n.lastSent[key]
	if ok && now.Sub(last) < n.cooldown {
		return false
	}
	n.lastSent[key] = now
	return true
}

func (n *EmailNotifier) send(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", n.smtpHost, n.smtpPort)
	var auth smtp.Auth
	if n.smtpUsername != "" || n.smtpPassword != "" {
		auth = smtp.PlainAuth("", n.smtpUsername, n.smtpPassword, n.smtpHost)
	}

	subject = strings.TrimSpace(subject)
	if subject == "" {
		subject = "Alert"
	}
	finalSubject := fmt.Sprintf("%s %s", n.subjectPrefix, subject)
	message := strings.Join([]string{
		fmt.Sprintf("From: %s", n.emailFrom),
		fmt.Sprintf("To: %s", n.emailTo),
		fmt.Sprintf("Subject: %s", finalSubject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
		"",
	}, "\r\n")

	return smtp.SendMail(addr, auth, n.emailFrom, []string{n.emailTo}, []byte(message))
}

func parsePort(raw string, fallback int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port <= 0 {
		return fallback
	}
	return port
}

func parseCooldown(raw string, fallbackMinutes int) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Duration(fallbackMinutes) * time.Minute
	}
	minutes, err := strconv.Atoi(raw)
	if err != nil || minutes <= 0 {
		return time.Duration(fallbackMinutes) * time.Minute
	}
	return time.Duration(minutes) * time.Minute
}
