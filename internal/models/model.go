package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	OrderUID string `json:"order_uid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	OrderUID     string `json:"order_uid"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	OrderUID    string `json:"order_uid"`
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (o *Order) Validate() error {
	if strings.TrimSpace(o.OrderUID) == "" {
		return errors.New("order_uid не должен быть пустым")
	}
	if strings.TrimSpace(o.TrackNumber) == "" {
		return errors.New("track_number не должен быть пустым")
	}
	if o.DateCreated.IsZero() {
		return errors.New("date_created обязателен")
	}
	if len(o.Items) == 0 {
		return errors.New("заказ должен содержать хотя бы один item")
	}
	if err := o.Delivery.Validate(); err != nil {
		return fmt.Errorf("delivery: %e", err)
	}
	if err := o.Payment.Validate(); err != nil {
		return fmt.Errorf("payment: %e", err)
	}
	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item[%d]: %e", i, err)
		}
	}
	return nil
}

func (d *Delivery) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("delivery.name не должен быть пустым")
	}
	if strings.TrimSpace(d.Phone) == "" {
		return errors.New("delivery.phone не должен быть пустым")
	}
	return nil
}

func (p *Payment) Validate() error {
	if strings.TrimSpace(p.Transaction) == "" {
		return errors.New("payment.transaction не должен быть пустым")
	}
	if p.Amount < 0 {
		return errors.New("payment.amount не может быть отрицательным")
	}
	return nil
}

func (i *Item) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("item.name не должен быть пустым")
	}
	if i.Price < 0 {
		return errors.New("item.price не может быть отрицательным")
	}
	if i.TotalPrice < 0 {
		return errors.New("item.total_price не может быть отрицательным")
	}
	return nil
}
