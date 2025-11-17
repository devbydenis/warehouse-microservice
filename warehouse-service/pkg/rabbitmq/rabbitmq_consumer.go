package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"micro-warehouse/warehouse-service/repository"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	repo    repository.WarehouseProductRepositoryInterface
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
	RoutingKey   = "stock.reduction"
)

func NewRabbitMQConsumer(rabbitMQURL string, repo repository.WarehouseProductRepositoryInterface) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// declare exchange
	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// declare queue
	q, err := ch.QueueDeclare(
		QueueName, // name
		true,      //durable
		false,     // delete when unused
		false,     // exclusive
		false,     //no-wait
		nil,       // arguments
	)

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		RoutingKey,   // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		repo:    repo,
	}, nil
}

func (rc *RabbitMQConsumer) StartConsuming(ctx context.Context) error {
	msgs, err := rc.channel.Consume(
		QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func ()  {
		for {
			select {
			case <-ctx.Done():
				log.Infof("[RabbitMQConsumer] Stopping consumer due to context cancellation")
				return
			case msg := <-msgs:
				rc.handleMessage(ctx, msg)
			}
		}
	}()

	return nil
}

func (rc *RabbitMQConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	var event StockReductionEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Errorf("[RabbitMQConsumer] handleMessage - 1: %v", err)
		return
	}

	if err := rc.processStockReduction(ctx, event); err != nil {
		log.Errorf("[RabbitMQConsumer] handleMessage - 2: %v", err)
		return
	}

	msg.Ack(false)
}

func (rc *RabbitMQConsumer) processStockReduction(ctx context.Context, event StockReductionEvent) error {
	warehouseProduct, err := rc.repo.GetWarehouseProductByWarehouseIDAndProductID(ctx, event.WarehouseID, event.ProductID)
	if err != nil {
		log.Errorf("[RabbitMQConsumer] processStockReduction - 1: %v", err)
		return err
	}

	newStock := warehouseProduct.Stock - event.Stock
	if newStock < 0 {
		return errors.New("stock not enough")
	}

	warehouseProduct.Stock = newStock

	if err := rc.repo.UpdateWarehouseProduct(ctx, warehouseProduct); err != nil {
		log.Errorf("[RabbitMQConsumer] processStockReduction - 2: %v", err)
		return err
	}

	return nil
}