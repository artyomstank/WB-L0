package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOrder(t *testing.T) {
	order := GenerateOrder()

	// Проверяем обязательные поля
	assert.NotEmpty(t, order.OrderUID)
	assert.NotEmpty(t, order.TrackNumber)
	assert.Equal(t, "WBIL", order.Entry)
	assert.NotEmpty(t, order.Delivery.Name)
	assert.NotEmpty(t, order.Delivery.Phone)
	assert.NotEmpty(t, order.Payment.Transaction)
	assert.Greater(t, order.Payment.Amount, 0)
	assert.NotEmpty(t, order.Items)

	// Проверяем корректность сумм
	var itemsTotal int
	for _, item := range order.Items {
		itemsTotal += item.TotalPrice
		assert.GreaterOrEqual(t, item.Price, item.TotalPrice) // Цена со скидкой не больше изначальной
	}

	assert.Equal(t, itemsTotal, order.Payment.GoodsTotal)
	assert.Equal(t, order.Payment.Amount, order.Payment.GoodsTotal+order.Payment.DeliveryCost+order.Payment.CustomFee)
}

func TestGenerateItems(t *testing.T) {
	items := generateItems(3)
	assert.Len(t, items, 3)

	for _, item := range items {
		assert.NotEmpty(t, item.Name)
		assert.NotEmpty(t, item.Brand)
		assert.NotEmpty(t, item.Size)
		assert.Greater(t, item.Price, 0)
		assert.GreaterOrEqual(t, item.Price, item.TotalPrice)
		assert.Equal(t, 200, item.Status)
	}
}

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid phone",
			input:    "1234567890",
			expected: "+71234567890",
		},
		{
			name:     "phone with special chars",
			input:    "+7(999)123-45-67",
			expected: "+79991234567",
		},
		{
			name:     "short phone",
			input:    "12345",
			expected: "+71234500000",
		},
		{
			name:     "non-digits",
			input:    "abc12def34",
			expected: "+71234000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPhoneNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
