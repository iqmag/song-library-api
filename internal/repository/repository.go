package repository

import (
    "gorm.io/gorm"
    "song-library/internal/models"
)

// SongRepository - интерфейс для работы с песнями
type SongRepository interface {
    GetAllFiltered(song_group, song_name, release_date string, page, pageSize int) ([]models.Song, error)
    Create(song models.Song) error
    Update(song models.Song) error
    Delete(id uint) error
    GetByID(id uint) (models.Song, error) // Новый метод
}

// songRepository - структура, реализующая интерфейс SongRepository
type songRepository struct {
    db *gorm.DB
}

// NewSongRepository - конструктор для создания нового репозитория
func NewSongRepository(db *gorm.DB) SongRepository {
    return &songRepository{db: db}
}

// GetAllFiltered - получение всех песен с фильтрацией и пагинацией
func (r *songRepository) GetAllFiltered(song_group, song_name, release_date string, page, pageSize int) ([]models.Song, error) {
    var songs []models.Song
    offset := (page - 1) * pageSize
    err := r.db.Where("song_group LIKE ? AND song_name LIKE ? AND release_date LIKE ?", "%"+song_group+"%", "%"+song_name+"%", "%"+release_date+"%").
        Limit(pageSize).Offset(offset).Find(&songs).Error
    return songs, err
}

func (r *songRepository) Create(song models.Song) error {
    return r.db.Create(&song).Error
}

func (r *songRepository) Update(song models.Song) error {
    return r.db.Save(&song).Error
}

func (r *songRepository) Delete(id uint) error {
    return r.db.Delete(&models.Song{}, id).Error
}

func (r *songRepository) GetByID(id uint) (models.Song, error) {
    var song models.Song
    err := r.db.First(&song, id).Error
    return song, err
}