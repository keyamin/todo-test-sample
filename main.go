package main

import "fmt"

// ///////////// entity ///////////////
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ////////////////////////////////////

// /////////// repository /////////////
var _ TodoRepository = (*TodoRepositoryImpl)(nil)

//go:generate moq -fmt goimports -out todo_repository_mock.go . TodoRepository
type TodoRepository interface {
	FindByID(id int) (*Todo, error)
}

type TodoRepositoryImpl struct{}

func (tr *TodoRepositoryImpl) FindByID(id int) (*Todo, error) {
	// Do something

	return &Todo{
		ID:        id,
		Title:     "title",
		Completed: false,
	}, nil
}

// ////////////////////////////////////

// ///////////// usecase //////////////
var _ TodoUsecase = (*TodoUsecaseImpl)(nil)

//go:generate moq -fmt goimports -out todo_usecase_mock.go . TodoUsecase
type TodoUsecase interface {
	One(id int) (*Todo, error)
}

func NewTodoUsecase(todoRepository TodoRepository) TodoUsecase {
	return &TodoUsecaseImpl{
		todoRepository: todoRepository,
	}
}

type TodoUsecaseImpl struct {
	todoRepository TodoRepository
}

func (tu *TodoUsecaseImpl) One(id int) (*Todo, error) {
	return tu.todoRepository.FindByID(id)
}

// ////////////////////////////////////

func main() {
	todoRepository := &TodoRepositoryImpl{}
	todoUsecase := NewTodoUsecase(todoRepository)

	todo, _ := todoUsecase.One(1)

	fmt.Println(todo.ID)
	fmt.Println(todo.Title)
	fmt.Println(todo.Completed)
}
