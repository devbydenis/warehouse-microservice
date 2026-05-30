package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-warehouse/notification-service/configs"
	"micro-warehouse/notification-service/pkg/email"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type RabbitMQServiceInterface interface {
	ConsumeEmail(ctx context.Context, emailService email.EmailServiceInterface) error
	Close() error
}

type rabbitmqService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     configs.Config
}

// Close implements [RabbitMQServiceInterface].
func (r *rabbitmqService) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}

// ConsumeEmail implements [RabbitMQServiceInterface].
func (r *rabbitmqService) ConsumeEmail(ctx context.Context, emailService email.EmailServiceInterface) error {
	// Declare queue untuk memastikan queue exists sebelum consume
	// Ini mencegah error jika queue belum dibuat oleh producer
	_, err := r.channel.QueueDeclare(
		"email_queue", // queue name
		true,          // durable (survive server restart)
		false,         // auto-delete
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Errorf("[RabbitMQService] ConsumeEmail - failed to declare queue: %v", err)
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// auto-ack = false (manual acknowledgment)
	// kenapa? agar kita bisa kontrol kapan message dihapus dari queue
	// jika auto-ack=true, message langsung hilang meski proses gagal
	msgs, err := r.channel.Consume("email_queue", "", false, false, false, false, nil)
	if err != nil {
		log.Errorf("[RabbitMQService] ConsumeEmail - 1: %v", err)
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Infof("[RabbitMQService] ConsumeEmail: %v", "Email consumer context cancelled")
				return
			case msg := <-msgs:
				// validasi pesan agar tidak kosong
				if len(msg.Body) == 0 {
					log.Warnf("[RabbitMQService] ConsumeEmail - 1: Empty message received")
					msg.Nack(false, false)
					continue
				}

				// fix: check jika channel ditutup (ConsumerTag bakal kosong)
				// ketika channel ditutup, amqp menutup channel msgs dan mengirim zero-value Delivery berupa struct kosong bukan nil
				if msg.ConsumerTag == "" {
					log.Errorf("[RabbitMQService] ConsumeEmail - 2: channel closed")
					return
				}

				var emailPayload email.EmailPayload
				if err := json.Unmarshal(msg.Body, &emailPayload); err != nil {
					log.Errorf("[RabbitMQService] ConsumeEmail - 3: %v", err)
					// nack dengan requeue=false (discard) karena JSON invalid tidak akan pernah berhasil
					msg.Nack(false, false)
					continue
				}

				// proceed email based on type
				var err error
				switch emailPayload.Type {
				case "welcome", "welcome_email":
					err = emailService.SendWelcomeEmail(ctx, emailPayload)
				default:
					log.Errorf("[RabbitMQService] ConsumeEmail - 4: unknown email type: %s", emailPayload.Type)
					// nack dengan requeue=false (discard) karena type unknown tidak akan pernah berhasil
					msg.Nack(false, false)
					continue
				}

				if err != nil {
					log.Errorf("[RabbitMQService] ConsumeEmail - 5: %v", err)
					// nack dengan requeue=true agar bisa di-retry nanti
					msg.Nack(false, true)
				} else {
					log.Infof("[RabbitMQService] ConsumeEmail - 6: Email sent successfully to %s", emailPayload.Email)
					msg.Ack(false)
				}

			}
		}
	}()

	return nil
}

func NewRabbitMQService(config configs.Config) (RabbitMQServiceInterface, error) {
	conn, err := amqp.Dial(config.RabbitMQ.URL())
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 1: %v", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 2: %v", err)
		return nil, err
	}

	return &rabbitmqService{
		conn:    conn,
		channel: channel,
		cfg:     config,
	}, nil
}
