swag-user:
	cd ./user-service && swag init -g main.go --parseDependency

swag-product:
	cd ./product-service && swag init -g main.go --parseDependency

swag-merchant:
	cd ./merchant-service && swag init -g main.go --parseDependency

swag-warehouse:
	cd ./warehouse-service && swag init -g main.go --parseDependency

swag-transaction:
	cd ./transaction-service && swag init -g main.go --parseDependency


swag:
	$(MAKE) swag-user
	$(MAKE) swag-product
	$(MAKE) swag-merchant
	$(MAKE) swag-warehouse
	$(MAKE) swag-transaction