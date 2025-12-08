package requestModels

type CancelOrderRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Reason  string `json:"reason"`
}
