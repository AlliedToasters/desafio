package files

import (
	"api/app/mock"
  "api/app/models"
	"testing"
	"net/http"
	"net/http/httptest"
  "strings"
  "encoding/json"

	//"api/app/models"

	"github.com/gin-gonic/gin"
)

func getFilesFromDriveFn() ([]*models.File, error) {
  var result []*models.File
  var file models.File
  file.Titulo = "Pagos a prov"
  file.Descripcion = "Tengo que hacer un pago"
  file.DriveID = "abc123"
  result = append(result, &file)
  return result, nil
}

func storeFilesFn(files []*models.File) error {
  if files[0].DriveID != "abc123" || files[0].Titulo != "Pagos a prov" {
    var err error
    return err
  }
  return nil
}

func TestAuthenticate(t *testing.T) {
  router := gin.Default()
  Configure(router, nil)

  //Mock drive interface
  var fds mock.FileDriveService
  Fds = &fds
  //Mock db interface
  var fdbs mock.FileDBService
  Fdbs = &fdbs

  fds.GetClientFn = func(code string) error {
    if code != "abc123" {
      t.Log("Expected authentication code 'abc123', got: \n")
      t.Fatal(code)
    }
    return nil
  }

  fds.GetFilesFromDriveFn = getFilesFromDriveFn
  fdbs.StoreFilesFn = storeFilesFn

  //Mock user call
  msg := `{"auth_code":"abc123"}`
	w := httptest.NewRecorder()
  reader := strings.NewReader(msg)
	r, _ := http.NewRequest("POST", "/auth", reader)
	router.ServeHTTP(w, r)

  if w.Code != 200 {
    t.Log("Bad code, expected 200, got: ")
    t.Fatal(w.Code)
  }

  if !fds.GetClientInvoked {
    t.Fatal("function fdbs.GetClient not invoked.")
  }

  if !fds.GetFilesFromDriveInvoked {
    t.Fatal("function dfs.GetFilesFromDrive not invoked.")
  }

  if !fdbs.StoreFilesInvoked {
    t.Fatal("functions dfs.GetFilesFromDrive, fdbs.StoreFile not invoked.")
  }
}


func TestGetFile(t *testing.T) {
  router := gin.Default()
  Configure(router, nil)

  //Mock drive interface
  var fds mock.FileDriveService
  Fds = &fds
  //Mock db interface
  var fdbs mock.FileDBService
  Fdbs = &fdbs


  fds.GetFilesFromDriveFn = getFilesFromDriveFn
  fdbs.StoreFilesFn = storeFilesFn

  fdbs.FileFn = func(fileID string) (*models.File, error) {
    if fileID != "100" {
      t.Fatal("Expected fileID to be '100'")
    }
    var file models.File
    file.ID = "100"
    file.Titulo = "Pagos a prov"
    file.Descripcion = "Tengo que hacer un pago"
    file.DriveID = "abc123"
    return &file, nil
  }

  //Mock user call
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/file/100", nil)
	router.ServeHTTP(w, r)
  if !fdbs.FileInvoked {
    t.Fatal("function fdbs.File not invoked.")
  }

  if w.Code != 200 {
    t.Log("Bad code, expected 200, got: ")
    t.Fatal(w.Code)
  }

  result_file := models.File{}
  result_bytes := w.Body.Bytes()
  err := json.Unmarshal(result_bytes, &result_file)
  if err != nil {
    t.Fatal(err.Error())
  }
  if result_file.DriveID != "abc123" {
    t.Log("Got unexpected file returned: ")
    t.Fatal(result_file)
  }
}

func TestSearchInDrive(t *testing.T) {
    router := gin.Default()
    Configure(router, nil)

    //Mock db interface
    var fdbs mock.FileDBService
    Fdbs = &fdbs
    //Mock drive interface
    var fds mock.FileDriveService
    Fds = &fds

    fdbs.GetDriveIDFn = func(fileID string) (string, error) {
      if fileID != "100" {
        t.Fatal("Expected fileID to be '100'")
      }
      return "abc123", nil
    }

    fds.GetWordQueryFn = func(word string) ([]*string, error) {
      if word != "dev" {
        t.Log("Expected word 'dev', got: ")
        t.Fatal(word)
      }
      var matches []*string
      res1 := "abc123"
      matches = append(matches, &res1)
      return matches, nil
    }

    //Mock user call
  	w := httptest.NewRecorder()
  	r, _ := http.NewRequest("GET", "/search-in-drive/100?word=dev", nil)
  	router.ServeHTTP(w, r)
    if !fdbs.GetDriveIDInvoked {
      t.Fatal("function fdbs.GetDriveID not invoked.")
    }
    if !fds.GetWordQueryInvoked {
      t.Fatal("function fds.GetWordQuery not invoked.")
    }
    if w.Code != 200 {
      t.Log("Bad code, expected 200, got: ")
      t.Fatal(w.Code)
    }
}

func TestsearchForMatch(t *testing.T) {
  var matches []*string
  entry := "abc123"
  matches = append(matches, &entry)
  entry2 := "def456"
  matches = append(matches, &entry2)
  entry3 := "xyz999"
  matches = append(matches, &entry3)
  match := "def456"
  nonmatch := "hhd918"
  found, err := searchForMatch(match, matches)
  if err != nil {
    t.Fatal("Error matching: ", err.Error())
  }
  if !found {
    t.Fatal("Failed to match: ", match, matches)
  }
  found, err = searchForMatch(nonmatch, matches)
  if err != nil {
    t.Fatal("Error matching: ", err.Error())
  }
  if found {
    t.Fatal("False match: ", match, matches)
  }
}
