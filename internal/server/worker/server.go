package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/server/security"
	"github.com/Alexandrfield/Smaug/internal/server/storage"
)

type DatabaseVault interface {
	CreateNewUser(login string, password []byte, signKeyComplicated []byte, key []byte) error
	GetUserPassword(login string) ([]byte, error)
	GetUserKey(login string) ([]byte, []byte, error)
	AddData(login string, description string, data []byte) error
	GetData(login string, description string) ([][]byte, error)
	GetAllData(login string) ([][]byte, error)
}

var ErrUserIsNotExists = errors.New("user isn't exists")

func NewServerSafe(logger common.Logger, config Config) *ServerSafe {
	vault := security.NewSafetyVault(config.MasterKey, logger)
	db := storage.NewMemDatabaseStorage(logger, config.DatabasDSN)
	if db == nil {
		logger.Fatalf("can't createServerSafe. problem with databse")
	}
	s := ServerSafe{logger: logger, vault: vault, database: db}
	return &s
}

type ServerSafe struct {
	logger   common.Logger
	vault    *security.SafetyVault
	database DatabaseVault
}

func (serv *ServerSafe) Registration(login string, password []byte) error {
	serv.logger.Debugf("User:%s try register.", login)
	passwordForDB := security.HashPassword(password)
	serv.logger.Debugf("save passwordInDB:%v", password) //TODO:remove
	key, signKeyComplicated := serv.vault.CreateKeys([]byte(password))
	err := serv.database.CreateNewUser(login, passwordForDB, signKeyComplicated, key)
	if err != nil {
		err := fmt.Errorf("registration not complited. err:%w", err)
		serv.logger.Warnf("err:%s", err)
		return err
	}
	serv.logger.Infof("User:%s successfully registered.", login)
	return nil
}

func (serv *ServerSafe) Login(login string, password []byte) error {
	serv.logger.Debugf("User:%s try login.", login)
	password = security.HashPassword(password)
	passwordInDB, err := serv.database.GetUserPassword(login)
	if err != nil {
		err := fmt.Errorf("login not complited. err:%w", err)
		serv.logger.Warnf("err:%s", err)
		return err
	}
	if !bytes.Equal(passwordInDB, password) {
		serv.logger.Debugf("passwordInDB:%v<->password:%v", passwordInDB, password) //TODO:Remove
		err := fmt.Errorf("login not complited. err: Password incorrect")
		serv.logger.Warnf("err:%s", err)
		return err
	}
	serv.logger.Infof("User:%s successfully logined.", login)
	return nil
}

func (serv *ServerSafe) checkSign(note *common.Note, comlicatedSignKey []byte) bool {
	signKey := common.CalculateSignKey([]byte(comlicatedSignKey), note.TemporaryKey)
	return common.CheckHash(note.GetDataForSign(), note.Sign, signKey)
}

func getNoteFormByte(data string) *common.Note {
	var note common.Note
	note.Deserialize([]byte(data))
	return &note
}
func (serv *ServerSafe) AddData(login string, data string) error {
	serv.logger.Debugf("User:%s try addData.", login)
	comlicatedSignKey, cryptoKey, err := serv.database.GetUserKey(login)
	if err != nil {
		err := fmt.Errorf("problem with add data. err:%w", ErrUserIsNotExists)
		serv.logger.Warnf("err:%s", err)
		return err
	}
	note := getNoteFormByte(data)
	if !serv.checkSign(note, comlicatedSignKey) {
		err := fmt.Errorf("problem with add data. err:incorrect signature")
		serv.logger.Warnf("err:%s", err)
		return err
	}

	encryptedData := serv.vault.EncryptData([]byte(data), cryptoKey)
	err = serv.database.AddData(login, string(note.Info), encryptedData)
	if err != nil {
		serv.logger.Warnf("can't add data. err:%s", err)
	}
	return nil
}

func (serv *ServerSafe) GetData(login string, data string) ([]string, error) {
	serv.logger.Debugf("User:%s try getData.", login)
	var res []string
	comlicatedSignKey, cryptoKey, err := serv.database.GetUserKey(login)
	if err != nil {
		err := fmt.Errorf("problem with add data. err:%w", ErrUserIsNotExists)
		serv.logger.Warnf("err:%s", err)
		return res, err
	}
	note := getNoteFormByte(data)
	if !serv.checkSign(note, comlicatedSignKey) {
		err := fmt.Errorf("problem with add data. err:incorrect signature")
		serv.logger.Warnf("err:%s", err)
		return res, err
	}
	datas, err := serv.database.GetData(login, string(note.Info))
	if err != nil {
		err := fmt.Errorf("problem with get data. err:%w", err)
		serv.logger.Warnf("err:%s", err)
		return res, err
	}
	for _, val := range datas {
		decryptedData := serv.vault.DecryptData([]byte(val), cryptoKey)
		res = append(res, string(decryptedData))
	}
	return res, nil
}

func (serv *ServerSafe) GetAllData(login string, data string) ([]string, error) {
	serv.logger.Debugf("User:%s try GetallData.", login)
	var res []string
	comlicatedSignKey, cryptoKey, err := serv.database.GetUserKey(login)
	if err != nil {
		err := fmt.Errorf("problem with add data. err:%w", ErrUserIsNotExists)
		serv.logger.Warnf("err:%s", err)
		return res, err
	}
	note := getNoteFormByte(data)
	if !serv.checkSign(note, comlicatedSignKey) {
		err := fmt.Errorf("problem with add data. err:incorrect signature")
		serv.logger.Warnf("err:%s", err)
		return res, err
	}
	datas, err := serv.database.GetAllData(login)
	if err != nil {
		err := fmt.Errorf("problem with get data. err:%w", err)
		serv.logger.Warnf("err:%s", err)
		return res, err
	}
	for _, val := range datas {
		decryptedData := serv.vault.DecryptData([]byte(val), cryptoKey)
		res = append(res, string(decryptedData))
	}
	return res, nil
}

func ServerLoop(ctx context.Context, done chan struct{}, config Config, logger common.Logger) {
	serv := NewServerSafe(logger, config)
	go StartGRPCServer(logger, serv, config.Port)
	<-ctx.Done()
	logger.Infof("Stop ServerLoop")
	close(done)
	return
}
