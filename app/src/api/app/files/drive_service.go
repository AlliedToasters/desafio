package files

import (
	"net/http"
	"strings"
  "fmt"

	"api/app/models"
  "api/app/auth"

	//"github.com/gin-gonic/gin"
  "google.golang.org/api/drive/v3"
  "google.golang.org/api/googleapi"
)

// File drive service ...
type FileDriveService struct {
  Client *http.Client
}

func (fds *FileDriveService) GetClient(code string) error {
  tok, err := auth.AuthToToken(code)
  if err != nil {
    var auth_err auth.AuthError
    return auth_err
  }
  save_path := "token.json"
  auth.SaveToken(save_path, tok)
  client, err := auth.GetClient()
  fds.Client = client
  return nil
}

//Handles responses when authentication is needed
func (fds *FileDriveService) GetAuthenticateURL() string {
  return auth.GetAuthCodeURL()
}

func (fds *FileDriveService) GetFilesFromDrive() ([]*models.File, error) {
  client := fds.Client
  srv, err := drive.New(client)
  if err != nil {
    return nil, err
  }
  r, err := srv.Files.List().PageSize(10).
          Fields("nextPageToken, files(id, name, description)").Do()
  if err != nil{
    return nil, err
  }
  var result []*models.File
  if len(r.Files) == 0 {
    return nil, nil
  } else {
    for _, i := range r.Files {
      var file models.File
      file.Titulo = i.Name
      file.DriveID = i.Id
      file.Descripcion = i.Description
      result = append(result, &file)
    }
  }
  return result, nil
}

//Performs query, returns a list of matching fileIDs
func (fds *FileDriveService) GetWordQuery(word string) ([]*string, error) {
  var result []*string
  client := fds.Client
  srv, err := drive.New(client)
  if err != nil {
    return result, err
  }
  q_query := fmt.Sprintf(`fullText contains '"%s"'`, word)
  r, err := srv.Files.List().Fields("files(id)").Q(q_query).Do()
  if err != nil {
    return result, err
  }
  if len(r.Files) == 0 {
    fmt.Println("No matching files found.")
    return result, nil
  }
  for _, i := range r.Files {
    result = append(result, &i.Id)
  }
  return result, nil
}


//Posts file to drive
func (fds *FileDriveService) DrivePostFile(f *models.File) (*models.File, error) {
  client := fds.Client
  srv, err := drive.New(client)
  if err != nil {
    return nil, err
  }
  drive_file := &drive.File{
    Name:     f.Titulo,
    MimeType: "text/plain",
  }
  reader := strings.NewReader("")
  //r, err := srv.Files().Insert(drive_file)
  r, err := srv.Files.Create(drive_file).Media(reader, googleapi.ContentType("text/plain")).Do()
  if err != nil {
    return f, err
  }
  f.DriveID = r.Id
  return f, nil
}
