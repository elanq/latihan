package connection

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/latihan/model"
)

type MySQL struct {
	db *sql.DB
}

func (m *MySQL) Insert(ctx context.Context, query string) error {
	result, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println(rowAffected, " row(s) affected")
	return nil
}

func (m *MySQL) Select(ctx context.Context, query string) ([]model.Student, error) {
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var student model.Student
	var students []model.Student

	for rows.Next() {
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Class)
		if err != nil {
			continue
		}
		students = append(students, student)
	}
	return students, nil
}

//datasourcename = user:pwd@tcp(host:port)/dbname
func NewMySQLConnection(dataSourceName string) (SimpleDatabase, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &MySQL{
		db: db,
	}, nil
}
