package files

import (
	"api/app/models"
	"database/sql"
	"strconv"
)

// FileService ...
type FileService struct {
	DB *sql.DB
}

// File ...
func (s *FileService) File(id string) (*models.File, error) {
	var f models.File
	row := s.DB.QueryRow(`SELECT id, titulo, descripcion FROM files WHERE id = ?`, id)
	if err := row.Scan(&f.ID, &f.Titulo, &f.Descripcion); err != nil {
		return nil, err
	}
	return &f, nil
}

// Files ...
func (s *FileService) Files() ([]*models.File, error) {
	rows, err := s.DB.Query(`SELECT id, titulo, descripcion FROM files`)
  defer rows.Close()
	if err != nil {
		return nil, err
	}
	var files []*models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Titulo, &file.Descripcion); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return files, nil
}

// CreateFile ...
func (s *FileService) CreateFile(i *models.File) error {
	stmt, err := s.DB.Prepare(`INSERT INTO files(titulo,descripcion) values(?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(i.Titulo, i.Descripcion)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	i.ID = strconv.Itoa(int(id))
	return nil
}
