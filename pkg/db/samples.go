package db

import (
	"database/sql"
	"fmt"

	"github.com/Arash-Afshar/gohtmx/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func AddSample(db *sql.DB, sample *models.Sample) error {
	sqlString := `INSERT INTO samples(name) VALUES (?)`
	statement, err := db.Prepare(sqlString)
	if err != nil {
		return fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	res, err := statement.Exec(sample.Name)
	if err != nil {
		return fmt.Errorf("executing %s with (%s): %v", sqlString, sample.Name, err)
	}
	sample.Id, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %v", err)
	}
	return nil
}

func SampleByName(db *sql.DB, name string) (*models.Sample, error) {
	sqlString := `SELECT id FROM samples WHERE name = ?`
	statement, err := db.Prepare(sqlString)
	if err != nil {
		return nil, fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	row, err := statement.Query(name)
	defer row.Close()
	if row.Next() {
		var id int64
		if err = row.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		if row.Next() {
			return nil, fmt.Errorf("more than one entry found")
		}
		return &models.Sample{Id: id, Name: name}, nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}

func DeleteSample(db *sql.DB, sample *models.Sample) error {
	sqlString := `DELETE FROM samples WHERE id = ?`
	statement, err := db.Prepare(sqlString)
	if err != nil {
		return fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	_, err = statement.Exec(sample.Id)
	if err != nil {
		return fmt.Errorf("executing %s with (%d): %v", sqlString, sample.Id, err)
	}
	return nil
}

func ListSamples(db *sql.DB) ([]*models.Sample, error) {
	sqlString := `SELECT id, name FROM samples`
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	defer rows.Close()
	samples := make([]*models.Sample, 0)
	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		samples = append(samples, &models.Sample{
			Id:   id,
			Name: name,
		})
	}
	return samples, nil
}
