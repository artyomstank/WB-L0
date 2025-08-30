package models

//split
import "time"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Items []Item

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             Items     `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type OrderResponse struct {
	OrderUID        string          `json:"order_uid"`
	TrackNumber     string          `json:"track_number"`
	Delivery        Delivery        `json:"delivery"`
	Payment         PaymentResponse `json:"payment"`
	Items           ItemsResponse   `json:"items"`
	Locale          string          `json:"locale"`
	DeliveryService string          `json:"delivery_service"`
	DateCreated     time.Time       `json:"date_created"`
}

type PaymentResponse struct {
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type ItemsResponse []ItemResponse

type ItemResponse struct {
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	Brand       string `json:"brand"`
}

// convertToOrderResponse конвертирует Order в OrderResponse для отпраки модели пользователю
func (o *Order) ConvertToOrderResponse() *OrderResponse {
	itemsResponse := make(ItemsResponse, len(o.Items))
	for i, item := range o.Items {
		itemsResponse[i] = ItemResponse{
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			Brand:       item.Brand,
		}
	}

	return &OrderResponse{
		OrderUID:    o.OrderUID,
		TrackNumber: o.TrackNumber,
		Delivery:    o.Delivery,
		Payment: PaymentResponse{
			Currency:     o.Payment.Currency,
			Provider:     o.Payment.Provider,
			Amount:       o.Payment.Amount,
			PaymentDt:    o.Payment.PaymentDt,
			Bank:         o.Payment.Bank,
			DeliveryCost: o.Payment.DeliveryCost,
			GoodsTotal:   o.Payment.GoodsTotal,
			CustomFee:    o.Payment.CustomFee,
		},
		Items:           itemsResponse,
		Locale:          o.Locale,
		DeliveryService: o.DeliveryService,
		DateCreated:     o.DateCreated,
	}
}
