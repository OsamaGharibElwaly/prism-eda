package api

import (
	"errors"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"prismeda/internal/eda"
	"prismeda/internal/service"

	"github.com/gin-gonic/gin"
)

const Version = "1.0.0"

var started = time.Now()

func Health(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"version":   Version,
			"go":        runtime.Version(),
			"uptime_s":  int64(time.Since(started).Seconds()),
			"datasets":  svc.Count(),
		})
	}
}

func UploadDataset(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			fail(c, http.StatusBadRequest, "attach a CSV file using the 'file' field")
			return
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".csv" && ext != ".tsv" && ext != ".txt" {
			fail(c, http.StatusUnsupportedMediaType, "only .csv, .tsv and .txt files are supported")
			return
		}
		if file.Size > svc.MaxBytes() {
			fail(c, http.StatusRequestEntityTooLarge, "file exceeds the upload limit")
			return
		}
		f, err := file.Open()
		if err != nil {
			fail(c, http.StatusBadRequest, "could not read uploaded file")
			return
		}
		defer f.Close()

		profile, err := svc.Ingest(file.Filename, f)
		if err != nil {
			respondErr(c, err)
			return
		}
		c.JSON(http.StatusCreated, profile)
	}
}

func LoadSample(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		profile, err := svc.LoadSample()
		if err != nil {
			respondErr(c, err)
			return
		}
		c.JSON(http.StatusCreated, profile)
	}
}

func GetProfile(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		profile, err := svc.Profile(c.Param("id"))
		if err != nil {
			respondErr(c, err)
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}

func GetPreview(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
		if err != nil || limit < 1 || limit > 100 {
			fail(c, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
		preview, err := svc.Preview(c.Param("id"), limit)
		if err != nil {
			respondErr(c, err)
			return
		}
		c.JSON(http.StatusOK, preview)
	}
}

func GetScatter(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		x, y := c.Query("x"), c.Query("y")
		if x == "" || y == "" {
			fail(c, http.StatusBadRequest, "query parameters x and y are required")
			return
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "1500"))
		if err != nil || limit < 1 || limit > 5000 {
			fail(c, http.StatusBadRequest, "limit must be between 1 and 5000")
			return
		}
		points, err := svc.Scatter(c.Param("id"), x, y, limit)
		if err != nil {
			fail(c, http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"x": x, "y": y, "points": points})
	}
}

func respondErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		fail(c, http.StatusNotFound, "dataset not found")
	case errors.Is(err, eda.ErrInvalidCSV), errors.Is(err, eda.ErrMaxRows):
		fail(c, http.StatusBadRequest, err.Error())
	default:
		fail(c, http.StatusInternalServerError, err.Error())
	}
}

func fail(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
