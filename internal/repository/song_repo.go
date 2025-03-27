package repository

import (
	"github.com/jaam8/online_song_library/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SongRepository struct {
	db *gorm.DB
	l  *zap.Logger
}

func New(db *gorm.DB, log *zap.Logger) *SongRepository {
	return &SongRepository{db: db, l: log}
}

func (s *SongRepository) CreateSong(song *models.Song) (uint, error) {
	s.l.Debug("starting create song", zap.Any("song", song))
	err := s.db.Create(&song).Error
	if err != nil {
		s.l.Error("create song failed", zap.Error(err))
		return 0, err
	}
	s.l.Debug("song created", zap.Uint("id", song.ID))
	return song.ID, nil
}

func (s *SongRepository) GetAllSongs(limit, offset int, filters map[string]interface{}) ([]models.Song, int64, error) {
	s.l.Debug("starting get all songs",
		zap.Any("filters", filters),
		zap.Int("limit", limit),
		zap.Int("offset", offset))
	var songs []models.Song
	var totalCount int64

	baseQuery := s.db.Model(&models.Song{})
	for key, value := range filters {
		if key == "group" {
			baseQuery = baseQuery.Where(`"group" = ?`, value)
		} else {
			baseQuery = baseQuery.Where(key+" = ?", value)
		}
	}

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		s.l.Error("failed to count songs", zap.Error(err))
		return nil, 0, err
	}

	query := baseQuery.Limit(limit).Offset((offset - 1) * limit)
	if err := query.Find(&songs).Error; err != nil {
		s.l.Error("failed to get songs", zap.Error(err))
		return nil, 0, err
	}
	s.l.Debug("retrieved songs",
		zap.Int("count", len(songs)),
		zap.Int64("total", totalCount))
	return songs, totalCount, nil
}

func (s *SongRepository) GetSong(id uint) (*models.Song, error) {
	s.l.Debug("starting get song", zap.Uint("id", id))
	var song models.Song
	result := s.db.First(&song, id)
	if result.Error != nil {
		s.l.Error("failed to get song",
			zap.Uint("id", id),
			zap.Error(result.Error))
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		s.l.Warn("no song found", zap.Uint("id", id))
		return nil, gorm.ErrRecordNotFound
	}
	s.l.Debug("song retrieved", zap.Uint("id", song.ID))
	return &song, nil
}

func (s *SongRepository) UpdateSong(id uint, updatedSong models.Song) error {
	s.l.Debug("starting update song",
		zap.Uint("id", id),
		zap.Any("update data", updatedSong))
	result := s.db.Where("id = ?", id).Updates(updatedSong)
	if result.Error != nil {
		s.l.Error("failed to update song",
			zap.Uint("id", id),
			zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		s.l.Warn("no song updated", zap.Uint("id", id))
		return gorm.ErrRecordNotFound
	}
	s.l.Debug("song updated successfully",
		zap.Uint("id", id),
		zap.Int64("rowsAffected", result.RowsAffected))
	return nil
}

func (s *SongRepository) DeleteSong(id uint) error {
	s.l.Debug("starting delete song", zap.Uint("id", id))
	result := s.db.Where("id = ?", id).Delete(&models.Song{})
	if result.Error != nil {
		s.l.Error("failed to delete song",
			zap.Uint("id", id),
			zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		s.l.Warn("no song deleted", zap.Uint("id", id))
		return gorm.ErrRecordNotFound
	}
	s.l.Debug("song deleted successfully", zap.Uint("id", id))
	return nil
}
