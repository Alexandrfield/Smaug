package worker

import (
	"errors"
	"testing"

	clientSecurity "github.com/Alexandrfield/Smaug/internal/client/security"
	mock "github.com/Alexandrfield/Smaug/internal/client/worker/mock"
	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("testError")

func TestRegistration(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTunnel := mock.NewMockTunnelToServer(ctrl)

	secyrityManager := clientSecurity.NewNoteManager(&common.FakeLogger{})
	cl := SmaugClient{logger: &common.FakeLogger{}, secyrityManager: secyrityManager, networkClient: mockTunnel}
	cl.localCache = make(map[string][]*common.Note)
	mockTunnel.EXPECT().Registration(login, gomock.Any()).Return(nil)
	mockTunnel.EXPECT().Login(login, gomock.Any()).Return(nil)
	err := cl.Registration(login, password)
	require.NoError(t, err)
}

func TestRegistrationErr(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTunnel := mock.NewMockTunnelToServer(ctrl)

	secyrityManager := clientSecurity.NewNoteManager(&common.FakeLogger{})
	cl := SmaugClient{logger: &common.FakeLogger{}, secyrityManager: secyrityManager, networkClient: mockTunnel}
	cl.localCache = make(map[string][]*common.Note)
	mockTunnel.EXPECT().Registration(login, gomock.Any()).Return(errTest)
	err := cl.Registration(login, password)
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
func TestRegistrationErrLogin(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTunnel := mock.NewMockTunnelToServer(ctrl)

	secyrityManager := clientSecurity.NewNoteManager(&common.FakeLogger{})
	cl := SmaugClient{logger: &common.FakeLogger{}, secyrityManager: secyrityManager, networkClient: mockTunnel}
	cl.localCache = make(map[string][]*common.Note)
	mockTunnel.EXPECT().Registration(login, gomock.Any()).Return(nil)
	mockTunnel.EXPECT().Login(login, gomock.Any()).Return(errTest)
	err := cl.Registration(login, password)
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}

func TestLoginErr(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTunnel := mock.NewMockTunnelToServer(ctrl)

	secyrityManager := clientSecurity.NewNoteManager(&common.FakeLogger{})
	cl := SmaugClient{logger: &common.FakeLogger{}, secyrityManager: secyrityManager, networkClient: mockTunnel}
	cl.localCache = make(map[string][]*common.Note)
	mockTunnel.EXPECT().Login(login, gomock.Any()).Return(errTest)
	err := cl.login(login, password)
	if err == nil {
		t.Errorf("Unexpected error. expected:not nil; actual:nil")
	}
}
func TestLogin2(t *testing.T) {
	login := "testUser"
	password := []byte{0x00, 0xff}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTunnel := mock.NewMockTunnelToServer(ctrl)

	secyrityManager := clientSecurity.NewNoteManager(&common.FakeLogger{})
	cl := SmaugClient{logger: &common.FakeLogger{}, secyrityManager: secyrityManager, networkClient: mockTunnel}
	cl.localCache = make(map[string][]*common.Note)
	mockTunnel.EXPECT().Login(login, gomock.Any()).Return(nil)
	err := cl.Login(login, password)
	require.NoError(t, err)
}
