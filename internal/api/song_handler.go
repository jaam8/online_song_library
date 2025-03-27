package api

import (
	"errors"
	"github.com/jaam8/online_song_library/internal/models"
	"github.com/jaam8/online_song_library/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type SongHandler struct {
	service *service.SongService
	l       *zap.Logger
}
type ErrorResponse struct {
	Error string `json:"error" swaggertype:"string" example:"error text"`
}

func New(service *service.SongService, log *zap.Logger) *SongHandler {
	return &SongHandler{service: service, l: log}
}

// @Summary Добавление новой песни
// @Description Добавляет новую песню, получая информацию о ней через запрос к стороннему API, возвращает ID песни
// @Tags songs
// @Accept json
// @Produce json
// @Param song body api.CreateSongHandler.request true "название группы и песни"
// @Success 201 {object} api.CreateSongHandler.successResponse "successfully created" example:{"id": 1}
// @Failure 400 {object} ErrorResponse "invalid request" example:{"error": "invalid request"}
// @Failure 404 {object} ErrorResponse "song not found" example:{"error": "song not found"}
// @Failure 422 {object} ErrorResponse "all field are required" example:{"error": "all field are required"}
// @Failure 500 {object} ErrorResponse "internal server error" example:{"error": "internal server error"}
// @Router / [post]
func (h *SongHandler) CreateSongHandler(c echo.Context) error {
	h.l.Debug("starting create song")
	type request struct {
		Group string `json:"group" example:"Muse" swaggertype:"string"`
		Song  string `json:"song" example:"Supermassive Black Hole" swaggertype:"string"`
	}
	type successResponse struct {
		ID uint `json:"id" example:"1" swaggertype:"integer"`
	}
	var req request
	if err := c.Bind(&req); err != nil {
		h.l.Debug("failed to bind request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{"invalid request"})
	}
	h.l.Debug("parsed request body",
		zap.String("group", req.Group),
		zap.String("song", req.Song))

	if req.Group == "" || req.Song == "" {
		h.l.Debug("validation failed: missing required fields")
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"all field are required"})
	}

	id, err := h.service.CreateSong(req.Group, req.Song)
	if err != nil {
		h.l.Debug("error creating song", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}
	if id == 0 {
		h.l.Warn("song not found")
		return c.JSON(http.StatusNotFound, ErrorResponse{"song not found"})
	}
	h.l.Info("song created successfully", zap.Uint("id", id))
	return c.JSON(http.StatusCreated, successResponse{id})
}

// @Summary Получение всех песен с фильтрацией и пагинацией
// @Description Получение всех песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false " "
// @Param song query string false " "
// @Param release_date query string false " "
// @Param text query string false " "
// @Param link query string false " "
// @Param page query int false " " default(1)
// @Param per_page query int false " " default(5)
// @Success 200 {object} api.GetAllSongsHandler.successResponse "received successfully"
// @Failure 404 {object} ErrorResponse "songs not found" example:{"error": "songs not found"}
// @Failure 422 {object} ErrorResponse "invalid page" example:{"error": "invalid page"}
// @Failure 422 {object} ErrorResponse "invalid per_page" example:{"error": "invalid per_page"}
// @Failure 422 {object} ErrorResponse "invalid release_date" example:{"error": "invalid release_date"}
// @Failure 500 {object} ErrorResponse "internal server error" example:{"error": "internal server error"}
// @Router / [get]
func (h *SongHandler) GetAllSongsHandler(c echo.Context) error {
	filters := make(map[string]interface{})
	if group := c.QueryParam("group"); group != "" {
		filters["group"] = group
	}
	if song := c.QueryParam("song"); song != "" {
		filters["song"] = song
	}
	if releaseDate := c.QueryParam("release_date"); releaseDate != "" {
		filters["release_date"] = releaseDate
	}
	if text := c.QueryParam("text"); text != "" {
		filters["text"] = text
	}
	if link := c.QueryParam("link"); link != "" {
		filters["link"] = link
	}

	page := 1
	perPage := 5

	if pageStr := c.QueryParam("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			h.l.Debug("failed to parse page", zap.Error(err))
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid page"})
		}
		page = p
	}

	if perPageStr := c.QueryParam("per_page"); perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err != nil || pp < 1 {
			h.l.Debug("failed to parse per_page", zap.Error(err))
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid per_page"})
		}
		perPage = pp
	}

	h.l.Debug("fetching songs",
		zap.Any("filters", filters),
		zap.Int("page", page),
		zap.Int("per_page", perPage))

	songs, totalCount, err := h.service.GetAllSong(perPage, page, filters)
	if errors.Is(err, service.ErrParsingTime) {
		h.l.Debug("failed to parse release_date", zap.Error(err))
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid release_date"})
	}
	if err != nil {
		h.l.Error("failed to get all songs", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}
	if len(songs) == 0 {
		h.l.Warn("no songs found")
		return c.JSON(http.StatusNotFound, ErrorResponse{"songs not found"})
	}
	h.l.Info("retrieved songs", zap.Int("count", len(songs)))

	type pagination struct {
		Page    int   `json:"page" example:"1"`
		PerPage int   `json:"per_page" example:"10"`
		Total   int64 `json:"total" example:"100"`
	}
	type successResponse struct {
		Pagination pagination    `json:"pagination"`
		Songs      []models.Song `json:"songs"`
	}

	response := successResponse{
		Songs: songs,
		Pagination: pagination{
			Page:    page,
			PerPage: perPage,
			Total:   totalCount,
		},
	}
	return c.JSON(http.StatusOK, response)
}

