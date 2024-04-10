package services

type SmptInterface interface {
	SendEmail(EmailData) error
}

type EmailData struct {
	Email   string
	BaseUrl string
	Code    string
}
