package security

import (
	"crypto/rand"
	"io"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/common/security"
)

type NoteManager struct {
	logger    common.Logger
	cred      *security.Credential
	cryptoKey []byte
	signKey   []byte
}

var noteManager *NoteManager = nil

func NewNoteManager(logger common.Logger, cred *security.Credential) *NoteManager {
	if noteManager == nil {
		// TODO: исправить  на нормальную генерацию ключей
		cryptoKey := []byte("passphrasewhichneedstobe32bytes!")
		signKey := []byte("passphrasewhichneedstobe32bytes!")
		noteManager = &NoteManager{logger: logger, cred: cred,
			cryptoKey: cryptoKey, signKey: signKey}
	}
	return noteManager
}

func (manager *NoteManager) GetLogin() string {
	return manager.cred.GetLogin()
}
func (manager *NoteManager) CreateNewNotes(plainText string, description string, metadata *clientCommon.Metadata) *Note {
	serializedMetadata := clientCommon.SerializeMetadata(metadata)
	dataForEncrypt := make([]byte, sizeField+len(serializedMetadata)+sizeField+len(plainText))
	actualIndex := 0
	saveDataToStream(dataForEncrypt, &actualIndex, serializedMetadata)
	saveDataToStream(dataForEncrypt, &actualIndex, []byte(plainText))

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
	var temp []byte
	temp = append(temp, info...)
	temp = append(temp, temporyKey...)
	temp = append(temp, cipherText...)
	sign, _ := common.Sign(temp, manager.signKey)

	note := CreatNote(info, temporyKey, cipherText, sign)
	return note
}

func (manager *NoteManager) OpenNote(data []byte) *Note {

	var note Note
	err := note.Deserialize(data)
	if err != nil {
		manager.logger.Warnf("can't deserialize data")
		return nil
	}
	return &note
}
func (manager *NoteManager) GetInfoFromNote(note *Note) ([]byte, *clientCommon.Metadata) {
	decryptText, err := common.DecryptAES(note.cipherText, manager.cryptoKey)
	if err != nil {
		manager.logger.Warnf("problem with create new Note. err:%s", err)
		return []byte{}, nil
	}
	actualIndex := 0
	metadata, _ := loadDataFromStream(decryptText, &actualIndex)
	met, _ := clientCommon.DeserializeMetadata(metadata)
	plainText, _ := loadDataFromStream(decryptText, &actualIndex)
	return plainText, met
}
