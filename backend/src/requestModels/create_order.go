package requestModels

type OrderProductRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type AddOrderRequest struct {
	CustomerId          uint                  `json:"customer_id"`
	SellerId      		uint                	`json:"seller_id"`
	SellerAddress		string					`json:"seller_address"` 
	SellerLatitude		float64					`json:"seller_latitude"`
	SellerLongitude		float64					`json:"seller_longitude"`
	DeliveryAddress		string					`json:"delivery_address"`
	DeliveryLatitude	float64					`json:"delivery_latitude"`	
	DeliveryLongitude	float64					`json:"delivery_longitude"`
	Products         []OrderProductRequest `json:"products"` // Nested array of products
}