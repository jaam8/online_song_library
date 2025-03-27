package models

import "time"

type Song struct {
	ID          uint      `json:"id"  example:"1" gorm:"primaryKey"`
	Group       string    `json:"group" example:"Muse"`
	Song        string    `json:"song" example:"Supermassive Black Hole"`
	ReleaseDate time.Time `json:"release_date" example:"2006-06-19T00:00:00Z"`
	Text        string    `json:"text" example:"Ooh baby, don't you know I suffer?\n..."`
	Link        string    `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

// SongRaw нужен, чтобы правильно парсить ReleaseDate из json
type SongRaw struct {
	//ID          uint   `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
