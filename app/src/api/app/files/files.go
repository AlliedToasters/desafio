package files

import (
	"api/app/models"
	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	//Is file service
	Fs models.FileServiceInterface
)

// Configure for files
func Configure(r *gin.Engine, db *sql.DB) {
	Fs = &FileService{DB: db}

	r.GET("/file/:id", GetFile)
	r.POST("/file", PostFile)
	r.GET("/file", GetFiles)
}
