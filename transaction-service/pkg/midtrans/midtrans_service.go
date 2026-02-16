package midtrans

import (
	"micro-warehouse/transaction-service/configs"

	"github.com/gofiber/fiber/v2/log"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransServiceInterface interface {
	CreateTransaction(req CreateTransactionRequest) (*CreateTransactionResponse, error)
}

type TransactionItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
}

type CreateTransactionRequest struct {
	OrderID       string            `json:"order_id"`
	Amount        int64             `json:"amount"`
	Items         []TransactionItem `json:"items"`
	CustomerName  string            `json:"customer_name"`
	CustomerEmail string            `json:"customer_email"`
	CustomerPhone string            `json:"customer_phone"`
	Notes         string            `json:"notes"`
}

type CreateTransactionResponse struct {
	PaymentToken string `json:"payment_token"`
	OrderID      string `json:"order_id"`
}

type MidtransService struct {
	config *configs.Config
}

// CreateTransaction implements [MidtransServiceInterface].
func (m *MidtransService) CreateTransaction(req CreateTransactionRequest) (*CreateTransactionResponse, error) {
	midtrans.ServerKey = m.config.Midtrans.ServerKey
	if m.config.Midtrans.IsProduction {
		midtrans.Environment = midtrans.EnvironmentType(midtrans.Production)
	} else {
		midtrans.Environment = midtrans.EnvironmentType(midtrans.Sandbox)
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
		},
	}

	snapRes, err := snap.CreateTransaction(snapReq)
	if err != nil {
		log.Errorf("[MidtransService] CreateTransaction - 1: %v", err)
		return nil, err
	}

	return &CreateTransactionResponse{
		PaymentToken: snapRes.RedirectURL,
		OrderID:      req.OrderID,
	}, nil

}

func NewMidtransService(config *configs.Config) MidtransServiceInterface {
	return &MidtransService{config: config}
}
