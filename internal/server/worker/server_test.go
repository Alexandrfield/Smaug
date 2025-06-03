package worker

import (
	"errors"
	"testing"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/server/security"
	mock "github.com/Alexandrfield/Smaug/internal/server/worker/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckSignError(t *testing.T) {
	comlicatedSignKey1 := []byte{0x00, 0xff, 0xaa, 0xee}
	comlicatedSignKey2 := []byte{0x00, 0xff, 0xaa}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey1, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	res := checkSign(note, comlicatedSignKey2)
	if res {
		t.Errorf("error check sign. expectedd error")
	}
}

func TestGetNoteFormByte(t *testing.T) {
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	serNote := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x01, 0xff, 0xab, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, 0xff, 0xab, 0xff, 0xff, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	note := getNoteFormByte(string(serNote))
	assert.ElementsMatch(t, note.Info, info)
	assert.ElementsMatch(t, note.CipherText, cipherText)
	assert.ElementsMatch(t, note.TemporaryKey, temporyKey)
}

func TestRegistration(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().CreateNewUser(login, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := serV.Registration(login, password)
	require.NoError(t, err)
}

var ErrTest = errors.New("testError")

func TestRegistrationError(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().CreateNewUser(login, gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrTest)

	err := serV.Registration(login, password)
	if !errors.Is(err, ErrTest) {
		t.Errorf("Unexpected error. ecpected:%s; actual:%s", ErrTest, err)
	}
}

func TestLogin(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	databasePassword := []byte{
		0xbf, 0xfe, 0x04, 0x1b, 0xe9, 0x39, 0xec, 0x36, 0x0d, 0xf5, 0x1c, 0x7b, 0x2d, 0x94, 0x4d, 0x28,
		0x5a, 0x08, 0x80, 0xc2, 0x82, 0x7c, 0xa5, 0xfe, 0x79, 0xfd, 0x03, 0xd4, 0x55, 0xad, 0xbe, 0x2c,
	}
	mockVault.EXPECT().GetUserPassword(login).Return(databasePassword, nil)

	err := serV.Login(login, password)
	require.NoError(t, err)
}
func TestLoginIncorrectPassword(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	databasePassword := []byte{0xbf}
	mockVault.EXPECT().GetUserPassword(login).Return(databasePassword, nil)

	err := serV.Login(login, password)
	if err == nil {
		t.Errorf("shoud be problem (incorrect password)")
	}
}
func TestLoginErr(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	databasePassword := []byte{
		0xbf, 0xfe, 0x04, 0x1b, 0xe9, 0x39, 0xec, 0x36, 0x0d, 0xf5, 0x1c, 0x7b, 0x2d, 0x94, 0x4d, 0x28,
		0x5a, 0x08, 0x80, 0xc2, 0x82, 0x7c, 0xa5, 0xfe, 0x79, 0xfd, 0x03, 0xd4, 0x55, 0xad, 0xbe, 0x2c,
	}
	mockVault.EXPECT().GetUserPassword(login).Return(databasePassword, ErrTest)

	err := serV.Login(login, password)
	if !errors.Is(err, ErrTest) {
		t.Errorf("Unexpected error. ecpected:%s; actual:%s", ErrTest, err)
	}
}

func TestCheckSign(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	res := checkSign(note, comlicatedSignKey)
	if !res {
		t.Errorf("error check sign")
	}
}

func TestAddData(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)
	mockVault.EXPECT().AddData(login, string(info), gomock.Any()).Return(nil)

	err := serV.AddData(login, string(note.Serialize()))
	require.NoError(t, err)
}
func TestAddDataErr(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, ErrTest)

	err := serV.AddData(login, string(note.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
func TestAddDataErr2(t *testing.T) {
	comlicatedSignKey1 := []byte{0x00, 0xff, 0xaa, 0xee}
	comlicatedSignKey2 := []byte{0xff, 0xff, 0xff, 0xff}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey1, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey2, cryptoKey, nil)

	err := serV.AddData(login, string(note.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
func TestAddDataErr3(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)
	mockVault.EXPECT().AddData(login, string(info), gomock.Any()).Return(ErrTest)
	err := serV.AddData(login, string(note.Serialize()))
	require.NoError(t, err)
}

func createServiceNotes(description []byte, signKey []byte, temporyKey []byte) *common.Note {
	RandDataAsCrypt := []byte{
		0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff,
		0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff,
	}
	info := description
	note := common.CreatNote(info, temporyKey, RandDataAsCrypt)

	temp := note.GetDataForSign()
	signKeyN := common.CalculateSignKey(signKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKeyN)
	return note
}

func TestAGetData(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)
	ret := [][]byte{note.Serialize()}
	mockVault.EXPECT().GetData(login, string(info)).Return(ret, nil)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetData(login, string(serviceNote.Serialize()))
	require.NoError(t, err)
}
func TestAGetDataErr(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, ErrTest)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetData(login, string(serviceNote.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
func TestAGetDataErr2Sign(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	comlicatedSignKey2 := []byte{0xff, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey2, cryptoKey, ErrTest)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetData(login, string(serviceNote.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}

func TestAGetDataErr3(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)
	ret := [][]byte{note.Serialize()}
	mockVault.EXPECT().GetData(login, string(info)).Return(ret, ErrTest)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetData(login, string(serviceNote.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}

func TestGetAllData(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)
	ret := [][]byte{note.Serialize()}
	mockVault.EXPECT().GetAllData(login).Return(ret, nil)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetAllData(login, string(serviceNote.Serialize()))
	require.NoError(t, err)
}

func TestGetAllDataErr(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, ErrTest)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetAllData(login, string(serviceNote.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}

func TestGetAllDataErrSign(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	comlicatedSignKey2 := []byte{0xff, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)

	serviceNote := createServiceNotes(info, comlicatedSignKey2, temporyKey)
	_, err := serV.GetAllData(login, string(serviceNote.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
func TestGetAllDataErr2(t *testing.T) {
	comlicatedSignKey := []byte{0x00, 0xff, 0xaa, 0xee}
	info := []byte{0x01, 0xff, 0xab}
	temporyKey := []byte{0x01, 0xff, 0xab, 0xff, 0xff}
	cryptoKey := []byte{0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0x55, 0xff, 0xab, 0xff, 0xaa, 0xff}
	cipherText := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	note := common.CreatNote(info, temporyKey, cipherText)
	temp := note.GetDataForSign()
	signKey := common.CalculateSignKey(comlicatedSignKey, temporyKey)
	note.Sign, _ = common.Sign(temp, signKey)
	login := "testUser"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockDatabaseVault(ctrl)
	masterKey := []byte{0x00, 0xff, 0xe3, 0x54}
	vault := security.NewSafetyVault(masterKey, &common.FakeLogger{})
	serV := ServerSafe{logger: &common.FakeLogger{}, database: mockVault, vault: vault}

	mockVault.EXPECT().GetUserKey(login).Return(comlicatedSignKey, cryptoKey, nil)
	ret := [][]byte{note.Serialize()}
	mockVault.EXPECT().GetAllData(login).Return(ret, ErrTest)

	serviceNote := createServiceNotes(info, comlicatedSignKey, temporyKey)
	_, err := serV.GetAllData(login, string(serviceNote.Serialize()))
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