// @Summary Получение песни и пагинация текста
// @Description Получение песни и пагинация текста по куплетам
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "song id"
// @Param page query int false "page number" default(1)
// @Param per_page query int false "items per page" default(5)
// @Success 200 {object} api.GetSongHandler.successResponse "received successfully"
// @Failure 404 {object} ErrorResponse "song not found" example:{"error": "song not found"}
// @Failure 422 {object} ErrorResponse "invalid id" example:{"error": "invalid id"}
// @Failure 422 {object} ErrorResponse "invalid page" example:{"error": "invalid page"}
// @Failure 422 {object} ErrorResponse "invalid per_page" example:{"error": "invalid per_page"}
// @Failure 500 {object} ErrorResponse "internal server error" example:{"error": "internal server error"}
// @Router /{id} [get]
func (h *SongHandler) GetSongHandler(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.l.Warn("failed to parse song id",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid id"})
	}
	id := uint(id64)
	h.l.Debug("starting get song", zap.Uint("id", id))

	page := 1
	perPage := 5

	if pageStr := c.QueryParam("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			h.l.Debug("failed to parse page", zap.Error(err))
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid page"})
		}
		page = p
	}

	if perPageStr := c.QueryParam("per_page"); perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err != nil || pp < 1 {
			h.l.Debug("failed to parse per_page", zap.Error(err))
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid per_page"})
		}
		perPage = pp
	}

	h.l.Debug("starting pagination for song text",
		zap.Uint("id", id),
		zap.Int("page", page),
		zap.Int("per_page", perPage))

	song, err := h.service.GetSong(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		h.l.Warn("song not found", zap.Uint("id", id))
		return c.JSON(http.StatusNotFound, ErrorResponse{"song not found"})
	}
	if err != nil {
		h.l.Error("failed to get song", zap.Uint("id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}

	type successResponse struct {
		Verses []string `json:"verses" example:"[\"ooh baby, don't you know i suffer?\", \"ooh baby, can you hear me moan?\"]"`
		Page   int      `json:"page" example:"1"`
		Total  int      `json:"total" example:"10"`
	}

	verses := strings.Split(song.Text, "\n")
	totalVerses := len(verses)
	startIdx := (page - 1) * perPage
	if startIdx >= totalVerses {
		return c.JSON(http.StatusOK, successResponse{Verses: []string{}, Page: page, Total: totalVerses})
	}
	endIdx := startIdx + perPage
	if endIdx > totalVerses {
		endIdx = totalVerses
	}

	h.l.Info("retrieved song text", zap.Uint("id", id))
	return c.JSON(http.StatusOK, successResponse{Verses: verses[startIdx:endIdx], Page: page, Total: totalVerses})
}

// @Summary Обновление песни
// @Description Обновление песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "song id"
// @Param song body models.SongRaw true "song update data"
// @Success 200 {object} api.UpdateSongHandler.successResponse "updated successfully" example:{"success": true}
// @Failure 404 {object} ErrorResponse "song not found" example:{"error": "song not found"}
// @Failure 400 {object} ErrorResponse "invalid data" example:{"error": "invalid data"}
// @Failure 422 {object} ErrorResponse "invalid id" example:{"error": "invalid id"}
// @Failure 422 {object} ErrorResponse "all fields are required" example:{"error": "all fields are required"}
// @Failure 500 {object} ErrorResponse "internal server error" example:{"error": "internal server error"}
// @Router /{id} [put]
func (h *SongHandler) UpdateSongHandler(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.l.Warn("failed to parse song id",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid id"})
	}
	id := uint(id64)
	h.l.Debug("starting update song", zap.Uint("id", id))

	var updatedSong models.SongRaw
	if err = c.Bind(&updatedSong); err != nil {
		h.l.Debug("failed to bind request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{"invalid data"})
	}

	if updatedSong.Song == "" || updatedSong.Group == "" || updatedSong.ReleaseDate == "" ||
		updatedSong.Text == "" || updatedSong.Link == "" {
		h.l.Debug("validation failed: one or more fields are empty",
			zap.Any("data", updatedSong))
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"all fields are required"})
	}

	err = h.service.UpdateSong(id, updatedSong)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		h.l.Warn("song not found", zap.Uint("id", id))
		return c.JSON(http.StatusNotFound, ErrorResponse{"song not found"})
	}
	if err != nil {
		h.l.Error("failed to update song",
			zap.Uint("id", id),
			zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}

	type successResponse struct {
		Success bool `json:"success" example:"true"`
	}
	h.l.Info("song updated successfully", zap.Uint("id", id))
	return c.JSON(http.StatusOK, successResponse{true})
}

// @Summary Удаление песни
// @Description Удаление песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "song id"
// @Success 200 {object} api.DeleteSongHandler.successResponse "deleted successfully" example:{"success": true}
// @Failure 404 {object} ErrorResponse "song not found" example:{"error": "song not found"}
// @Failure 422 {object} ErrorResponse "invalid id" example:{"error": "invalid id"}
// @Failure 500 {object} ErrorResponse "internal server error" example:{"error": "internal server error"}
// @Router /{id} [delete]
func (h *SongHandler) DeleteSongHandler(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.l.Warn("failed to parse song id",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid id"})
	}
	id := uint(id64)
	h.l.Debug("starting delete song", zap.Uint("id", id))
	err = h.service.DeleteSong(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		h.l.Warn("song not found", zap.Uint("id", id))
		return c.JSON(http.StatusNotFound, ErrorResponse{"song not found"})
	}
	if err != nil {
		h.l.Error("failed to delete song", zap.Uint("id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}
	type successResponse struct {
		Success bool `json:"success" example:"true"`
	}
	h.l.Info("song deleted successfully", zap.Uint("id", id))
	return c.JSON(http.StatusOK, successResponse{true})
}
