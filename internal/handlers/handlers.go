package handlers

import (
	"log"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/repository"
	"song-library/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SongHandler - структура для обработки запросов к песням
type SongHandler struct {
    repo repository.SongRepository
    svc  service.SongService
}

func NewSongHandler(repo repository.SongRepository, svc service.SongService) *SongHandler {
    return &SongHandler{repo: repo, svc: svc}
}

//Аннотации
// @Summary Получить список песен
// @Description Получить все песни с возможностью фильтрации и пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Название музыкальной группы" example("Muse")
// @Param song query string false "Название песни" example("Supermassive Black Hole")
// @Param releaseDate query string false "Дата релиза" example("2006-06-19")
// @Param page query int false "Номер страницы" default(1)
// @Param pageSize query int false "Количество элементов на странице" default(10)
// @Success 200 {array} models.Song
// @Router /api/v1/songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
    group := c.DefaultQuery("group", "")
    songName := c.DefaultQuery("song", "")
    releaseDate := c.DefaultQuery("releaseDate", "")
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("pageSize", "10")

    // Преобразуем строки в целые числа
    pageInt, err := strconv.Atoi(page)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page number"})
        return
    }

    pageSizeInt, err := strconv.Atoi(pageSize)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
        return
    }

    // Получаем песни с фильтрацией
    songs, err := h.repo.GetAllFiltered(group, songName, releaseDate, pageInt, pageSizeInt)
    if err != nil {
        log.Println("Error fetching songs:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching songs"})
        return
    }

    c.JSON(http.StatusOK, songs)
}

// @Summary Создать новую песню
// @Description Создает новую песню с обогащенными данными из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные песни (группа и название обязательны)"
// @Success 201 {object} models.Song
// @Failure 400 {object} gin.H "Неверный ввод"
// @Router /api/v1/songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
    var song models.Song
    if err := c.ShouldBindJSON(&song); err != nil {
        log.Printf("Enriching song data for group: %s, song: %s", song.Group, song.Song)
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    // Запрос к внешнему API для получения дополнительной информации
    enrichedSong, err := h.svc.GetSongInfo(song.Group, song.Song)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to enrich song data"})
        return
    }

    // Сохраняем обогащенную песню
    if err := h.repo.Create(enrichedSong); err != nil {
        log.Println("Error creating song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating song"})
        return
    }

    c.JSON(http.StatusCreated, enrichedSong)
}

// @Summary Обновить данные песни
// @Description Обновляет существующую песню по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Обновленные данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} gin.H "Неверный ввод"
// @Failure 404 {object} gin.H "Песня не найдена"
// @Router /api/v1/songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid song ID"})
        return
    }

    var song models.Song
    if err := c.ShouldBindJSON(&song); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    song.ID = uint(id)

    if err := h.repo.Update(song); err != nil {
        log.Println("Error updating song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating song"})
        return
    }

    c.JSON(http.StatusOK, song)
}

// @Summary Удалить песню
// @Description Удаляет песню по ID
// @Tags songs
// @Param id path int true "ID песни"
// @Success 204 "Нет содержимого"
// @Failure 404 {object} gin.H "Песня не найдена"
// @Router /api/v1/songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid song ID"})
        return
    }

    if err := h.repo.Delete(uint(id)); err != nil {
        log.Println("Error deleting song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting song"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

// @Summary Получить текст песни
// @Description Получить текст песни по ID с поддержкой пагинации куплетов
// @Tags songs
// @Param id path int true "ID песни"
// @Success 200 {array} string "Массив куплетов песни"
// @Failure 404 {object} gin.H "Песня не найдена"
// @Router /api/v1/songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid song ID"})
        return
    }

    // Получение информации о песне
    song, err := h.repo.GetByID(uint(id))
    if err != nil {
        log.Println("Error fetching song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching song"})
        return
    }

    // Разделение текста на куплеты
    verses := splitTextIntoVerses(song.Text)

    c.JSON(http.StatusOK, verses)
}

// Разделение текста на куплеты
func splitTextIntoVerses(text string) []string {
    return strings.Split(text, "\n\n") // Разделение по двум новым строкам
}