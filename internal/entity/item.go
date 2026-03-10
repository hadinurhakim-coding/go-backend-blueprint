package entity

import "time"

// NowFunc dipakai untuk mendapatkan waktu sekarang (bisa di-override saat testing).
var NowFunc = time.Now

// Item adalah entitas contoh untuk CRUD. Nantinya bisa diganti atau ditambah field.
type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
