package mailer

import "embed"

const (
	FromName            = "Welcome to Gobali Where You Can Rent A Good Villa !"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(TemplateFile, username, email string, data any, isSandbox bool) (int, error)
}
