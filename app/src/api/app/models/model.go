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

// Interface...
type FileServiceInterface interface {
	File(id string) (*File, error)
	Files() ([]*File, error)
	CreateFile(i *File) error
}
