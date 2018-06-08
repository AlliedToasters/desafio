package models

// Item ...
type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ItemServiceInterface ...
type ItemServiceInterface interface {
	Item(id string) (*Item, error)
	Items() ([]*Item, error)
	CreateItem(i *Item) error
	DeleteItem(id string) (*Item, error)
}

// File type.
type File struct {
	ID          string `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
  DriveID     string `json:"drive_id"`
}

// Database Interface...
type FileDBServiceInterface interface {
	File(id string) (*File, error)
  GetDriveID(id string) (string, error)
	Files() ([]*File, error)
	CreateFile(f *File) error
	StoreFiles(files []*File) error
}

// Drive Interface...
type FileDriveServiceInterface interface {
  GetClient(code string) error
  GetAuthenticateURL() string
  GetFilesFromDrive() ([]*File, error)
	GetWordQuery(word string) ([]*string, error)
  DrivePostFile(f *File) error
}
