package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mattn/go-sqlite3"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// subject is empty, return error
	if subject == "" {
		return nil, sqlite3.Error{Code: sqlite3.ErrConstraint}
	}

	// execute insert query
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// execute confirm query
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// todoのidを設定
	todo.ID = id

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var (
		rows *sql.Rows
		err  error
	)

	// prevIDがある場合はreadWithIDを実行し、ない場合はreadを実行
	if prevID != 0 {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		rows, err = s.db.QueryContext(ctx, read, size)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// rowsをscanするためのTODOのスライスを作成
	todos := make([]*model.TODO, 0)
	for rows.Next() {
		var todo model.TODO
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	// todosが空の場合は空のスライスを返す
	if len(todos) == 0 {
		todos = []*model.TODO{}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// id is empty, return ErrNotFound
	if id == 0 {
		return nil, &model.ErrNotFound{}
	}

	// subject empty, return error
	if subject == "" {
		return nil, sqlite3.Error{Code: sqlite3.ErrConstraint}
	}

	// execute update query
	row, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	// row affected is 0, return ErrNotFound
	if affected, _ := row.RowsAffected(); affected == 0 {
		return nil, &model.ErrNotFound{}
	}

	// execute confirm query
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// todoのidを設定
	todo.ID = id

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	// idsが空の場合はnilを返す
	if len(ids) == 0 {
		return nil
	}

	// クエリのプレースホルダーを生成
	placeholders := strings.Repeat(",?", len(ids)-1)

	// クエリを生成
	query := fmt.Sprintf(deleteFmt, placeholders)

	// クエリの引数を生成
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// execute delete query
	rows, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// rows affected is 0, return ErrNotFound
	if affected, _ := rows.RowsAffected(); affected == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
