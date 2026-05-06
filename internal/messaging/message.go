package messaging

import "encoding/json"

const (
	TaskExchange            = "task.events"
	UserExchange            = "user.events"
	CreditResultExchange    = "credit.results"
	CaptchaVerify           = "captcha.events"
	QueueAuditCaptchaEvents = "captcha.events.queue"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
