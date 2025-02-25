package gateway

/*
gateway パッケージは，DB操作に対するアダプターです．
*/

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"

	"github.com/arkuchy/clean-architecture-sample-sample/entity"
	"github.com/arkuchy/clean-architecture-sample-sample/usecase/port"
)

type UserRepository struct {
	conn *sql.DB
}

// NewUserRepository はUserRepositoryを返します．
func NewUserRepository(conn *sql.DB) port.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

// GetUserByID はDBからデータを取得します．
func (u *UserRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	conn := u.GetDBConn()
	row := conn.QueryRowContext(ctx, "SELECT * FROM `user` WHERE id=?", userID)
	user := entity.User{}
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user Not Found. UserID = %s", userID)
		}
		log.Println(err)
		return nil, errors.New("internal Server Error. adapter/gateway/GetUserByID")
	}
	return &user, nil
}

func (u *UserRepository) PostUserByName(ctx context.Context, userName string) (*entity.User, error) {
	conn := u.GetDBConn()
	userId := uuid.New().String()
	_, err := conn.ExecContext(ctx, "INSERT INTO `user` (`id`,`name`) VALUES (?, ?)", userId, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("couldn't create user: %s", userName)
		}
		log.Println(err)
		return nil, errors.New("internal Server Error. adapter/gateway/PostUserByName")
	}
	user := entity.User{}
	user.ID = userId
	user.Name = userName
	return &user, nil
}

// GetDBConn はconnectionを取得します．
func (u *UserRepository) GetDBConn() *sql.DB {
	return u.conn
}
