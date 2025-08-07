package hanlders

import (
	"encoding/json"
	"net/http"
	"testex/internal/cache"
	"testex/internal/db"
	"testex/internal/service"
)

func GetOrder(orderCache *cache.OrderCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("order_id")
		if orderID == "" {
			http.Error(w, "Нет параметра order_id", http.StatusBadRequest)
			return
		}

		order, found := orderCache.Get(orderID)
		if !found {
			var err error
			order, err = service.GetOrderByID(db.DB, orderID)
			if err != nil {
				http.Error(w, "заказ не найден", http.StatusNotFound)
				return
			}

			orderCache.Set(order)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "Ошибка кодировки json", http.StatusBadRequest)
		}
	}
}
