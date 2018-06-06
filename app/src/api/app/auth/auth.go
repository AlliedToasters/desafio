package auth

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "os"
        //"reflect"

        "api/app/data"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
)

//Invoked when authorization fails
type authError struct {
  url string
}

func (err authError) Error() string {
  return err.url
}


func clientFromFile(config *oauth2.Config) (*http.Client, error) {
  tokenFile := "token.json"
  tok, err := tokenFromFile(tokenFile)
  return config.Client(context.Background(), tok), err
}

func GetAuthCodeURL() string {
  config, err := getConfig()
  if err != nil {
    log.Fatalf("Could not get config.")
  }
  authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  return authURL
}

func GetClient() (*http.Client, error) {
  config, err := getConfig()
  fmt.Print("Looking for token in file...")
  client, err := clientFromFile(config)
  if err != nil {
    var auth_err authError
    auth_err.url = GetAuthCodeURL()
    return client, auth_err
  }
  return client, err
}

func getConfig() (*oauth2.Config, error) {
    b, err := data.Asset("src/api/app/data/client_secret.json")
    if err != nil {
            log.Fatalf("Unable to read client secret file: %v", err)
    }
    config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
    return config, err
}

func AuthToToken(authCode string) (*oauth2.Token, error) {
    config, err := getConfig()
    if err != nil {
            log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    tok, err := config.Exchange(oauth2.NoContext, authCode)
    if err != nil {
            log.Panic("Unable to retrieve token from web")
            return tok, err
    }
    return tok, err
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    defer f.Close()
    if err != nil {
            return nil, err
    }
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    defer f.Close()
    if err != nil {
            log.Fatalf("Unable to cache oauth token: %v", err)
    }
    json.NewEncoder(f).Encode(token)
}

func DoPeek() {
    client, err := GetClient()
    srv, err := drive.New(client)
    if err != nil {
            log.Fatalf("Unable to retrieve Drive client: %v", err)
    }

    r, err := srv.Files.List().PageSize(10).
            Fields("nextPageToken, files(id, name)").Do()
    if err != nil {
            log.Fatalf("Unable to retrieve files: %v", err)
    }
    fmt.Println("Files:")
    if len(r.Files) == 0 {
            fmt.Println("No files found.")
    } else {
            for _, i := range r.Files {
                    fmt.Printf("%s (%s)\n", i.Name, i.Id)
            }
    }
}
