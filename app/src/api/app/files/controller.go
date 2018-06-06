package files

import (
	"net/http"
	"strings"
  "fmt"
  "io/ioutil"
  "regexp"

	"api/app/models"
  "api/app/auth"

	"github.com/gin-gonic/gin"
  "google.golang.org/api/drive/v3"
)

type authCode struct {
  Value       string `json:"auth_code"`
}

// Authenticate...
func getAuthCode(c *gin.Context) (string, error) {
  var code authCode
  if err := c.ShouldBindJSON(&code); err != nil {
    return "", err
  }
  return code.Value, nil
}

func Authenticate(c *gin.Context) {
  code, err := getAuthCode(c)
  if err != nil {
    c.JSON(404, gin.H{"error":"authentication_error", "description":err.Error()})
  }
  tok, err := auth.AuthToToken(code)
  if err != nil {
    url := auth.GetAuthCodeURL()
    c.JSON(401, gin.H{"error":"authentication_error", "auth_url":url})
    return
  }
  save_path := "token.json"
  auth.SaveToken(save_path, tok)
  c.JSON(200, gin.H{"success":"authentication_success"})
  return
}

// Retrives file metadata from drive, adds them to database
func getFilesFromDrive(c *gin.Context) error {
  client, err := auth.GetClient()
  if err != nil {
    //c.JSON(401, gin.H{"error":"authentication_error", "description":err.Error()})
    return err
  }
  srv, err := drive.New(client)
  if err != nil {
    //c.JSON(404, gin.H{"error":"find_error", "description":err.Error()})
    return err
  }
  r, err := srv.Files.List().PageSize(10).
          Fields("nextPageToken, files(id, name, description)").Do()
  if err != nil{
    //c.JSON(404, gin.H{"error":"find_error", "description":err.Error()})
    return err
  }
  if len(r.Files) == 0 {
    fmt.Print("No files found in drive.")
    return err
  } else {
    for _, i := range r.Files {
      var file models.File
      file.Titulo = i.Name
      file.DriveID = i.Id
      file.Descripcion = i.Description
      storeFile(&file)
    }
  }
  return nil
}

// GetFile ...
func GetFile(c *gin.Context) {
    fmt.Print(c.Param("id"))
	fileID := strings.TrimSpace(c.Param("id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}
  err := getFilesFromDrive(c)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "description": err.Error()})
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

func getDriveID(fileID string) (string, error) {
  file, err := Fs.File(fileID)
  return file.DriveID, err
}

func SearchInDrive(c *gin.Context) {
	fileID := strings.TrimSpace(c.Param("id"))
  word := c.Request.URL.Query()["word"][0]
  driveID, err := getDriveID(fileID)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":"db_error","description":"file not found"})
    return
  }
  content, err := getFileContent(driveID)
  if err!= nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":"drive_error","description":err.Error()})
    return
  }
  found, err := searchFile(content, word)
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

func searchFile(content string, word string) (bool, error) {
  exp := fmt.Sprintf(" %s ", word)
  match, err := regexp.Match(exp, []byte(content))
  return match, err
}

func getFileContent(driveID string) (string, error) {
  client, err := auth.GetClient()
  if err != nil {
    return "", err
  }
  srv, err := drive.New(client)
  if err != nil {
    return "", err
  }
  r, err := srv.Files.Export(driveID, "text/plain").Download()//.Do("alts=media")
  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  bodyString := string(body)
  return bodyString, err
}

// GetFiles ...
func GetFiles(c *gin.Context) {
  err := getFilesFromDrive(c)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "description": err.Error()})
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

// PostFile ...
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

func storeFile(file *models.File) error {
	err := Fs.CreateFile(file)
  return err
}