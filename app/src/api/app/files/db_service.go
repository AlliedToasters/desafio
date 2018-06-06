package files

import (
	"api/app/models"
	"database/sql"
	"strconv"
  "google.golang.org/api/drive/v3"
)

// FileService ...
type FileService struct {
  //tok *oauth2.Token
  d_srv *drive.Service
	DB *sql.DB
}

// File ...
func (s *FileService) File(id string) (*models.File, error) {
	var f models.File
	row := s.DB.QueryRow(
    `SELECT id, titulo, descripcion, drive_id
    FROM files
    WHERE id = ?`,
    id)
	if err := row.Scan(&f.ID, &f.Titulo, &f.Descripcion, &f.DriveID); err != nil {
		return nil, err
	}
	return &f, nil
}

// Files ...
func (s *FileService) Files() ([]*models.File, error) {
	rows, err := s.DB.Query(`SELECT id, titulo, descripcion, drive_id FROM files`)
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
func (s *FileService) CreateFile(f *models.File) error {
	stmt, err := s.DB.Prepare(`INSERT INTO files(titulo, descripcion, drive_id)
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
