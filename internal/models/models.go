package models

import (
	"log"

	"gorm.io/gorm"
)

// Модель песни
type Song struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	Group       string `gorm:"column:song_group;index" json:"song_group"` // Индекс
	Song        string `gorm:"column:song_name;index" json:"song_name"`   // Индекс
	ReleaseDate string `gorm:"column:release_date;index" json:"release_date"` // Индекс
	Text        string `gorm:"column:text" json:"text"`
	Links       string `gorm:"column:links" json:"links"`
}

// Функция миграции
func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Song{}); err != nil {
		panic("не удалось выполнить миграцию базы данных")
	}
}

// AddInitialSongs добавляет начальные записи в таблицу songs
func AddInitialSongs(db *gorm.DB) {
	songs := []Song{
		{
			Group:       "Muse",
			Song:        "Supermassive Black Hole",
			ReleaseDate: "2006-06-19",
			Text:        "Текст песни...",
			Links:       "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		},
		{
			Group:       "Nirvana",
			Song:        "Smells Like Teen Spirit",
			ReleaseDate: "1991-09-10",
			Text:        "Текст песни...",
			Links:       "https://www.youtube.com/watch?v=hTWKbfoikeg",
		},
		{
			Group:       "The Beatles",
			Song:        "Hey Jude",
			ReleaseDate: "1968-08-26",
			Text:        "Текст песни...",
			Links:       "https://www.youtube.com/watch?v=A_MjCqQoLLA",
		},
	}

	for _, song := range songs {
		var existingSong Song
		// Проверяем, существует ли уже песня с такими же параметрами
		if err := db.Where("song_group = ? AND song_name = ? AND release_date = ?", song.Group, song.Song, song.ReleaseDate).First(&existingSong).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Если запись не найдена, добавляем новую
				if err := db.Create(&song).Error; err != nil {
					log.Printf("Ошибка при добавлении песни %s: %v", song.Song, err)
				} else {
					log.Printf("Песня %s добавлена в базу данных", song.Song)
				}
			} else {
				log.Printf("Ошибка при проверке существования песни %s: %v", song.Song, err)
			}
		} else {
			log.Printf("Песня %s уже существует в базе данных", song.Song)
		}
	}
}