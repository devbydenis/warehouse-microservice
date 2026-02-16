package request

type GetAllTransactionRequest struct {
	Page       int    `form:"page" query:"page" validate:"omitempty,min=1"`
	Limit      int    `form:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
	Search     string `form:"search" query:"search" validate:"omitempty"`
	SortBy     string `form:"sort_by" query:"sort_by" validate:"omitempty,oneof=id name created_at"`
	SortOrder  string `form:"sort_order" query:"sort_order" validate:"omitempty,oneof=asc desc"`
	MerchantID string `form:"merchant_id" query:"merchant_id" validate:"omitempty"`
}

type CreateTransactionRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Address    string `json:"address" validate:"required"`
	MerchantID uint   `json:"merchant_id" validate:"required"`
	Notes      string `json:"notes" validate:"omitempty"`
	Currency   string `json:"currency" validate:"omitempty,oneof=IDR"`
}

type CreateTransactionProductRequest struct {
	ProductID uint  `json:"product_id" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required,min=1"`
	Price     int64 `json:"price" validate:"required,min=1"`
}

type CreateTransactionWithProductRequest struct {
	CreateTransactionRequest
	Products []CreateTransactionProductRequest `json:"products" validate:"required,min=1,dive"`
}

type MidtransCallbackRequest struct {
	OrderID           string `json:"order_id" validate:"required"`
	TransactionID     string `json:"transaction_id" validate:"required"`
	TransactionStatus string `json:"transaction_status" validate:"required"`
	PaymentType       string `json:"payment_type" validate:"required"`
	FraudStatus       string `json:"fraud_status" validate:"required"`
	StatusCode        string `json:"status_code" validate:"required"`
	SignatureKey      string `json:"signature_key" validate:"required"`
}
