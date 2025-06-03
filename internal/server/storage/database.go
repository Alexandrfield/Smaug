package storage

import (
	"context"
	"database/sql"
	"embed"
	b64 "encoding/base64"
	"fmt"
	"time"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseStorage struct {
	logger      common.Logger
	db          *sql.DB
	databaseDsn string
}

func NewMemDatabaseStorage(logger common.Logger, dsn string) *DatabaseStorage {
	memStorage := DatabaseStorage{logger: logger, databaseDsn: dsn}
	err := memStorage.Start()
	if err != nil {
		logger.Errorf("can't create database.err:%s", err)
		return nil
	}
	return &memStorage
}

//go:embed migrations/*.sql
var migrationFolder embed.FS

func (st *DatabaseStorage) Migrate() error {
	d, err := iofs.New(migrationFolder, "migrations")
	if err != nil {
		return fmt.Errorf("problem with migration. err:%w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, st.databaseDsn)

	if err != nil {
		st.logger.Errorf("problem with migrate.NewWithDatabaseInstance. st.databaseDsn:%s; err:%s", st.databaseDsn, err)
	}
	err = m.Down()
	if err != nil {
		st.logger.Errorf("problem m.Up. err:%s", err)
	}

	err = m.Up()
	if err != nil {
		st.logger.Errorf("problem m.Up. err:%s", err)
	}
	return nil
}
func (st *DatabaseStorage) Close() {
	if st.db != nil {
		st.db.Close()
	}
}

func (st *DatabaseStorage) Start() error {
	var err error
	st.db, err = sql.Open("pgx", st.databaseDsn)
	if err != nil {
		return fmt.Errorf("can not open database. err:%w", err)
	}
	st.logger.Infof("Connect to db open")
	err = st.Migrate()
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
	st.logger.Debugf("isUserLoginExist %s check:%t", login, err == nil)
	return err == nil
}

func (st *DatabaseStorage) CreateNewUser(login string, password []byte, signKeyComplicated []byte, key []byte) error {
	if st.isUserLoginExist(login) {
		return fmt.Errorf("user with login:%s is already exists", login)
	}

	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("can not create transaction. err:%w", err)
	}
	st.logger.Debugf("CreateNewUser login:%s;", login)
	st.logger.Debugf(">>>> password:%s; signKeyComplicated:%s; key:%s;", password, signKeyComplicated, key) //TODO: Remove
	query := `INSERT INTO Users (login, password, signKeyComplicated, key) VALUES ($1, $2, $3, $4)`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := tx.ExecContext(ctx, query, login, b64.StdEncoding.EncodeToString(password),
		b64.StdEncoding.EncodeToString(signKeyComplicated), b64.StdEncoding.EncodeToString(key)); err != nil {
		errRol := tx.Rollback()
		if errRol != nil {
			return fmt.Errorf("error create new user login:%s. err:%w; and error rollback err:%w",
				login, err, errRol)
		}
		return fmt.Errorf("error create new user. login:%s: %w", login, err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error with commit transactiom CreateNewUser. err:%w", err)
	}
	return nil
}

func (st *DatabaseStorage) GetUserPassword(login string) ([]byte, error) {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT password FROM Users WHERE login = $1", login)
	var passwd string
	err := row.Scan(&passwd)
	if err != nil {
		return []byte{}, fmt.Errorf("error scan value from row. err:%w", err)
	}
	pas, _ := b64.StdEncoding.DecodeString(passwd)
	return pas, nil
}

func (st *DatabaseStorage) GetUserKey(login string) ([]byte, []byte, error) { // return signKeyComplicated, cryptoKey
	row := st.db.QueryRowContext(context.Background(),
		"SELECT signKeyComplicated, key FROM Users WHERE login = $1", login)
	var signKeyComplicated string
	var cryptoKey string
	err := row.Scan(&signKeyComplicated, &cryptoKey)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("error scan value from row. err:%w", err)
	}
	sigKey, _ := b64.StdEncoding.DecodeString(signKeyComplicated)
	key, _ := b64.StdEncoding.DecodeString(cryptoKey)
	return sigKey, key, nil
}

func (st *DatabaseStorage) AddData(login string, description string, data []byte) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("can not create transaction AddData. err:%w", err)
	}
	query := `INSERT INTO EncryptedData (login, info, data ) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := tx.ExecContext(ctx, query, login, description, b64.StdEncoding.EncodeToString(data)); err != nil {
		return fmt.Errorf("tx, error while trying to ord. err: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error with commit transaction AddData. err:%w", err)
	}
	return nil
}

func (st *DatabaseStorage) GetData(login string, description string) ([][]byte, error) {
	var res [][]byte
	rows, err := st.db.QueryContext(context.Background(),
		"SELECT data FROM EncryptedData WHERE login=$1 and info=$2", login, description)
	if err != nil {
		return res, fmt.Errorf("problem GetData. err:%w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var dat string
		err := rows.Scan(&dat)
		if err != nil {
			st.logger.Warnf("error scan value from row. err:%s", err)
		}
		d, _ := b64.StdEncoding.DecodeString(dat)
		res = append(res, d)
	}
	err = rows.Err()
	if err != nil {
		st.logger.Warnf("error rows. err:%s", err)
	}
	return res, nil
}

func (st *DatabaseStorage) GetAllData(login string) ([][]byte, error) {
	var res [][]byte
	rows, err := st.db.QueryContext(context.Background(),
		"SELECT data FROM EncryptedData WHERE login=$1", login)
	if err != nil {
		return res, fmt.Errorf("problem GetData. err:%w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var dat string
		err := rows.Scan(&dat)
		if err != nil {
			st.logger.Warnf("error scan value from row. err:%s", err)
		}
		d, _ := b64.StdEncoding.DecodeString(dat)
		res = append(res, d)
	}
	err = rows.Err()
	if err != nil {
		st.logger.Warnf("error rows. err:%s", err)
	}
	return res, nil
}
