package handler

import (
	"applicationDesignTest/pkg/domain"
	"applicationDesignTest/pkg/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateOrderHandler(t *testing.T) {

	orderService := service.NewOrderService()
	orderHandler := NewOrderHandler(orderService)

	order := domain.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "guest@mail.ru",
		From:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// POST
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orderHandler.CreateOrder)

	// Запускаем
	handler.ServeHTTP(rr, req)

	// Check status
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Ответ
	expected := domain.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "guest@mail.ru",
		From:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	// маршал, чтобы не возиться с проблами
	var returnedOrder domain.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &returnedOrder); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if returnedOrder != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			returnedOrder, expected)
	}
}

func TestCreateUnavailableOrderHandler(t *testing.T) {

	orderService := service.NewOrderService()
	orderHandler := NewOrderHandler(orderService)

	order := domain.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "guest@mail.ru",
		From:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// POST
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orderHandler.CreateOrder)

	// Запускаем
	handler.ServeHTTP(rr, req)

	// Check status
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := service.ErrRoomUnavailable.Error()
	if strings.TrimSuffix(rr.Body.String(), "\n") != expected {

		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateNoQuotaOrderHandler(t *testing.T) {

	orderService := service.NewOrderService()
	orderHandler := NewOrderHandler(orderService)

	order := domain.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "guest@mail.ru",
		From:      time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// POST
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orderHandler.CreateOrder)

	// Запускаем
	handler.ServeHTTP(rr, req)

	// Check status
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := service.ErrRoomUnavailable.Error()
	if strings.TrimSuffix(rr.Body.String(), "\n") != expected {

		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
