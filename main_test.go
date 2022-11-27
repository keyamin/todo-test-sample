package main

import (
	_ "embed"
	"errors"
	"reflect"
	"testing"
)

//go:embed testdata/sql/todo/find_by_id.sql
var findByIDInitSQL string

func TestMain(m *testing.M) {
	// ここでDBの初期化処理を実行し、deferでクローズ処理を登録

	m.Run()
}

func TestTodoRepository_FindByID(t *testing.T) {
	t.Parallel()

	type args struct {
		id int
	}
	tests := map[string]struct {
		args    args
		want    *Todo
		wantErr bool
	}{
		"ok":        {args{1}, DefaultTodo(), false},
		"not found": {args{2}, nil, true},
	}
	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// ここでトランザクションを張り、(*testing.T).Cleanup()でロールバック処理を登録
			// ここでヘルパー関数を使ってfindByIDInitSQLを実行

			tr := &TodoRepositoryImpl{}

			got, err := tr.FindByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TodoRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TodoRepository_FindByIDの引数と戻り値
type TodoRepository_FindByID_Expectation struct {
	Input  TodoRepository_FindByID_Input
	Output TodoRepository_FindByID_Output
}
type TodoRepository_FindByID_Input struct {
	ID int
}
type TodoRepository_FindByID_Output struct {
	Todo *Todo
	Err  error
}

// デフォルト値のTodo構造体を生成するヘルパーと各種オプション関数
type todoOptFunc func(*Todo)

func WithID(id int) todoOptFunc                { return func(t *Todo) { t.ID = id } }
func WithTitle(title string) todoOptFunc       { return func(t *Todo) { t.Title = title } }
func WithCompleted(completed bool) todoOptFunc { return func(t *Todo) { t.Completed = completed } }

func DefaultTodo(opts ...todoOptFunc) *Todo {
	t := &Todo{
		ID:        1,
		Title:     "title",
		Completed: false,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func TestTodoUsecaseImpl_One(t *testing.T) {
	t.Parallel()

	// 冗長な記述になるのを避けるため型エイリアス作成
	type tfe = TodoRepository_FindByID_Expectation
	type tfi = TodoRepository_FindByID_Input
	type tfo = TodoRepository_FindByID_Output

	type args struct {
		id int
	}
	tests := map[string]struct {
		args              args
		TodoRepo_FindByID tfe
		want              *Todo
		wantErr           bool
	}{
		"ok":        {args{1}, tfe{tfi{1}, tfo{DefaultTodo(), nil}}, DefaultTodo(), false},
		"not found": {args{2}, tfe{tfi{2}, tfo{nil, errors.New("test")}}, nil, true},
		// ↓こう書きたくない！
		// "bad": {
		// 	args: args{
		// 		id: 1,
		// 	},
		// 	TodoRepo_FindByID: TodoRepository_FindByID_Expectation{
		// 		Input: TodoRepository_FindByID_Input{
		// 			ID: 1,
		// 		},
		// 		Output: TodoRepository_FindByID_Output{
		// 			Todo: &Todo{
		// 				ID:        1,
		// 				Title:     "title",
		// 				Completed: false,
		// 			},
		// 		},
		// 	},
		// 	want: &Todo{
		// 		ID:        1,
		// 		Title:     "title",
		// 		Completed: false,
		// 	},
		// 	wantErr: false,
		// },
	}
	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// モック関数の登録
			tr := &TodoRepositoryMock{
				FindByIDFunc: func(_ int) (*Todo, error) {
					return tt.TodoRepo_FindByID.Output.Todo, tt.TodoRepo_FindByID.Output.Err
				},
			}
			tu := NewTodoUsecase(tr)

			got, err := tu.One(tt.args.id)

			// モックが想定通り呼び出されたかどうかの確認
			calls := tr.FindByIDCalls()
			if len(calls) != 1 {
				t.Fatalf("TodoRepository.FindByID() calls = %v, want 1", len(calls))
			}
			if calls[0].ID != tt.TodoRepo_FindByID.Input.ID {
				t.Fatalf("TodoRepository.FindByID() calls[0].ID = %v, want %v", calls[0].ID, tt.TodoRepo_FindByID.Input.ID)
			}

			// 戻り値の確認
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoUsecaseImpl.One() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TodoUsecaseImpl.One() = %v, want %v", got, tt.want)
			}
		})
	}
}
