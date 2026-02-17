package app

import (
	"micro-warehouse/notification-service/controller"
	"micro-warehouse/notification-service/pkg/email"
	"micro-warehouse/notification-service/pkg/rabbitmq"
	"micro-warehouse/notification-service/usecase"
)

type Container struct {
	EmailController *controller.EmailController
	EmailUsecase    *usecase.EmailUsecase
	RabbitMQService rabbitmq.RabbitMQServiceInterface
	EmailService    email.EmailServiceInterface
}

func BuildContainer(rabbitMQService rabbitmq.RabbitMQServiceInterface, emailService email.EmailServiceInterface) *Container {
	emailUsecase := usecase.NewEmailUsecase(emailService)
	emailController := controller.NewEmailController(emailUsecase)
	
	return &Container{
		EmailController: emailController,
		EmailUsecase:    emailUsecase,
		RabbitMQService: rabbitMQService,
		EmailService:    emailService,
	}
}