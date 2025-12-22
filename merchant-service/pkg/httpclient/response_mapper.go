package httpclient

import "micro-warehouse/merchant-service/controller/response"

func MapProductResponseToMerchantProduct(product *ProductResponse) response.MerchantProduct {
	return response.MerchantProduct{
		ID:                   product.ID,
		ProductID:            product.ID,
		ProductName:          product.Name,
		ProductAbout:         product.About,
		ProductPhoto:         product.Thumbnail,
		ProductPrice:         int(product.Price),
		ProductCategory:      product.Category.Name,
		ProductCategoryPhoto: product.Category.Photo,
	}
}

func MapWarehouseResponseToMerchantProduct(warehouse *WarehouseResponse) response.MerchantProduct {
	return response.MerchantProduct{
		WarehouseID:    warehouse.ID,
		WarehouseName:  warehouse.Name,
		WarehousePhoto: warehouse.Photo,
		WarehousePhone: warehouse.Phone,
	}
}
