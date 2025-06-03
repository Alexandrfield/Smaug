package worker

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetadata(t *testing.T) {
	user := "testUser"
	timeCreate := time.Now()
	typeTest := "testType"
	metadata := NewMetadata(typeTest, user, timeCreate)
	assert.Equal(t, user, metadata.User)
	assert.Equal(t, timeCreate.Format("2006-01-02 15:04:05"), metadata.TimeCreate)
	assert.Equal(t, typeTest, metadata.TypeNote)
}

func TestGetInfo(t *testing.T) {
	user := "testUser"
	timeCreate := time.Now()
	typeTest := "testType"
	expected := fmt.Sprintf("type:testType; create:%s; user:testUser;", timeCreate.Format("2006-01-02 15:04:05"))
	metadata := NewMetadata(typeTest, user, timeCreate)
	actual := metadata.GetInfo()
	assert.Equal(t, expected, actual)
}

func TestSerializeDeserializeMetadata(t *testing.T) {
	user := "testUser"
	timeCreate := time.Now()
	typeTest := "testType"
	metadata := NewMetadata(typeTest, user, timeCreate)
	serByteStream := SerializeMetadata(metadata)
	newMetadata, err := DeserializeMetadata(serByteStream)
	require.NoError(t, err)
	assert.Equal(t, user, newMetadata.User)
	assert.Equal(t, timeCreate.Format("2006-01-02 15:04:05"), newMetadata.TimeCreate)
	assert.Equal(t, typeTest, newMetadata.TypeNote)
}
func TestSerializeDeserializeMetadataError(t *testing.T) {
	serByteStream := []byte{0x00, 0x11}
	_, err := DeserializeMetadata(serByteStream)
	if err == nil {
		t.Errorf("expected not nil error!")
	}
}
