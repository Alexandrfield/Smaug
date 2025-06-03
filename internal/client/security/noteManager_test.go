package security

import (
	"testing"
	"time"

	clientCom "github.com/Alexandrfield/Smaug/internal/client/common"
	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestNewNoteManager(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	if p == nil {
		t.Errorf("ptr to NoteManager nil")
	}
	globalNoteManager = nil
}

func TestInitParametrs(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	assert.Equal(t, "testLogin", p.login)
	assert.ElementsMatch(t, p.signKey, common.ComplicatedPasswordForPrepareSign(password))
	assert.ElementsMatch(t, p.cryptoKey, generateCryptoKey(login, password))
}

func TestGetLogin(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	actual := p.GetLogin()
	assert.Equal(t, "testLogin", actual)
}

func TestCreateNewNotes(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	plainText := "test Plain Text"
	description := "test description"
	metadata := clientCom.NewMetadata("testType", "testUser", time.Now())
	note := p.CreateNewNotes([]byte(plainText), description, metadata)
	if note == nil {
		t.Errorf("expected not nil note")
	}
}
func TestCreateNewNotesWithoutParametrs(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	plainText := "test Plain Text"
	description := "test description"
	metadata := clientCom.NewMetadata("testType", "testUser", time.Now())
	note := p.CreateNewNotes([]byte(plainText), description, metadata)
	if note != nil {
		t.Errorf("expected nil note (not set parametrs)")
	}
}

func TestCreateServiceNotes(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	description := "test description"
	note := p.CreateServiceNotes(description)
	if note == nil {
		t.Errorf("expected not nil note")
	}
}

func TestCreateServiceNotesWithoutParametrs(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	description := "test description"
	note := p.CreateServiceNotes(description)
	if note == nil {
		t.Errorf("expected not nil note (not set parametrs)")
	}
}

func TestOpenNote(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	plainText := "test Plain Text"
	description := "test description"
	metadata := clientCom.NewMetadata("testType", "testUser", time.Now())
	note := p.CreateNewNotes([]byte(plainText), description, metadata)
	serByteStream := note.Serialize()
	newNote := p.OpenNote(serByteStream)
	if note == nil {
		t.Errorf("expected not nil note")
	}
	if newNote.Info == nil {
		t.Errorf("expected not nil info")
		return
	}
	if newNote.Sign == nil {
		t.Errorf("expected not nil Sign")
		return
	}
	assert.ElementsMatch(t, note.Info, newNote.Info)
	assert.ElementsMatch(t, note.Sign, newNote.Sign)
}
func TestOpenNoteError(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	data := []byte{0x00, 0xff}
	note := p.OpenNote(data)
	if note != nil {
		t.Errorf("expected nil note")
	}
}

func TestGetInfoFromNote(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	plainText := "test Plain Text"
	description := "test description"
	metadata := clientCom.NewMetadata("testType", "testUser", time.Now())
	note := p.CreateNewNotes([]byte(plainText), description, metadata)
	actualPlainText, actualMetadata := p.GetInfoFromNote(note)
	assert.Equal(t, plainText, string(actualPlainText))
	assert.Equal(t, "testType", actualMetadata.TypeNote)
}

func TestGetInfoFromNoteWitoutParametrs(t *testing.T) {
	globalNoteManager = nil
	logger := common.FakeLogger{}
	p := NewNoteManager(&logger)
	defer func() { globalNoteManager = nil }()
	login := []byte("testLogin")
	password := []byte{0x00, 0xff, 0x00, 0xff}
	p.InitParametrs(login, password)
	plainText := "test Plain Text"
	description := "test description"
	metadata := clientCom.NewMetadata("testType", "testUser", time.Now())
	note := p.CreateNewNotes([]byte(plainText), description, metadata)
	p.cryptoKey = nil
	actualPlainText, actualMetadata := p.GetInfoFromNote(note)
	if len(actualPlainText) != 0 {
		t.Errorf("expected empty text. actual:%s", string(actualPlainText))
	}
	if actualMetadata != nil {
		t.Errorf(" expected nil pointer actualMetadata")
	}
}
