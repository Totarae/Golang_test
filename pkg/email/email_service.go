package email

import "fmt"

type SMTPEmailService struct {
	// SMTP config (адрес и прочее)
}

func NewSMTPEmailService( /* config params */ ) *SMTPEmailService {
	return &SMTPEmailService{
		// Конструктор сам
	}
}

func (s *SMTPEmailService) SendConfirmation(email string, orderID string) error {
	// Логика рассылки
	fmt.Printf("Sending email confirmation to %s for order %s\n", email, orderID)
	return nil
}
