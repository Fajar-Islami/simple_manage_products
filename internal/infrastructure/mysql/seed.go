package mysql

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
)

var (
	orderItemsSeed = make([]daos.OrderItems, 0)
	maxPrice       = 100000
	minPrice       = 10000
)

func maxRandomPrice() int {
	return rand.Intn(maxPrice+minPrice) + minPrice
}

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		expireTime := time.Now().AddDate(0, 0, i)
		orderItemsSeed = append(orderItemsSeed, daos.OrderItems{
			Name:      fmt.Sprintf("Example Positive order items %d", i),
			Price:     maxRandomPrice(),
			ExpiredAt: &expireTime,
		})
	}

	for i := 5; i < 10; i++ {
		expireTime := time.Now().AddDate(0, 0, -i)
		orderItemsSeed = append(orderItemsSeed, daos.OrderItems{
			Name:      fmt.Sprintf("Example Negative order items %d", i),
			Price:     maxRandomPrice(),
			ExpiredAt: &expireTime,
		})
	}
}
