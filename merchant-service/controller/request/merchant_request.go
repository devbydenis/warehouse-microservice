package request

type CreateMerchantRequest struct {
	Name     string `json:"name" validate:"required"`
	KeeperID uint   `json:"keeper_id" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Photo    string `json:"photo" validate:"required"`
}

