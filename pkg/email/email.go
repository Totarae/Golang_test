package email

type EmailService interface {
	SendConfirmation(email string, orderID string) error
}
