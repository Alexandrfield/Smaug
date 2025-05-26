package worker

import (
	"context"
	"fmt"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/server/security"
)

type TunnelFromClient interface {
	Registration(login string, password string) error
	Login(login string, password string) ([]byte, error)
	SaveData(description string, dat []byte) error
	GetData(name string) ([][]byte, error)
	GetAllData() ([][]byte, error)
}
type DatabaseVault interface {
	CreateNewUser(login string, password string, signKeyComplicated string, key string) error
	GetUser(login string) (string, string, error)
	AddData(login string, description string, data string) error
	GetData(login string, description string) ([]string, error)
	GetAllData(login string) ([]string, error)
}

type ServerSafe struct {
	logger   common.Logger
	tun      TunnelFromClient
	vault    *security.SafetyVault
	database DatabaseVault
}

func (serv *ServerSafe) Registration(login string, password string) error {
	password = security.HashPassword(password)
	key, signKeyComplicated := serv.vault.CreateKeys()
	err := serv.database.CreateNewUser(login, password, string(signKeyComplicated), string(key))
	if err != nil {
		return fmt.Errorf("Registration not complited. err:%w")
	}
	serv.logger.Infof("User:%s successfully registered.", login)
	return nil
}

func (serv *ServerSafe) Login(login string, password string) error {
	password = security.HashPassword(password)
	passwordInDB, _, err := serv.database.GetUser(login)
	if err != nil {
		return fmt.Errorf("Login not complited. err:%w")
	}
	if passwordInDB != password {
		return fmt.Errorf("Login not complited. err: Password incorrect")
	}
	serv.logger.Infof("User:%s successfully logined.", login)
	return nil
}

func (serv *ServerSafe) checkSign(login string, data []byte) bool {
	_, comlicatedSignKey, err := serv.database.GetUser(login)
	if err != nil {
		serv.logger.Warnf("Problem with check sign. Uder doesent exists. err:%w")
		return false
	}
	var note common.Note
	note.Deserialize(data)
	signKey := common.CalculateSignKey([]byte(comlicatedSignKey), note.TemporaryKey)
	return common.CheckHash(note.GetDataForSign(), note.Sign, signKey)
}

func (serv *ServerSafe) AddData(login string, description string, data string) (string, error) {
	_, comlicatedSignKey, err := serv.database.GetUser(login)
	if err != nil {
		return "", fmt.Errorf("Login not complited. err:%w")
	}
	if passwordInDB != password {
		return "", fmt.Errorf("Login not complited. err: Password incorrect")
	}
	serv.logger.Infof("User:%s successfully logined.", login)
	return comlicatedSignKey, nil
}
func ServerLoop(ctx context.Context, done chan struct{}, logger common.Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Stop ServerLoop")
			close(done)
			return
		default:
			continue
		}
	}
}
