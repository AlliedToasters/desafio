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

type FileDBService struct {
	FileFn      func(id string) (*models.File, error)
	FileInvoked bool

  GetDriveIDFn func(string) (string, error)
  GetDriveIDInvoked bool

	FilesFn      func() ([]*models.File, error)
	FilesInvoked bool

	CreateFileFn      func(f *models.File) error
	CreateFileInvoked bool

  StoreFilesFn func(files []*models.File) error
  StoreFilesInvoked bool
}

func (fdbs *FileDBService) File(id string) (*models.File, error) {
	fdbs.FileInvoked = true
	return fdbs.FileFn(id)
}

func (fdbs *FileDBService) GetDriveID(id string) (string, error) {
	fdbs.GetDriveIDInvoked = true
	return fdbs.GetDriveIDFn(id)
}

func (fdbs *FileDBService) Files() ([]*models.File, error) {
	fdbs.FilesInvoked = true
	return fdbs.FilesFn()
}

func (fdbs *FileDBService) CreateFile(f *models.File) error {
	fdbs.CreateFileInvoked = true
	return fdbs.CreateFileFn(f)
}

func (fdbs *FileDBService) StoreFiles(files []*models.File) error {
	fdbs.StoreFilesInvoked = true
	return fdbs.StoreFilesFn(files)
}

type FileDriveService struct {

  GetClientFn func(code string) error
  GetClientInvoked bool

	GetAuthenticateURLFn      func() string
	GetAuthenticateURLInvoked bool

	GetFilesFromDriveFn      func() ([]*models.File, error)
	GetFilesFromDriveInvoked bool

  GetWordQueryFn func(word string) ([]*string, error)
  GetWordQueryInvoked bool

  DrivePostFileFn func(f *models.File) (*models.File, error)
  DrivePostFileInvoked bool
}

func (fds *FileDriveService) GetClient(code string) error {
	fds.GetClientInvoked = true
	return fds.GetClientFn(code)
}

func (fds *FileDriveService) GetAuthenticateURL() string {
	fds.GetAuthenticateURLInvoked = true
	return fds.GetAuthenticateURLFn()
}

func (fds *FileDriveService) GetFilesFromDrive() ([]*models.File, error) {
	fds.GetFilesFromDriveInvoked = true
	return fds.GetFilesFromDriveFn()
}

func (fds *FileDriveService) GetWordQuery(word string) ([]*string, error) {
	fds.GetWordQueryInvoked = true
	return fds.GetWordQueryFn(word)
}

func (fds *FileDriveService) DrivePostFile(f *models.File) (*models.File, error) {
	fds.DrivePostFileInvoked = true
	return fds.DrivePostFileFn(f)
}
