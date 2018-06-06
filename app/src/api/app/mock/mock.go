package mock

import "api/app/models"

// ItemService ...
type ItemService struct {
	ItemFn      func(id string) (*models.Item, error)
	ItemInvoked bool

	ItemsFn      func() ([]*models.Item, error)
	ItemsInvoked bool

	CreateItemFn      func(i *models.Item) error
	CreateItemInvoked bool

	DeleteItemFn      func(id string) (*models.Item, error)
	DeleteItemInvoked bool
}

// Item ...
func (is *ItemService) Item(id string) (*models.Item, error) {
	is.ItemInvoked = true
	return is.ItemFn(id)
}

// Items ...
func (is *ItemService) Items() ([]*models.Item, error) {
	is.ItemsInvoked = true
	return is.ItemsFn()
}

// CreateItem ...
func (is *ItemService) CreateItem(i *models.Item) error {
	is.CreateItemInvoked = true
	return is.CreateItemFn(i)
}

// DeleteItem ...
func (is *ItemService) DeleteItem(id string) (*models.Item, error) {
	is.DeleteItemInvoked = true
	return is.DeleteItemFn(id)
}

type FileService struct {
	FileFn      func(id string) (*models.File, error)
	FileInvoked bool

	FilesFn      func() ([]*models.File, error)
	FilesInvoked bool

	CreateFileFn      func(i *models.File) error
	CreateFileInvoked bool
}

// File ...
func (fs *FileService) File(id string) (*models.File, error) {
	fs.FileInvoked = true
	return fs.FileFn(id)
}

// Files ...
func (fs *FileService) Files() ([]*models.File, error) {
	fs.FilesInvoked = true
	return fs.FilesFn()
}

// CreateFile ...
func (fs *FileService) CreateFile(f *models.File) error {
	fs.CreateFileInvoked = true
	return fs.CreateFileFn(f)
}



/*
type MockDriveMeta struct{
  Id       string
  Name     string
}
*/
