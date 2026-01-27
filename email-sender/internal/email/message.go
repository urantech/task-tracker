package email

type Message struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
}

type UserRegisteredEvent struct {
	UserID int64  `json:"id"`
	Email  string `json:"email"`
}

type DailyReportEvent struct {
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}
