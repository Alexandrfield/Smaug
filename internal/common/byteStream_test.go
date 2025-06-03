package common

import (
	"bytes"
	"testing"
)

func TestSaveDataToStream(t *testing.T) {
	data := []byte{0x01, 0x02, 0x10, 0x20, 0xff}
	actualIndex := 0
	expectedLength := SizeField + len(data)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x10, 0x20, 0xff}
	actual := make([]byte, expectedLength)
	SaveDataToStream(actual, &actualIndex, data)
	if actualIndex != expectedLength {
		t.Errorf("error length. expectd:%d, actual^%d", expectedLength, actualIndex)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Result was incorrect. expected:%x; actual:%x", expected, actual)
	}
}

func TestLoadDataFromStream(t *testing.T) {
	actualIndex := 0
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x10, 0x20, 0xff}
	expected := []byte{0x01, 0x02, 0x10, 0x20, 0xff}
	actual, err := LoadDataFromStream(data, &actualIndex)
	if err != nil {
		t.Errorf("Result was incorrect. expected error:nil; actual:%s", err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Result was incorrect. expected:%x; actual:%x", expected, actual)
	}
}

func TestLoadDataFromStreamErrorLength(t *testing.T) {
	actualIndex := 0
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x10}
	_, err := LoadDataFromStream(data, &actualIndex)
	if err == nil {
		t.Errorf("Result was incorrect. expected error:not nil; actual:nil")
	}
}
