package cache

import (
	"sync"
	"testex/internal/models"
)

type OrderCache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]models.Order),
	}
}

func (c *OrderCache) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *OrderCache) Get(orderUID string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	order, ok := c.orders[orderUID]
	return order, ok
}

func (c *OrderCache) GetAll() map[string]models.Order {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make(map[string]models.Order)
	for k, v := range c.orders {
		result[k] = v
	}
	return result
}
