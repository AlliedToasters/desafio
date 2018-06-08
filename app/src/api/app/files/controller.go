package files

import (
	"net/http"
	"strings"
  "fmt"

	"api/app/models"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
  code, err := getAuthCode(c)
  if err != nil {
    c.JSON(400, gin.H{"error":"authentication_error", "description":err.Error()})
    return
  }
  err = Fds.GetClient(code)
  if err != nil {
    url := Fds.GetAuthenticateURL()
    c.JSON(401, gin.H{"error":"unauthorized", "description":url})
    return
  }
  err = getFiles(c)
  if err != nil {
    c.JSON(200, gin.H{"success":"authentication success.","getFiles_error":err.Error()})
    return
  }
  c.JSON(200, gin.H{"success":"authentication_success"})
  return
}

//Handles extracting authentication from query in JSON format
type authCode struct {
  Value       string `json:"auth_code"`
}

//invokes authCode and BindJSON, returns code
func getAuthCode(c *gin.Context) (string, error) {
  var code authCode
  if err := c.ShouldBindJSON(&code); err != nil {
    return "", err
  }
  return code.Value, nil
}

func getFiles(c *gin.Context) error {
  files, err := Fds.GetFilesFromDrive()
  if err != nil {
    return err
  }
  if len(files) == 0 {
    return nil
  }
  if err := Fdbs.StoreFiles(files); err != nil {
    return err
  }
  return nil

}

// Retrieves file from db
func GetFile(c *gin.Context) {
  err := getFiles(c)
  if err != nil {
    url := Fds.GetAuthenticateURL()
    c.JSON(401, gin.H{"error":"unauthorized", "description":url})
    return
  }
  fmt.Print(c.Param("id"))
	fileID := strings.TrimSpace(c.Param("id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}
	file, err := Fdbs.File(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "find_error", "description": err.Error()})
		return
	}
	c.JSON(200, file)
	return
}

//Handles search task
func SearchInDrive(c *gin.Context) {
	fileID := strings.TrimSpace(c.Param("id"))
  word := c.Request.URL.Query()["word"][0]
  driveID, err := Fdbs.GetDriveID(fileID)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":"db_error","description":"file not found"})
    return
  }
  matches, err := Fds.GetWordQuery(word)
  if err != nil {
    c.JSON(401, gin.H{"error":"drive_error","description":err.Error()})
    return
  }
  found, err := searchForMatch(driveID, matches)
  if !found {
    c.Status(404)
    return
  }
  if err != nil {
    c.JSON(403, gin.H{"error":"search error","description":err.Error()})
    return
  }
  c.Status(200)
  return
}

//Checks if fileID is in returned list
func searchForMatch(driveID string, matches []*string) (bool, error) {
  for _, v := range matches {
    if driveID == *v {
      return true, nil
    }
  }
  return false, nil
}


// Retrieve all files in db
func GetFiles(c *gin.Context) {
  err := getFiles(c)
  if err != nil {
    url := Fds.GetAuthenticateURL()
    c.JSON(401, gin.H{"error":"unauthorized", "description":url})
    return
  }
	files, err := Fdbs.Files()
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
		c.JSON(400, gin.H{"error": "parameter error", "description": err.Error()})
		return
	}
  f, err := Fds.DrivePostFile(f)
  if err != nil {
    c.JSON(500, gin.H{"error":"could not create", "description":err.Error()})
    return
  }
	err = Fdbs.CreateFile(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "save_error", "description": err.Error()})
		return
	}
	c.JSON(200, f)
  return
}
