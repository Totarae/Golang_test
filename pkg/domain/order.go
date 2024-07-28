package domain

import "time"

type Order struct {
	//ID        string    `json:"id"`
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	//.Rooms []Room `json:"rooms"`
}

// Для заказа нескольких номеров
/*type Room struct {
	HotelID string    `json:"hotelName"`
	RoomID  string    `json:"roomType"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}*/
