package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"todo"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func (r *TodoItemPostgres) Update(userId, itemId int, item todo.UpdateItem) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if item.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *item.Title)
		argId++
	}

	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *item.Description)
		argId++
	}

	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *item.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		"UPDATE %s ti "+
			"SET %s "+
			"FROM %s li, %s ul "+
			"WHERE ti.id = li.item_id AND ul.list_id = li.list_id AND ul.user_id = $%d AND ti.id = $%d",
		todoItemsTable,
		setQuery,
		listsItemsTable,
		usersListTable,
		argId,
		argId+1,
	)

	args = append(args, userId, itemId)

	logrus.Infof("update query: %s", query)
	logrus.Infof("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	deleteItem := fmt.Sprintf("DELETE "+
		"FROM %s ti "+
		"USING %s li, %s ul "+
		"WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ti.id = $1 AND ul.user_id = $2",
		todoItemsTable,
		listsItemsTable,
		usersListTable,
	)

	_, err = tx.Exec(deleteItem, itemId, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, list todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err = row.Scan(&itemId); err != nil {
		if rollbackErr := tx.Rollback(); err != nil {
			return 0, rollbackErr
		}

		return 0, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	if _, err = tx.Exec(createUserListQuery, listId, itemId); err != nil {
		if rollbackErr := tx.Rollback(); err != nil {
			return 0, rollbackErr
		}

		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(listId, userId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.done "+
			"FROM %s ti "+
			"INNER JOIN %s li ON ti.id = li.item_id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE li.list_id = $1 AND ul.user_id = $2",
		todoItemsTable,
		listsItemsTable,
		usersListTable,
	)

	err := r.db.Select(&items, query, listId, userId)
	return items, err
}

func (r *TodoItemPostgres) GetById(userId int, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.done "+
			"FROM %s ti "+
			"INNER JOIN %s li ON ti.id = li.item_id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE ti.id = $1 AND ul.user_id = $2",
		todoItemsTable,
		listsItemsTable,
		usersListTable,
	)

	err := r.db.Get(&item, query, itemId, userId)
	return item, err
}
