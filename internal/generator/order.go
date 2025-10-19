package generator

import (
	"L0-wb/internal/models"
	"fmt"
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateOrder() *models.Order {
	fake := gofakeit.New(0)
	now := time.Now()

	// Генерируем items с корректными ценами
	items := generateItems(fake.IntRange(1, 5))

	// Считаем общую стоимость товаров без округления, т.к. TotalPrice уже int
	var goodsTotal int
	for _, item := range items {
		goodsTotal += item.TotalPrice
	}

	deliveryCost := fake.IntRange(300, 2000)
	customFee := int(math.Round(float64(goodsTotal) * 0.05)) // 5% от стоимости товаров

	// Общая сумма заказа (все значения уже в int)
	totalAmount := goodsTotal + deliveryCost + customFee

	// Генерируем трек-номер в формате WBILXXXXXXXX
	trackNumber := fmt.Sprintf("WBIL%d", fake.IntRange(10000000, 99999999))

	return &models.Order{
		OrderUID:    fake.UUID(),
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    fake.Name(),
			Phone:   formatPhoneNumber(fake.Phone()),
			Zip:     fmt.Sprintf("%05d", fake.IntRange(100000, 999999)),
			City:    fake.City(),
			Address: formatAddress(fake.Street(), fake.StreetNumber()),
			Region:  fake.State(),
			Email:   fake.Email(),
		},
		Payment: models.Payment{
			Transaction:  fake.UUID(),
			RequestID:    "", // Может быть пустым
			Currency:     "USD",
			Provider:     fake.RandomString([]string{"stripe", "paypal", "bank_transfer"}),
			Amount:       totalAmount,
			PaymentDt:    int(now.Unix()),
			Bank:         fake.RandomString([]string{"Sber", "Tinkoff", "Alpha", "VTB"}),
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    customFee,
		},
		Items:             items,
		Locale:            fake.RandomString([]string{"en", "ru"}),
		InternalSignature: "", // Может быть пустым
		CustomerID:        fmt.Sprintf("customer_%s", fake.Username()),
		DeliveryService:   fake.RandomString([]string{"DHL", "FedEx", "UPS", "SDEK", "Russian Post"}),
		Shardkey:          fmt.Sprintf("%d", fake.IntRange(1, 10)),
		SmID:              fake.IntRange(1, 999),
		DateCreated:       now,
		OofShard:          fmt.Sprintf("%d", fake.IntRange(1, 10)),
	}
}

func generateItems(count int) []models.Item {
	fake := gofakeit.New(0)
	items := make([]models.Item, count)

	for i := 0; i < count; i++ {
		basePrice := fake.IntRange(100, 10000)
		sale := fake.IntRange(0, 50) // Скидка до 50%

		// Расчет цены со скидкой
		totalPrice := float64(basePrice) * (100 - float64(sale)) / 100
		// Округляем до целого числа для TotalPrice
		roundedTotalPrice := int(math.Round(totalPrice))

		items[i] = models.Item{
			ChrtID:      fake.IntRange(1000000, 9999999),
			TrackNumber: fmt.Sprintf("WBIL%d", fake.IntRange(10000000, 99999999)),
			Price:       basePrice,
			Rid:         fake.UUID(),
			Name:        generateProductName(fake),
			Sale:        sale,
			Size:        generateSize(fake),
			TotalPrice:  roundedTotalPrice,
			NmID:        fake.IntRange(1000000, 9999999),
			Brand:       generateBrand(fake),
			Status:      200, // 200 - доставлен
		}
	}

	return items
}

func generateProductName(fake *gofakeit.Faker) string {
	categories := []string{
		"Футболка", "Джинсы", "Куртка", "Платье", "Кроссовки",
		"Рубашка", "Юбка", "Свитер", "Брюки", "Пальто",
	}
	colors := []string{
		"черный", "белый", "синий", "красный", "зеленый",
		"серый", "бежевый", "розовый", "желтый", "коричневый",
	}
	return fmt.Sprintf("%s %s", fake.RandomString(colors), fake.RandomString(categories))
}

func generateBrand(fake *gofakeit.Faker) string {
	brands := []string{
		"NIKE", "ADIDAS", "PUMA", "REEBOK", "NEW BALANCE",
		"ZARA", "H&M", "UNIQLO", "LEVIS", "TOMMY HILFIGER",
	}
	return fake.RandomString(brands)
}

func generateSize(fake *gofakeit.Faker) string {
	sizes := []string{"XS", "S", "M", "L", "XL", "XXL"}
	return fake.RandomString(sizes)
}

func formatPhoneNumber(phone string) string {
	// Извлекаем только цифры из номера
	var digits []rune
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}

	// Берем последние 10 цифр, если их больше 10
	// или дополняем нулями, если меньше
	result := make([]rune, 10)
	if len(digits) >= 10 {
		copy(result, digits[len(digits)-10:])
	} else {
		copy(result, digits)
		for i := len(digits); i < 10; i++ {
			result[i] = '0'
		}
	}

	return "+7" + string(result)
}

func formatAddress(street, number string) string {
	return fmt.Sprintf("%s, д. %s", street, number)
}
