package models

import (
	"context"
	"fmt"

	"examples/restApi/config"
)

type Todo struct {
	ID      int
	Title   int
	Content int
}

func GetTodoByID(context context.Context, todoID int) (*Todo, error) {
	var todo *Todo

	config, err := config.FromContext(context)
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`SELECT id, title, content FROM todos WHERE id = ? LIMIT 1;`, todoID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		todo = &Todo{}
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Content)
		return todo, nil
	}

	return nil, fmt.Errorf("Todo not found")
}

func GetTodos(context context.Context) ([]*Todo, error) {
	var todo *Todo

	config, err := config.FromContext(context)
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`SELECT id, title, content FROM todos;`)
	if err != nil {
		return nil, err
	}

	todos := []*Todo{}
	for rows.Next() {
		todo = &Todo{}

		if err = rows.Scan(&todo.ID, &todo.Title, &todo.Content); err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func CreateTodo(context context.Context, todo *Todo) (*Todo, error) {

	config, err := config.FromContext(context)
	if err != nil {
		return nil, err
	}

	_ = config.DB

	return nil, fmt.Errorf("Todo not found")
}
