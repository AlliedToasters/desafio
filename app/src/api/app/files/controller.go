package files

import (
	"net/http"
	"strings"

	"api/app/models"

	"github.com/gin-gonic/gin"
    "fmt"
)

// GetFile ...
func GetFile(c *gin.Context) {
    fmt.Print(c.Param("id"))
	fileID := strings.TrimSpace(c.Param("id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}

	file, err := Fs.File(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "find_error", "description": err.Error()})
		return
	}
	c.JSON(200, file)
	return
}

// GetFiles ...
func GetFiles(c *gin.Context) {
  fmt.Print("Getting files...")
	files, err := Fs.Files()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "find_error", "description": err.Error()})
		return
	}
	c.JSON(200, files)
	return
}

// PostFile ...
func PostFile(c *gin.Context) {
	f := &models.File{}
	if err := c.BindJSON(f); c.Request.ContentLength == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind_error", "description": err.Error()})
		return
	}
	err := Fs.CreateFile(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "save_error", "description": err.Error()})
		return
	}
	c.JSON(201, f)
}
