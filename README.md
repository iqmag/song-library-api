# Song Library API

## Введение
Song Library API предназначен для разработчиков, которые хотят интегрировать функциональность управления музыкальными библиотеками в свои приложения. API позволяет легко добавлять, обновлять, удалять и получать информацию о песнях.

## Описание
Song Library API - это RESTful API для управления песнями в музыкальной библиотеке. Он позволяет пользователям выполнять операции CRUD (создание, чтение, обновление, удаление) над песнями, а также получать текст песен.

## Технологии
- Go
- Gin
- GORM
- PostgreSQL
- Swagger

## Установка
1. Клонируйте репозиторий:
    ```bash
    git clone https://github.com/iqmag/song-library-api.git
    ```

2. Перейдите в директорию проекта:
    ```bash
    cd song-library-api
    ```

3. Установите зависимости:
    ```bash
    go mod tidy
    ```

4. Настройте базу данных PostgreSQL и обновите конфигурацию подключения в файле `.env`

5. Запустите приложение:
    ```bash
    go run main.go
    ```
    
## Эндпоинты
- Получить все песни: GET /api/v1/songs
- Создать новую песню: POST /api/v1/songs
- Обновить песню по ID: PUT /api/v1/songs/{id}
- Удалить песню по ID: DELETE /api/v1/songs/{id}
- Получить текст песни по ID: GET /api/v1/songs/{id}/text

**Пример запроса в Postman:**

1. Откройте Postman.
2. Создайте новый запрос.
3. Выберите метод `GET`.
4. Введите URL: http://localhost:8080/api/v1/songs
5. Нажмите кнопку **"Send"**
6. Пример успешного ответа (200 OK):

Пример вывода:

    ```json
[
 {
     "id": 1,
     "song_group": "Muse",
     "song_name": "Supermassive Black Hole",
     "release_date": "2006-06-19",
     "text": "Текст песни...",
     "links": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
 },
 {
     "id": 2,
     "song_group": "Nirvana",
     "song_name": "Smells Like Teen Spirit",
     "release_date": "1991-09-10",
     "text": "Текст песни...",
     "links": "https://www.youtube.com/watch?v=hTWKbfoikeg"
 },
 {
     "id": 3,
     "song_group": "The Beatles",
     "song_name": "Hey Jude",
     "release_date": "1968-08-26",
     "text": "Текст песни...",
     "links": "https://www.youtube.com/watch?v=A_MjCqQoLLA"
 }
]
    ```

## Ошибки
API может возвращать следующие коды ошибок:
- `400 Bad Request` - Неверный запрос (например, отсутствуют обязательные поля)
- `404 Not Found` - Песня с указанным ID не найдена
- `500 Internal Server Error` - Ошибка на сервере

## Документация API
Документация API доступна по адресу [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html). Вы можете использовать Swagger UI для тестирования эндпоинтов и просмотра доступных операций