package response

import "micro-warehouse/transaction-service/pkg/pagination"

type TransactionResponse struct {
	ID                  uint                         `json:"id"`
	Name                string                       `json:"name"`
	Phone               string                       `json:"phone"`
	Email               string                       `json:"email"`
	Address             string                       `json:"address"`
	SubTotal            int64                        `json:"sub_total"`
	TaxTotal            int64                        `json:"tax_total"`
	GrandTotal          int64                        `json:"grand_total"`
	MerchantID          uint                         `json:"merchant_id"`
	MerchantName        string                       `json:"merchant_name"`
	PaymentStatus       string                       `json:"payment_status"`
	PaymentMethod       string                       `json:"payment_method"`
	TransactionCode     string                       `json:"transaction_code"`
	OrderID             string                       `json:"order_id"`
	Notes               string                       `json:"notes"`
	TransactionProducts []TransactionProductResponse `json:"transaction_products"`
}

type TransactionProductResponse struct {
	ID            uint   `json:"id"`
	ProductID     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductPhoto  string `json:"product_photo"`
	ProductAbout  string `json:"product_about"`
	Quantity      int64  `json:"quantity"`
	Price         int64  `json:"price"`
	SubTotal      int64  `json:"sub_total"`
	TransactionID uint   `json:"transaction_id"`
	Category      struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
	} `json:"category"`
}

type GetAllTransactionsResponse struct {
	Transactions []TransactionResponse         `json:"transactions"`
	Pagination   pagination.PaginationResponse `json:"pagination"`
}

type DashboardResponse struct {
	TotalRevenue      int64 `json:"total_revenue"`
	TotalTransactions int64 `json:"total_transactions"`
	ProductsSold      int64 `json:"products_sold"`
}

type MerchantSummary struct {
	ID                uint              `json:"id"`
	Name              string            `json:"name"`
	TotalRevenue      int64             `json:"total_revenue"`
	TotalTransactions int64             `json:"total_transactions"`
	ProductsSold      int64             `json:"products_sold"`
	Products          []MerchantProduct `json:"products,omitempty"`
}

type MerchantProduct struct {
	ID                   uint   `json:"id"`
	MerchantID           uint   `json:"merchant_id"`
	ProductID            uint   `json:"product_id"`
	ProductName          string `json:"product_name"`
	ProductAbout         string `json:"product_about"`
	ProductPhoto         string `json:"product_photo"`
	ProductPrice         int64  `json:"product_price"`
	ProductCategory      string `json:"product_category"`
	ProductCategoryPhoto string `json:"product_category_photo"`
	Stock                int64  `json:"stock"`
	WarehouseID          uint   `json:"warehouse_id"`
	WarehouseName        string `json:"warehouse_name"`
	WarehousePhoto       string `json:"warehouse_photo"`
	WarehousePhone       string `json:"warehouse_phone"`
}

type DashboardByMerchantResponse struct {
	DashboardResponse
	Merchant MerchantSummary `json:"merchant"`
}
