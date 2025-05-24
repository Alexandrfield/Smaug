package security

func CreatNote(info []byte, temporyKey []byte, cipherText []byte, sign []byte) *Note {
	return &Note{info: info, temporyKey: temporyKey, cipherText: cipherText, sign: sign}
}

type Note struct {
	info       []byte
	temporyKey []byte
	cipherText []byte
	sign       []byte
}

func (note *Note) Serialize() []byte {
	length := (sizeField + len(note.info)) + (sizeField + len(note.temporyKey)) +
		(sizeField + len(note.cipherText)) + (sizeField + len(note.sign))

	bytesStream := make([]byte, length)
	actualIndex := 0
	saveDataToStream(bytesStream, &actualIndex, note.info)
	saveDataToStream(bytesStream, &actualIndex, note.temporyKey)
	saveDataToStream(bytesStream, &actualIndex, note.cipherText)
	saveDataToStream(bytesStream, &actualIndex, note.sign)
	return bytesStream
}

func (note *Note) Deserialize(data []byte) error {
	actualIndex := 0
	note.info, _ = loadDataFromStream(data, &actualIndex)
	note.temporyKey, _ = loadDataFromStream(data, &actualIndex)
	note.cipherText, _ = loadDataFromStream(data, &actualIndex)
	note.sign, _ = loadDataFromStream(data, &actualIndex)
	return nil
}
