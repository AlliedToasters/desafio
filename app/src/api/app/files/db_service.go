package files

import (
	"api/app/models"
	"database/sql"
	"strconv"
)

// File database service ...
type FileDBService struct {
	DB *sql.DB
}

// File ...
func (fdbs *FileDBService) File(id string) (*models.File, error) {
	var f models.File
	row := fdbs.DB.QueryRow(
    `SELECT id, titulo, descripcion, drive_id
    FROM files
    WHERE id = ?`,
    id)
	if err := row.Scan(&f.ID, &f.Titulo, &f.Descripcion, &f.DriveID); err != nil {
		return nil, err
	}
	return &f, nil
}

//Lookup drive ID given :id
func (fdbs *FileDBService) GetDriveID(fileID string) (string, error) {
  file, err := Fdbs.File(fileID)
  return file.DriveID, err
}

// Files ...
func (fdbs *FileDBService) Files() ([]*models.File, error) {
	rows, err := fdbs.DB.Query(`SELECT id, titulo, descripcion, drive_id FROM files`)
  defer rows.Close()
	if err != nil {
		return nil, err
	}
	var files []*models.File
	for rows.Next() {
		var f models.File
		if err := rows.Scan(&f.ID, &f.Titulo, &f.Descripcion, &f.DriveID); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil
}

// CreateFile ...
func (fdbs *FileDBService) CreateFile(f *models.File) error {
	stmt, err := fdbs.DB.Prepare(`INSERT INTO files(titulo, descripcion, drive_id)
  values(?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(f.Titulo, f.Descripcion, f.DriveID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	f.ID = strconv.Itoa(int(id))
	return nil
}

func (fdbs *FileDBService) StoreFiles(files []*models.File) error {
  for _, file := range files {
    fdbs.CreateFile(file)
  }
  return nil
}
