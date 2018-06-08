package files

import (
	"api/app/models"
	"database/sql"
  "net/http"

	"github.com/gin-gonic/gin"
)

var (
	Fdbs models.FileDBServiceInterface
  Fds models.FileDriveServiceInterface
)

// Configure for files
func Configure(r *gin.Engine, db *sql.DB) {
	Fdbs = &FileDBService{DB: db}
  Fds = &FileDriveService{Client: &http.Client{}}

  r.POST("/auth", Authenticate)
	r.GET("/file/:id", GetFile)
  r.GET("/search-in-drive/:id", SearchInDrive)
	r.POST("/file", PostFile)
	r.GET("/file", GetFiles)
}
