package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alexandrfield/Smaug/internal/common"

	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type KeysCreator interface {
	CreateKeys() (string, string)
}
type DatabaseStorage struct {
	Logger      common.Logger
	db          *sql.DB
	DatabaseDsn string
	creator     KeysCreator
}

var ErrPasswordNotValidForUser = errors.New("for this user password not valids")

func NewMemDatabaseStorage(logger common.Logger, dsn string, creator KeysCreator) *DatabaseStorage {
	memStorage := DatabaseStorage{Logger: logger, DatabaseDsn: dsn, creator: creator}
	return &memStorage
}
func (st *DatabaseStorage) createTable(ctx context.Context) error {
	const queryUsers = `CREATE TABLE if NOT EXISTS Users (id text PRIMARY KEY, 
	login text, password text, signKeyComplicated text, key text)`
	if _, err := st.db.ExecContext(ctx, queryUsers); err != nil {
		return fmt.Errorf("error while trying to create table: %w", err)
	}
	const queryData = `CREATE TABLE if NOT EXISTS Data (id text PRIMARY KEY, 
	user text, info text, data text)`
	if _, err := st.db.ExecContext(ctx, queryData); err != nil {
		return fmt.Errorf("error while trying to create table: %w", err)
	}
	return nil
}
func (st *DatabaseStorage) Close() {
	if st.db != nil {
		st.db.Close()
	}
}

func (st *DatabaseStorage) Start(databaseDsn string) error {
	var err error
	st.db, err = sql.Open("pgx", databaseDsn)
	if err != nil {
		return fmt.Errorf("can not open database. err:%w", err)
	}
	st.Logger.Infof("Connect to db open")
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = st.createTable(ctx)
	if err != nil {
		errClose := st.db.Close()
		if errClose != nil {
			return fmt.Errorf("can not create table err:%w; end close connection to database err:%w",
				err, errClose)
		}
		return fmt.Errorf("can not create table. err:%w", err)
	}
	return nil
}

func (st *DatabaseStorage) isUserLoginExist(login string) bool {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT id FROM Users WHERE login = $1", login)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return false
	}
	return true
}

func (st *DatabaseStorage) CreateNewUser(login string, password string, signKeyComplicated string, key string) (string, error) {
	if !st.isUserLoginExist(login) {
		return "", fmt.Errorf("user with login:%s is already exists", login)
	}

	tx, err := st.db.Begin()
	if err != nil {
		return "", fmt.Errorf("can not create transaction. err:%w", err)
	}
	st.Logger.Debugf("CreateNewUser login:%s;", login)

	query := `INSERT INTO Users (login, password, signKeyComplicated, key) VALUES ($1, $2, $3, $4)`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := tx.ExecContext(ctx, query, login, password, signKeyComplicated, key); err != nil {
		errRol := tx.Rollback()
		if errRol != nil {
			return "", fmt.Errorf("error create new user login:%s. err:%w; and error rollback err:%w",
				login, err, errRol)
		}
		return "", fmt.Errorf("error create new user. login:%s: %w", login, err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("error with commit transactiom CreateNewUser. err:%w", err)
	}
	return signKeyComplicated, nil
}

func (st *DatabaseStorage) LoginUser(login string, password string) error { // return user_id, signKeyComplicated
	row := st.db.QueryRowContext(context.Background(),
		"SELECT id, password, signKeyComplicated FROM Users WHERE login = $1", login)
	var userID int
	var passwd string
	err := row.Scan(&userID, &passwd)
	if err != nil {
		return fmt.Errorf("error scan value from row. err:%w", err)
	}
	if password != passwd {
		return ErrPasswordNotValidForUser
	}
	return nil
}
