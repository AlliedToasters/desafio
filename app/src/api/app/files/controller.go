package files

import (
	"net/http"
	"strings"
  "fmt"

	"api/app/models"
  "api/app/auth"

	"github.com/gin-gonic/gin"
)

// Retrieves file from db
func GetFile(c *gin.Context) {
  have_token := auth.HaveToken()
  if !have_token {
    promptAuthenticate(c)
    return
  }
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

//Lookup drive ID given :id
func getDriveID(fileID string) (string, error) {
  file, err := Fs.File(fileID)
  return file.DriveID, err
}

// Retrieve all files in db
func GetFiles(c *gin.Context) {
  have_token := auth.HaveToken()
  if !have_token {
    promptAuthenticate(c)
    return
  }
	files, err := Fs.Files()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "find_error", "description": err.Error()})
		return
	}
  if len(files) == 0 {
    c.JSON(404, gin.H{"error": "find_error", "description": "no files found."})
    return
  }
	c.JSON(200, files)
	return
}

// Create new file in db
func PostFile(c *gin.Context) {
	f := &models.File{}
	if err := c.BindJSON(f); c.Request.ContentLength == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind_error", "description": err.Error()})
		return
	}
	err := storeFile(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "save_error", "description": err.Error()})
		return
	}
	c.JSON(201, f)
}

// Helper function: stores a file object in db
func storeFile(file *models.File) error {
	err := Fs.CreateFile(file)
  return err
}
