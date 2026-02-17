package usecase

import (
	"context"
	"micro-warehouse/notification-service/controller/request"
	"micro-warehouse/notification-service/pkg/email"
)

type EmailUsecase struct {
	emailService email.EmailServiceInterface
}

func NewEmailUsecase(emailService email.EmailServiceInterface) *EmailUsecase {
	return &EmailUsecase{
		emailService: emailService,
	}
}

func (e *EmailUsecase) SendEmail(ctx context.Context, req request.SendEmailRequest) error {
	return e.emailService.SendCustomEmail(ctx, req.To, req.Subject, req.Body)
}

func (e *EmailUsecase) SendWelcomeEmail(ctx context.Context, req request.SendWelcomeEmailRequest) error {
	payload := email.EmailPayload{
		Email:    req.Email,
		Password: req.Password,
		Type:     "welcome",
		UserId:   req.UserId,
		Name:     req.Name,
	}

	return e.emailService.SendWelcomeEmail(ctx, payload)
}
