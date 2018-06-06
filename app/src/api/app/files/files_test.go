package files

import (
	//"api/app/mock"
	"net/http"
	"net/http/httptest"
	"testing"
  "strings"
  //"fmt"

	//"api/app/models"

	"github.com/gin-gonic/gin"
)

func TestBindAuthCode(t *testing.T) {
  router := gin.Default()
  Configure(router, nil)

  w := httptest.NewRecorder()
  body := strings.NewReader(`{"auth_code":"test_auth_code"}`)
  r, _ := http.NewRequest("POST", "/auth", body)
  router.ServeHTTP(w, r)
}

func TestGetFilesFromDrive(t *testing.T) {
  router := gin.Default()
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("GET", "/file/100", nil)
  router.ServeHTTP(w, r)
}

func TestSearchForMatch(t *testing.T) {
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

/*
func TestGetID(t *testing.T) {
	router := gin.Default()
	Configure(router, nil)

	// Inject mock into handler
	var fs mock.FileService
	Fs = &fs

	// Mock call.
	fs.FileFn = func(id string) (*models.File, error) {
		if id != "100" {
			t.Fatalf(fmt.Sprintf("unexpected id: %s", id))
		}
		return &models.File{
      ID: "100",
      Titulo: "DaTitulo",
      Descripcion: "archivo de Elnesto",
      DriveID: "abc123"},
      nil
	}

	// Invoke the handler.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/file/100", nil)
	router.ServeHTTP(w, r)

	// Validate mock.
	if !fs.FileInvoked {
		t.Fatal("expected User() to be invoked")
	}
  test_file, err := fs.FileFn("100")
  if err != nil {
    t.Fatal("failed to return file.")
  }
  if test_file.DriveID != "abc123"{
    fmt.Print("DriveID:")
    fmt.Print(test_file.DriveID)
    t.Fatal("expected drive ID to be 100")
  }

}
*/
