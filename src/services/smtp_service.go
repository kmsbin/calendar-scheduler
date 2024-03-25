package services

type SmptInterface interface {
	SendEmail(email string) error
}
