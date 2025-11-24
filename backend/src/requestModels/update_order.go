package requestModels

type UpdateOrderRequest struct {
	OrderID				uint					`json:"order_id"`
	DeliveryAddress		string					`json:"delivery_address"`
	DeliveryLatitude	float64					`json:"delivery_latitude"`	
	DeliveryLongitude	float64					`json:"delivery_longitude"`
}

