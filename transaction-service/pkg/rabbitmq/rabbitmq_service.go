package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type StockReductionEvent struct {
	WarehouseID uint      `json:"warehouse_id"`
	ProductID   uint      `json:"product_id"`
	Stock       int       `json:"stock"`
	MerchantID  uint      `json:"merchant_id"`
	Timestamp   time.Time `json:"timestamp"`
}

const (
	ExchangeName = "warehouse_events"
	QueueName    = "stock_reduction_queue"
	RoutingKey   = "stock_reduction"
)

func NewRabbitMQService(rabbitMQUrl string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 1: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 2: %v", err)
		return nil, err
	}

	err = ch.ExchangeDeclare(
		ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 3: %v", err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 4: %v", err)
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 5: %v", err)
		return nil, err
	}

	return &RabbitMQService{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQService) PublishStockReductionEvent(ctx context.Context, event StockReductionEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		log.Errorf("[RabbitMQService] PublishStockReductionEvent - 1: %v", err)
		return err
	}

	err = r.ch.Publish(
		ExchangeName,
		RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Errorf("[RabbitMQService] PublishStockReductionEvent - 2: %v", err)
		return err
	}

	return nil
}

func (r *RabbitMQService) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}	
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}