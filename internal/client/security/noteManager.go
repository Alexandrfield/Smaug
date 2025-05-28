package security

import (
	"crypto/rand"
	"crypto/sha256"
	"io"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
	"github.com/Alexandrfield/Smaug/internal/common"
)

type NoteManager struct {
	logger    common.Logger
	login     string
	cryptoKey []byte
	signKey   []byte
}

var noteManager *NoteManager = nil

func NewNoteManager(logger common.Logger) *NoteManager {
	if noteManager == nil {

		noteManager = &NoteManager{logger: logger}
	}
	return noteManager
}

func (manager *NoteManager) InitParametrs(login []byte, password []byte) {
	manager.login = string(login)
	manager.saveSignKey(password)
	manager.saveCryptoKey(login, password)
}
func (manager *NoteManager) saveSignKey(password []byte) {
	manager.signKey = common.ComplicatedPasswordForPrepareSign(password)
	manager.logger.Debugf("manager.signKey:%v", manager.signKey)
}
func (manager *NoteManager) saveCryptoKey(login []byte, password []byte) {
	h := sha256.New()
	h.Write([]byte{0x02, 0xfa, 0xa4, 0x55})
	h.Write(login)
	h.Write(password)
	manager.cryptoKey = h.Sum(nil)
}
func (manager *NoteManager) GetLogin() string {
	return manager.login
}
func (manager *NoteManager) CreateNewNotes(plainText string, description string, metadata *clientCommon.Metadata) *common.Note {
	serializedMetadata := clientCommon.SerializeMetadata(metadata)
	dataForEncrypt := make([]byte, common.SizeField+len(serializedMetadata)+common.SizeField+len(plainText))
	actualIndex := 0
	common.SaveDataToStream(dataForEncrypt, &actualIndex, serializedMetadata)
	common.SaveDataToStream(dataForEncrypt, &actualIndex, []byte(plainText))

	cipherText, err := common.EncryptAES(dataForEncrypt, manager.cryptoKey)
	if err != nil {
		manager.logger.Warnf("problem with create new Note. err:%s", err)
		return nil
	}
	io.ReadFull(rand.Reader, dataForEncrypt)

	temporyKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, temporyKey); err != nil {
		manager.logger.Warnf("Can't create temporyKey: %s", err)
		return nil
	}
	info := []byte(description)
	note := common.CreatNote(info, temporyKey, cipherText)

	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(manager.signKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	return note
}

func (manager *NoteManager) CreateServiceNotes(description string) *common.Note {
	RandDataAsCrypt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, RandDataAsCrypt); err != nil {
		manager.logger.Warnf("Can't create tempData: %s", err)
		return nil
	}

	temporyKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, temporyKey); err != nil {
		manager.logger.Warnf("Can't create temporyKey: %s", err)
		return nil
	}
	info := []byte(description)
	note := common.CreatNote(info, temporyKey, RandDataAsCrypt)

	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(manager.signKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	return note
}

func (manager *NoteManager) OpenNote(data []byte) *common.Note {

	var note common.Note
	err := note.Deserialize(data)
	if err != nil {
		manager.logger.Warnf("can't deserialize data")
		return nil
	}
	return &note
}
func (manager *NoteManager) GetInfoFromNote(note *common.Note) ([]byte, *clientCommon.Metadata) {
	decryptText, err := common.DecryptAES(note.CipherText, manager.cryptoKey)
	if err != nil {
		manager.logger.Warnf("problem with create new Note. err:%s", err)
		return []byte{}, nil
	}
	actualIndex := 0
	metadata, _ := common.LoadDataFromStream(decryptText, &actualIndex)
	met, _ := clientCommon.DeserializeMetadata(metadata)
	plainText, _ := common.LoadDataFromStream(decryptText, &actualIndex)
	return plainText, met
}
