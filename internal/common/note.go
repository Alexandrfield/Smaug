package common

import "fmt"

func CreatNote(info []byte, temporyKey []byte, cipherText []byte) *Note {
	return &Note{Info: info, TemporaryKey: temporyKey, CipherText: cipherText}
}

type Note struct {
	Info         []byte
	TemporaryKey []byte
	CipherText   []byte
	Sign         []byte
}

func (note *Note) GetDataForSign() []byte {
	var temp []byte
	temp = append(temp, note.Info...)
	temp = append(temp, note.TemporaryKey...)
	temp = append(temp, note.CipherText...)
	return temp
}
func (note *Note) Serialize() []byte {
	length := (SizeField + len(note.Info)) + (SizeField + len(note.TemporaryKey)) +
		(SizeField + len(note.CipherText)) + (SizeField + len(note.Sign))

	bytesStream := make([]byte, length)
	actualIndex := 0
	SaveDataToStream(bytesStream, &actualIndex, note.Info)
	SaveDataToStream(bytesStream, &actualIndex, note.TemporaryKey)
	SaveDataToStream(bytesStream, &actualIndex, note.CipherText)
	SaveDataToStream(bytesStream, &actualIndex, note.Sign)
	return bytesStream
}

func (note *Note) Deserialize(data []byte) error {
	actualIndex := 0
	if len(data) < 32 {
		return fmt.Errorf("error data length for deserialize. data:%x", data)
	}
	note.Info, _ = LoadDataFromStream(data, &actualIndex)
	note.TemporaryKey, _ = LoadDataFromStream(data, &actualIndex)
	note.CipherText, _ = LoadDataFromStream(data, &actualIndex)
	note.Sign, _ = LoadDataFromStream(data, &actualIndex)
	return nil
}
