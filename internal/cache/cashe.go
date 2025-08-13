package cache

import (
	"database/sql"
	"hash/fnv"
	"log"
	"sync"
	"testex/internal/db"
	"testex/internal/models"
	"time"
)

type CasheOrder struct {
	Order   models.Order
	AddedAt time.Time
}

type shard struct {
	mu     sync.RWMutex
	orders map[string]CasheOrder
}

type OrderCache struct {
	shards []shard
	ttl    time.Duration
	count  uint32
}

func NewOrderCache(shardCount int, ttl time.Duration) *OrderCache {
	shards := make([]shard, shardCount)
	for i := range shards {
		shards[i] = shard{
			orders: make(map[string]CasheOrder),
		}
	}
	return &OrderCache{
		shards: shards,
		ttl:    ttl,
		count:  uint32(shardCount),
	}
}

func (c *OrderCache) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return &c.shards[h.Sum32()%c.count]
}

func (c *OrderCache) Set(order models.Order) error {
	s := c.getShard(order.OrderUID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUID] = CasheOrder{
		Order:   order,
		AddedAt: time.Now(),
	}
	return nil
}

func (c *OrderCache) Get(orderUID string) (models.Order, bool) {
	s := c.getShard(orderUID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	cached, ok := s.orders[orderUID]
	if !ok {
		return models.Order{}, false
	}
	if time.Since(cached.AddedAt) > c.ttl {
		return models.Order{}, false
	}
	return cached.Order, true
}

func (c *OrderCache) GetAll() map[string]models.Order {
	result := make(map[string]models.Order)
	for i := range c.shards {
		s := &c.shards[i]
		s.mu.RLock()
		for k, v := range s.orders {
			if time.Since(v.AddedAt) <= c.ttl {
				result[k] = v.Order
			}
		}
		s.mu.RUnlock()
	}
	return result
}

func WarmUpCache(dbConn *sql.DB, orderCache *OrderCache) error {
	orders, err := db.LoadAllOrdersFromDB(dbConn)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := orderCache.Set(order); err != nil {
			log.Fatal("Ошибка логирования кеша", err)
		}
	}

	return nil
}
