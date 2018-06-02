package items

import (
	"api/app/models"
	"database/sql"
	"strconv"
	"fmt"
)

// ItemService ...
type ItemService struct {
	DB *sql.DB
}

// Item ...
func (s *ItemService) Item(id string) (*models.Item, error) {
	var i models.Item
	row := s.DB.QueryRow(`SELECT id, name, description FROM items WHERE id = ?`, id)
	if err := row.Scan(&i.ID, &i.Name, &i.Description); err != nil {
		return nil, err
	}
	return &i, nil
}

// Items ...
func (s *ItemService) Items() ([]*models.Item, error) {
	rows, err := s.DB.Query(`SELECT id, name, description FROM items`)
  defer rows.Close()
	if err != nil {
		return nil, err
	}
	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// CreateItem ...
func (s *ItemService) CreateItem(i *models.Item) error {
	stmt, err := s.DB.Prepare(`INSERT INTO items(name,description) values(?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(i.Name, i.Description)
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

// DeleteItem ...
func (s *ItemService) DeleteItem(id string) (*models.Item, error) {
	//Verify entry exists...
	item, err := s.Item(id)
	if err != nil {
		return nil, err
	}
	stmt, err := s.DB.Prepare(`DELETE FROM items WHERE id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	response, err := stmt.Exec(id)
	fmt.Print(response)
	if err != nil {
		return nil, err
	}
	return item, nil
}
