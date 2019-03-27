package models



// PORT
type Configuration struct {
	MailServer   string
	MailPort     int
	MailUsername string
	MailPassword string
}

type NotificationPG struct {
	Model

	EmailTo string `form:"email"`
	Description string `form:"description"`
}


