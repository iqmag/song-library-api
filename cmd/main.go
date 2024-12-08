package main

import (
	"log"
	"song-library/internal/config"
	"song-library/internal/handlers"
	"song-library/internal/models"
	"song-library/internal/repository"
	"song-library/internal/service"

	_ "song-library/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Music Library API
// @version 1.0
// @description API для управления песнями в музыкальной библиотеке
// @BasePath /api/v1

func main() {
    // Загрузка конфигурации
    config.LoadConfig()

    // Настройка БД
    db, err := gorm.Open(postgres.Open(config.GetDatabaseURL()), &gorm.Config{})
        handleError(err, "не удалось подключиться к базе данных")

    // Миграции
    log.Println("Начинаем миграцию базы данных...")
	models.Migrate(db)
	log.Println("Миграция завершена успешно!")

    // Добавление начальных данных
	models.AddInitialSongs(db)

    // Обработчики и сервисы
    songRepo := repository.NewSongRepository(db)
    songService := service.NewSongService()
    songHandler := handlers.NewSongHandler(songRepo, songService)

    // Инициализация маршрутов
    r := gin.Default()
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    setupRoutes(r, songHandler)

    // Запуск сервера
    port := config.GetPort()
    if port == "" {
        log.Fatal("SERVER_PORT не установлен")
    }

	log.Printf("Сервер запущен на порту %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("не удалось запустить сервер", err)
    }
}

// handleError обрабатывает ошибки и выводит сообщение в лог
func handleError(err error, message string) {
    if err != nil {
        log.Fatal(message, err) // Завершение программы при ошибке
    }
}

// setupRoutes настраивает маршруты для API
func setupRoutes(r *gin.Engine, songHandler *handlers.SongHandler) {
    api := r.Group("/api/v1")
    {
        // Получить все песни
        // @Summary Получить все песни
        // @Description Получить список всех песен
        // @Tags песни
        // @Produce json
        // @Success 200 {array} models.Song "Список песен"
        // @Router /songs [get]
        api.GET("/songs", songHandler.GetSongs)

        // Создать новую песню
        // @Summary Создать новую песню
        // @Description Создать новую песню в библиотеке
        // @Tags песни
        // @Accept json
        // @Produce json
        // @Param song body models.Song true "Данные песни"
        // @Success 201 {object} models.Song "Песня успешно создана"
        // @Failure 400 {string} string "Неверный ввод"
        // @Router /songs [post]
        api.POST("/songs", songHandler.CreateSong)

        // Обновить песню по ID
        // @Summary Обновить песню по ID
        // @Description Обновить существующую песню в библиотеке
        // @Tags песни
        // @Accept json
        // @Produce json
        // @Param id path int true "ID песни"
        // @Param song body models.Song true "Обновленные данные песни"
        // @Success 200 {object} models.Song "Песня успешно обновлена"
        // @Failure 400 {string} string "Неверный ввод"
        // @Failure 404 {string} string "Песня не найдена"
        // @Router /songs/{id} [put]
        api.PUT("/songs/:id", songHandler.UpdateSong)

        // Удалить песню по ID
        // @Summary Удалить песню по ID
        // @Description Удалить существующую песню из библиотеки
        // @Tags песни
        // @Produce json
        // @Param id path int true "ID песни"
        // @Success 204 {string} string "Песня успешно удалена"
        // @Failure 404 {string} string "Песня не найдена"
        // @Router /songs/{id} [delete]
        api.DELETE("/songs/:id", songHandler.DeleteSong)

        // Получить текст песни по ID
        // @Summary Получить текст песни по ID
        // @Description Получить текст песни по ее ID
        // @Tags песни
        // @Produce json
        // @Param id path int true "ID песни"
        // @Success 200 {string} string "Текст песни"
        // @Failure 404 {string} string "Песня не найдена"
        // @Router /songs/{id}/text [get]
        api.GET("/songs/:id/text", songHandler.GetSongText)
    }
}