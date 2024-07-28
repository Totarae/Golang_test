// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"applicationDesignTest/pkg/handler"
	"applicationDesignTest/pkg/logging"
	"applicationDesignTest/pkg/service"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//loggerext := logging.NewLogger()
	orderService := service.NewOrderService()
	// emailService := email.NewSMTPEmailService()
	// emailService := email.NewMockEmailService() // mock для тестов
	// discountService := discount.newDiscountService()
	orderHandler := handler.NewOrderHandler(orderService /*, emailService, discountService, loggerext*/)
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", orderHandler.CreateOrder)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel чтобы перехватывать системные вызовы для остановки сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		logging.LogInfo("Server listening on localhost:8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.LogErrorf("Server failed: %s", err)
			os.Exit(1)
		}
	}()

	// Блокируем основной поток, пока не получим прерывание
	<-quit
	logging.LogInfo("Server is shutting down...")

	// Таймаут для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершаем в стандартнолм режиме
	if err := server.Shutdown(ctx); err != nil {
		logging.LogErrorf("Server forced to shutdown: %s", err)
	}

	logging.LogInfo("Server exiting")
}
