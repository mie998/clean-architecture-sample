package presenter

/*
presenter パッケージは，出力に対するアダプターです．

ここでは，アウトプットポートを実装します(interfaceを満たすようにmethodを追加するということ)
*/

import (
	"fmt"
	"net/http"

	"github.com/ari1021/clean-architecture/entity"
	"github.com/ari1021/clean-architecture/usecase/port"
)

type User struct {
	w http.ResponseWriter
}

func NewUser(w http.ResponseWriter) port.UserOutputPort {
	return &User{
		w: w,
	}
}

// usecase.UserOutputPortを実装している
func (u *User) Render(user *entity.User) {
	u.w.WriteHeader(http.StatusOK)
	// httpでentity.User.Nameを出力
	fmt.Fprint(u.w, user.Name)
}

func (u *User) RenderError(err error) {
	u.w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(u.w, err)
}
