package controller

/*
controller パッケージは，入力に対するアダプターです．

ここでは，インプットポートとアウトプットポートを組み立てて，
インプットポートを実行します．
ユースケースレイヤからの戻り値を受け取って出力する必要はなく，
純粋にhttpを受け取り，ユースケースを実行します．
*/

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arkuchy/clean-architecture-sample-sample/usecase/port"
)

type User struct {
	OutputFactory func(w http.ResponseWriter) port.UserOutputPort
	// -> presenter.NewUserOutputPort
	InputFactory func(o port.UserOutputPort, u port.UserRepository) port.UserInputPort
	// -> interactor.NewUserInputPort
	RepoFactory func(c *sql.DB) port.UserRepository
	// -> gateway.NewUserRepository
	Conn *sql.DB
}

// GetUserByID は，httpを受け取り，portを組み立てて，inputPort.GetUserByIDを呼び出します．
func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := strings.TrimPrefix(r.URL.Path, "/user/")
	outputPort := u.OutputFactory(w)
	repository := u.RepoFactory(u.Conn)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.GetUserByID(ctx, userID)
}

type PostUserByName struct {
	Name string
}

func (u *User) PostUserByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PostUserByName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	outputPort := u.OutputFactory(w)
	repository := u.RepoFactory(u.Conn)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.PostUserByName(ctx, req.Name)
}

func (u *User) HandleUserRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u.GetUserByID(w, r)
	case "POST":
		u.PostUserByName(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
