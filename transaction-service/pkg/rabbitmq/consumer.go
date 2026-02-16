package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type StockReducedEvent struct {
	MerchantID uint                       `json:"merchant_id"`
	Products   []StockReducedEventProduct `json:"products"`
	OrderID    string                     `json:"order_id"`
	Timestamp  time.Time                  `json:"timestamp"`
}

type StockReducedEventProduct struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type StockConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewStockConsumer(url string) (*StockConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 1: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 2: %v", err)
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"bussiness_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 3: %v", err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"merchant_stock_events",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 4: %v", err)
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"merchant.stock.*",
		"bussiness_events",
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 5: %v", err)
		return nil, err
	}

	return &StockConsumer{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQService) PublishStockReducedEvent(ctx context.Context, event StockReducedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = r.ch.Publish(
		"bussiness_events",       // exchange
		"merchant.stock.reduced", // routing key - lebih spesifik
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		})

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (sc *StockConsumer) Close() error {
	if sc.ch != nil {
		sc.ch.Close()
	}
	if sc.conn != nil {
		return sc.conn.Close()
	}
	return nil
}
