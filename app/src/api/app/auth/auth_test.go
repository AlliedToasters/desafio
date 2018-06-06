package auth

import (
  "testing"
  "fmt"
)

func TestGetConfig(t *testing.T) {
  config, err := getConfig()
  if err != nil {
    t.Fatal(fmt.Sprintf("Could not get config. Error: %s", err))
  }
  t.Log("get config: ")
  t.Log(config)
}

func TestGetURL(t *testing.T) {
  authURL := GetAuthCodeURL()
  if authURL == "" {
    t.Fatal("Did not retrieve authorization URL")
  }
  t.Log("Got url: ")
  t.Log(authURL)
}

func TestGetClient(t *testing.T) {
  _, err := GetClient()
  t.Log("GetClient returned error: ")
  t.Log(err)
}

func TestHaveToken(t *testing.T) {
  if HaveToken() {
    t.Fatal("Should not have token saved on disk.")
  }
}
