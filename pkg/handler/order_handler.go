package handler

import (
	"applicationDesignTest/pkg/domain"
	"applicationDesignTest/pkg/logging"
	"applicationDesignTest/pkg/service"
	"encoding/json"
	"errors"
	"net/http"
)

type OrderHandler struct {
	OrderService *service.OrderService
	// Пример как расширять обработчик
	// TODO: EmailService email.EmailService
	// TODO: DiscountService discount.DiscountService
	// TODO: Logger ext
}

func NewOrderHandler(orderService *service.OrderService /*, emailService email.EmailService*/) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
		// EmailService: emailService,
		// LoggerExt: loggerExt,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder domain.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		//h.Logger.LogWithFields("ERROR", "Invalid request payload", map[string]interface{}{
		//			"error": err,
		//		})
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logging.LogErrorf("Invalid request payload: %v", err)
		return
	}

	err = h.OrderService.CreateOrder(newOrder)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateRange) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logging.LogErrorf("Invalid date range: %v", err)
			return
		} //  Нежно обернем кривую дату
		if errors.Is(err, service.ErrRoomUnavailable) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logging.LogErrorf("Room unavailable: %v", err)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logging.LogErrorf("Error creating order: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)

	logging.LogInfo("Order successfully created: %v", newOrder)
	//h.Logger.LogWithFields("INFO", "Order successfully created", map[string]interface{}{
	//		"orderID": newOrder.ID,
	//		"email":   newOrder.Email,
	//	})
	// расширенный логнер сделан ради этого

	/* TODO:
	err = h.EmailService.SendConfirmation(newOrder.Email, newOrder.ID)
	if err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		logging.LogErrorf("Failed to send confirmation email: %v", err)
		return
	}*/
}
