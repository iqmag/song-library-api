package service

import (
    "encoding/json"
    "fmt"
    "net/http"
    "song-library/internal/models"
)

type SongService interface {
    GetSongInfo(group, songName string) (models.Song, error)
}

type songService struct{}

func NewSongService() SongService {
    return &songService{}
}

func (s *songService) GetSongInfo(group, songName string) (models.Song, error) {
    // URL для запроса к API
    url := fmt.Sprintf("http://localhost:8080/api/v1/songs?group=%s&song=%s", group, songName)
    fmt.Println("Requesting song info from:", url)
    
    resp, err := http.Get(url)
    if err != nil {
        return models.Song{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return models.Song{}, fmt.Errorf("Failed to get song info: %s", resp.Status)
    }

    var songInfo models.Song
    if err := json.NewDecoder(resp.Body).Decode(&songInfo); err != nil {
        return models.Song{}, err
    }

    return songInfo, nil
}