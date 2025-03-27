package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaam8/online_song_library/internal/models"
	"github.com/jaam8/online_song_library/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"time"
)

var ErrParsingTime = errors.New("error parsing time")

type SongService struct {
	repo       *repository.SongRepository
	l          *zap.Logger
	swaggerUrl string
}

func New(repo *repository.SongRepository, log *zap.Logger, swaggerUrl string) *SongService {
	return &SongService{repo: repo, l: log, swaggerUrl: swaggerUrl}
}

func (s *SongService) CreateSong(group, songName string) (uint, error) {
	s.l.Debug("starting create song",
		zap.String("group", group),
		zap.String("song", songName))

	var songRaw models.SongRaw
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", songName)
	reqURL := fmt.Sprintf("%s?%s", s.swaggerUrl, params.Encode())
	s.l.Debug("sending request to swagger", zap.String("url", reqURL))

	resp, err := http.Get(reqURL)
	if err != nil {
		s.l.Error("http get failed", zap.Error(err))
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.l.Error("failed to read swagger response", zap.Error(err))
		return 0, err
	}
	s.l.Debug("received response from swagger",
		zap.Int("status", resp.StatusCode),
		zap.String("body", string(body)))

	if resp.StatusCode != http.StatusOK {
		s.l.Error("swagger response not ok",
			zap.Int("status", resp.StatusCode))
		return 0, fmt.Errorf("swagger error: status %d", resp.StatusCode)
	}

	if err = json.Unmarshal(body, &songRaw); err != nil {
		s.l.Error("failed to unmarshal swagger response", zap.Error(err))
		return 0, err
	}
	s.l.Debug("unmarshaled swagger response",
		zap.String("releaseDate", songRaw.ReleaseDate))

	releaseDate, err := time.Parse("02.01.2006", songRaw.ReleaseDate)
	if err != nil {
		s.l.Error("failed to parse release_date",
			zap.String("releaseDate", songRaw.ReleaseDate),
			zap.Error(err))
		return 0, err
	}

	song := models.Song{
		Group:       group,
		Song:        songName,
		ReleaseDate: releaseDate,
		Text:        songRaw.Text,
		Link:        songRaw.Link,
	}
	s.l.Debug("creating song entity",
		zap.String("group", song.Group),
		zap.String("song", song.Song),
		zap.Time("releaseDate", song.ReleaseDate))

	id, err := s.repo.CreateSong(&song)
	if err != nil {
		s.l.Error("create song failed", zap.Error(err))
		return 0, err
	}
	s.l.Info("song created successfully", zap.Uint("id", id))
	return id, nil
}

func (s *SongService) GetSong(id uint) (*models.Song, error) {
	s.l.Debug("retrieving song", zap.Uint("id", id))
	song, err := s.repo.GetSong(id)
	if err != nil {
		s.l.Error("failed to retrieve song",
			zap.Uint("id", id),
			zap.Error(err))
	} else {
		s.l.Debug("song retrieved", zap.Uint("id", song.ID))
	}
	return song, err
}

func (s *SongService) GetAllSong(limit, offset int, filters map[string]interface{}) ([]models.Song, int64, error) {
	s.l.Debug("retrieving all songs",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Any("filters", filters))
	if val, ok := filters["release_date"]; ok {
		releaseDate, err := time.Parse("02.01.2006", val.(string))
		if err != nil {
			s.l.Error("failed to parse release_date",
				zap.String("release_date", val.(string)),
				zap.Error(err))
			return nil, 0, ErrParsingTime
		}
		filters["release_date"] = releaseDate
	}
	songs, totalCount, err := s.repo.GetAllSongs(limit, offset, filters)
	if err != nil {
		s.l.Error("failed to retrieve songs", zap.Error(err))
	} else {
		s.l.Debug("retrieved songs",
			zap.Int("count", len(songs)),
			zap.Int64("total", totalCount))
	}
	return songs, totalCount, err
}

func (s *SongService) UpdateSong(id uint, updatedSong models.SongRaw) error {
	s.l.Debug("starting update song", zap.Uint("id", id))
	releaseDate, err := time.Parse("02.01.2006", updatedSong.ReleaseDate)
	if err != nil {
		s.l.Error("failed to parse release_date",
			zap.String("release_date", updatedSong.ReleaseDate),
			zap.Error(err))
		return err
	}
	song := models.Song{
		Group:       updatedSong.Group,
		Song:        updatedSong.Song,
		ReleaseDate: releaseDate,
		Text:        updatedSong.Text,
		Link:        updatedSong.Link,
	}
	s.l.Debug("updating song entity",
		zap.String("group", song.Group),
		zap.String("song", song.Song),
		zap.Time("releaseDate", song.ReleaseDate))
	err = s.repo.UpdateSong(id, song)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.l.Warn("song not found", zap.Uint("id", id))
		} else {
			s.l.Error("update song failed",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return err
	}
	s.l.Info("song updated successfully", zap.Uint("id", id))
	return nil
}

func (s *SongService) DeleteSong(id uint) error {
	s.l.Debug("starting delete song", zap.Uint("id", id))
	err := s.repo.DeleteSong(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.l.Warn("song not found", zap.Uint("id", id))
		} else {
			s.l.Error("delete song failed",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return err
	}
	s.l.Info("song deleted successfully", zap.Uint("id", id))
	return nil
}
