package worker

import (
	"fmt"
	"os"
	"time"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
	clientSecurity "github.com/Alexandrfield/Smaug/internal/client/security"
	"github.com/Alexandrfield/Smaug/internal/common"
)

//go:generate mockgen -source=client.go -destination=mock/client.go
type TunnelToServer interface {
	Registration(login string, password string) error
	Login(login string, password []byte) error
	SaveData(login string, dat []byte) error
	GetData(login string, dat []byte) ([][]byte, error)
	GetAllData(login string, dat []byte) ([][]byte, error)
}

type SmaugClient struct {
	logger          common.Logger
	secyrityManager *clientSecurity.NoteManager
	networkClient   TunnelToServer
	localCache      map[string][]*common.Note
}

func NewSmaugClient(logger common.Logger, config *Config) *SmaugClient {
	secyrityManager := clientSecurity.NewNoteManager(logger)
	cl := SmaugClient{logger: logger, secyrityManager: secyrityManager, networkClient: &GRPCClient{serverAddr: config.ServerAddr}}
	cl.localCache = make(map[string][]*common.Note)
	return &cl
}
func (client *SmaugClient) Registration(login string, password []byte) error {
	passwordForServer := clientSecurity.HashInputPassword([]byte(password))
	err := client.networkClient.Registration(login, string(passwordForServer))
	if err != nil {
		return fmt.Errorf("problem with registration. can't continue. err:%w", err)
	}
	err = client.login(login, passwordForServer)
	if err != nil {
		return fmt.Errorf("problem with login after registration. err:%w", err)
	}
	return nil
}
func (client *SmaugClient) login(login string, password []byte) error {
	err := client.networkClient.Login(login, password)
	if err != nil {
		return fmt.Errorf("problem with login. err:%w", err)
	}
	client.secyrityManager.InitParametrs([]byte(login), []byte(password))
	return nil
}

func (client *SmaugClient) Login(login string, password []byte) error {
	passwordForServer := clientSecurity.HashInputPassword(password)
	return client.login(login, passwordForServer)
}
func (client *SmaugClient) saveData(metadata *clientCommon.Metadata, description string, data []byte) error {
	newNote := client.secyrityManager.CreateNewNotes(data, description, metadata)
	err := client.networkClient.SaveData(client.secyrityManager.GetLogin(), newNote.Serialize())
	if err != nil {
		return fmt.Errorf("saveData err:%w", err)
	}
	client.localCache[description] = append(client.localCache[description], newNote)
	return nil
}

func (client *SmaugClient) SaveText(description string, data []byte) error {
	metadata := clientCommon.NewMetadata("text", client.secyrityManager.GetLogin(), time.Now())
	err := client.saveData(metadata, description, data)
	if err != nil {
		client.logger.Errorf("problem with save data. err:%s", err)
		return fmt.Errorf("can't processing this note:%s", description)
	}
	return nil
}
func (client *SmaugClient) SaveFile(description string, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		client.logger.Errorf("problem with read file:%s; err:%s", filename, err)
		return fmt.Errorf("can't processing file:%s", filename)
	}
	metadata := clientCommon.NewMetadata("file", client.secyrityManager.GetLogin(), time.Now())
	metadata.Filename = filename
	err = client.saveData(metadata, description, data)
	if err != nil {
		return fmt.Errorf("SaveData err:%w", err)
	}
	return nil
}

func (client *SmaugClient) updateLocalCash(rawByteData [][]byte) {
	client.logger.Infof("updateLocalCash")
	for _, val := range rawByteData {
		note := client.secyrityManager.OpenNote(val)
		client.localCache[string(note.Info)] = append(client.localCache[string(note.Info)], note)
		client.logger.Debugf("update %s note", string(note.Info))
	}
}

func (client *SmaugClient) prepareOpenData(listDescriptions []string) []*OpenData {
	var res []*OpenData
	for _, v := range listDescriptions {
		notes := client.localCache[v]
		for _, note := range notes {
			client.logger.Debugf("prepare %s note", string(note.Info))
			plaintext, metadata := client.secyrityManager.GetInfoFromNote(note)
			if len(plaintext) == 0 || metadata == nil {
				continue
			}
			t := &OpenData{metadata: metadata, data: plaintext, description: v}
			res = append(res, t)
		}
	}
	return res
}
func (client *SmaugClient) GetData(description string) []*OpenData {
	client.logger.Debugf("try GetData description:%s", description)
	tempNote := client.secyrityManager.CreateServiceNotes(description)
	data, err := client.networkClient.GetData(client.secyrityManager.GetLogin(), tempNote.Serialize())
	if err != nil {
		client.logger.Warnf("Problem Get data. err:%s. Used local cash.", err)
	} else {
		client.updateLocalCash(data)
	}
	return client.prepareOpenData([]string{description})
}
func (client *SmaugClient) GetAllData() []*OpenData {
	client.logger.Debugf("try GetAllData")
	tempNote := client.secyrityManager.CreateServiceNotes("get all data")
	data, err := client.networkClient.GetAllData(client.secyrityManager.GetLogin(), tempNote.Serialize())
	if err != nil {
		client.logger.Warnf("Problem Get all data. err:%s. Used local cash.", err)
	} else {
		client.updateLocalCash(data)
	}
	descriptions := make([]string, 0, len(client.localCache))
	for k := range client.localCache {
		descriptions = append(descriptions, k)
	}
	return client.prepareOpenData(descriptions)
}
