package cache

import (
	"sync"
	"testex/internal/models"
	"time"
)

type OrderCache struct {
	mu     sync.RWMutex
	orders map[string]CasheOrder
	ttl    time.Duration
}
type CasheOrder struct {
	order models.Order
	added time.Time
}

func NewOrderCache(ttl time.Duration) (*OrderCache, error) {
	return &OrderCache{
		orders: make(map[string]CasheOrder),
		ttl:    ttl,
	}, nil
}

func (c *OrderCache) Set(order models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = CasheOrder{
		order: order,
		added: time.Now(),
	}
	return nil
}

func (c *OrderCache) Get(orderUID string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	order, ok := c.orders[orderUID]
	if !ok {
		return models.Order{}, false
	}
	if time.Since(order.added) > c.ttl {
		return models.Order{}, false
	}
	return order.order, true
}

func (c *OrderCache) GetAll() (map[string]models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make(map[string]models.Order)
	for k, v := range c.orders {
		result[k] = v.order
	}
	return result, nil
}
