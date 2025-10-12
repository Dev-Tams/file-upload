package actions

import (
	"math"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"totalPages"`
	NextPage   string      `json:"nextPage,omitempty"`
	PrevPage   string      `json:"prevPage,omitempty"`
	Data       any 			`json:"data"`
}

func Paginate[T any](c *gin.Context, db *gorm.DB, model T) (Pagination, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var total int64

	if err := db.Model(model).Count(&total).Error; err != nil {
		return Pagination{}, err
	}

	var data []T
	if err := db.Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return Pagination{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Build base URL
	baseURL := &url.URL{
		Scheme: c.Request.URL.Scheme,
		Host:   c.Request.Host,
		Path:   c.FullPath(),
	}

	q := c.Request.URL.Query()

	var nextPage, prevPage string

	if page < totalPages {
		q.Set("page", strconv.Itoa(page+1))
		q.Set("limit", strconv.Itoa(limit))
		baseURL.RawQuery = q.Encode()
		nextPage = baseURL.String()
	}

	if page > 1 {
		q.Set("page", strconv.Itoa(page-1))
		q.Set("limit", strconv.Itoa(limit))
		baseURL.RawQuery = q.Encode()
		prevPage = baseURL.String()
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		NextPage:   nextPage,
		PrevPage:   prevPage,
		Data:       data,
	}, nil
}
