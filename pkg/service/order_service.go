package service

import (
	"applicationDesignTest/pkg/domain"
	"errors"
	"math"
	"time"
)

type OrderService struct {
	orders       []domain.Order
	availability []domain.RoomAvailability
}

var ErrInvalidDateRange = errors.New("invalid date range")
var ErrRoomUnavailable = errors.New("hotel room is not available for selected dates")

func NewOrderService() *OrderService {
	return &OrderService{
		orders: []domain.Order{},
		availability: []domain.RoomAvailability{
			{"reddison", "lux", date(2024, 1, 1), 1},
			{"reddison", "lux", date(2024, 1, 2), 1},
			{"reddison", "lux", date(2024, 1, 3), 1},
			{"reddison", "lux", date(2024, 1, 4), 1},
			{"reddison", "lux", date(2024, 1, 5), 0},
			{"hilton", "eco", date(2024, 2, 2), 1},
			{"hilton", "eco", date(2024, 2, 3), 1},
			{"otherHotel", "stan", date(2024, 1, 1), 1},
			{"otherHotel", "stan", date(2024, 1, 2), 1},
		},
	}
}

func (s *OrderService) CreateOrder(newOrder domain.Order) error {
	// для бронирования нескольких номеров
	//for _, room := range newOrder.Rooms {
	//daysToBook := daysBetween(room.From, room.To)
	daysToBook := daysBetween(newOrder.From, newOrder.To)
	if len(daysToBook) == 0 {
		return ErrInvalidDateRange // Пробросим дальше во врапере в обработчик
	}

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	for _, dayToBook := range daysToBook {
		for i, availability := range s.availability {
			// Пройдемся по коллекции
			if availability.HotelID != newOrder.HotelID || availability.RoomID != availability.RoomID {
				continue
			}
			if !availability.Date.Equal(dayToBook) || availability.Quota < 1 {
				continue
			}
			availability.Quota -= 1
			s.availability[i] = availability
			delete(unavailableDays, dayToBook)
		}
	}

	if len(unavailableDays) != 0 {
		//http.Error(w, "Hotel room is not available for selected dates", http.StatusInternalServerError)
		//logging.LogErrorf("Hotel room is not available for selected dates:\n%v\n%v", newOrder, unavailableDays)
		//return
		return ErrRoomUnavailable
		//return errors.New("hotel room is not available for selected dates")
	}

	s.orders = append(s.orders, newOrder)

	return nil
}

func daysBetween(from time.Time, to time.Time) []time.Time {
	// Приведем к UTC  и не хотим делать каст в цикле
	// TODO: увеличить точность до час/минутв
	from = toDay(from)
	to = toDay(to)

	if from.After(to) {
		return []time.Time{}
	} // from после to вернем пустой

	// Рассчитаем разницу 'from' и 'to' с использованием Duration
	numDays := int(math.Ceil(to.Sub(from).Hours()/24)) + 1
	days := make([]time.Time, numDays) // не будем создавать слайс с 0 длиной

	// TODO: Перепилить завтра, баг с пустыми значениями в слайсе
	for i, d := 0, from; !d.After(to); d = d.AddDate(0, 0, 1) {
		days[i] = d
		i++
	}

	return days
}

// приведение к UTC
func toDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}

// Прсото конвыертер
func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
