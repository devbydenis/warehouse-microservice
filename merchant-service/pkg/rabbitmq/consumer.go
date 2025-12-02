package rabbitmq

import (
	"context"
	"encoding/json"
	"micro-warehouse/merchant-service/repository"
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
	merchantRepo repository.MerchantProductRepositoryInterface
}

func NewStockConsumer(url string, merchantRepo repository.MerchantProductRepositoryInterface) *StockConsumer {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 1: %v", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 2: %v", err)
		return nil
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
		return nil
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
		return nil
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
		return nil
	}

	return &StockConsumer{
		conn:         conn,
		ch:           ch,
		merchantRepo: merchantRepo,
	}
}

func (s *StockConsumer) ConsumeStockReductionEvents(ctx context.Context) error {
	msgs, err := s.ch.Consume(
		"merchant_stock_events",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] ConsumeStockReductionEvents - 1: %v", err)
		return err
	}

	for {
		select {
			case <-ctx.Done():
				log.Info("[StockConsumer] ConsumeStockReductionEvents - 2: Context done")
				return nil
			case msg := <-msgs:
				go s.handleStockReductionEvent(msg)
			}
	}
}

func (sc *StockConsumer) handleStockReductionEvent(msg amqp.Delivery) error {
	defer msg.Ack(false)

	var event StockReducedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Errorf("[StockConsumer] handleStockReductionEvent - 1: %v", err)
		return err
	}

	for _, product := range event.Products {
		if err := sc.reduceStock(event.MerchantID, product.ProductID, product.Quantity); err != nil {
			log.Errorf("[StockConsumer] handleStockReductionEvent - 2: %v", err)
			return err
		}

		log.Infof("Successfully reduce stock for product %d by %d", product.ProductID, product.Quantity)
	}

	return nil
}

func (sc *StockConsumer) reduceStock(merchantID uint, productID uint, quantity int) error {
	err := sc.merchantRepo.ReduceStock(context.Background(), merchantID, productID, int64(quantity))
	if err != nil {
		log.Errorf("[StockConsumer] reduceStock - 1: %v", err)
		return err
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

