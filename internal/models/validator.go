package models

import (
	"fmt"
	"regexp"
)

var (
	phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func (o *Order) Validate() error {
	if o == nil {
		return fmt.Errorf("order is nil")
	}

	if err := o.validateBasicFields(); err != nil {
		return err
	}
	if err := o.Delivery.Validate(); err != nil {
		return fmt.Errorf("delivery validation failed: %w", err)
	}
	if err := o.Payment.Validate(); err != nil {
		return fmt.Errorf("payment validation failed: %w", err)
	}
	if err := o.Items.Validate(); err != nil {
		return fmt.Errorf("items validation failed: %w", err)
	}

	return nil
}

func (o *Order) validateBasicFields() error {
	if o.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}
	if o.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}
	if o.Entry == "" {
		return fmt.Errorf("entry is required")
	}
	return nil
}

func (d *Delivery) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !phoneRegex.MatchString(d.Phone) {
		return fmt.Errorf("invalid phone format")
	}
	if d.Zip == "" {
		return fmt.Errorf("zip is required")
	}
	if d.City == "" {
		return fmt.Errorf("city is required")
	}
	if d.Address == "" {
		return fmt.Errorf("address is required")
	}
	if d.Email != "" && !emailRegex.MatchString(d.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (p *Payment) Validate() error {
	if p.Transaction == "" {
		return fmt.Errorf("transaction is required")
	}
	if p.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if p.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if p.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	return nil
}

func (items Items) Validate() error {
	if len(items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	for i, item := range items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d validation failed: %w", i+1, err)
		}
	}
	return nil
}

func (i *Item) Validate() error {
	if i.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}
	if i.Name == "" {
		return fmt.Errorf("name is required")
	}
	if i.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	if i.TotalPrice <= 0 {
		return fmt.Errorf("total_price must be positive")
	}
	return nil
}
