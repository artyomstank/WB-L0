package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Validate(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		wantErr bool
	}{
		{
			name:    "nil order",
			order:   nil,
			wantErr: true,
		},
		{
			name: "valid order",
			order: &Order{
				OrderUID:    "test-123",
				TrackNumber: "WBIL123456789",
				Entry:       "WBIL",
				Delivery: Delivery{
					Name:    "Test User",
					Phone:   "+79991234567",
					Zip:     "123456",
					City:    "Moscow",
					Address: "Test St, 1",
				},
				Payment: Payment{
					Transaction: "tx-123",
					Currency:    "USD",
					Provider:    "stripe",
					Amount:      100,
				},
				Items: []Item{
					{
						ChrtID:      1,
						TrackNumber: "WBIL123456789",
						Price:       50,
						Name:        "Test Item",
						TotalPrice:  45,
						Status:      200,
					},
				},
				DateCreated: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			order: &Order{
				OrderUID: "test-123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestDelivery_Validate(t *testing.T) {
	tests := []struct {
		name     string
		delivery Delivery
		wantErr  bool
	}{
		{
			name: "valid delivery",
			delivery: Delivery{
				Name:    "Test User",
				Phone:   "+79991234567",
				Zip:     "123456",
				City:    "Moscow",
				Address: "Test St, 1",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			delivery: Delivery{
				Phone: "+79991234567",
			},
			wantErr: true,
		},
		{
			name: "invalid phone",
			delivery: Delivery{
				Name:  "Test User",
				Phone: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.delivery.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
