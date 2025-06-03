package common

import "fmt"

func saveUint64ToSlise(slice []byte, startPos int, value uint64) {
	// temp :=make([]byte, 8)
	for i := 0; i < 8; i++ {
		slice[startPos+i] = uint8((value >> (8 * (7 - i))) & 0xff)
	}
}
func loadUint64FromSlise(slice []byte, startPos int) uint64 {
	var value uint64 = 0
	for i := 0; i < 8; i++ {
		value <<= 8
		value |= uint64(slice[startPos+i])
		// slice[startPos + i] = (value>>(8*(7-i)))&0xff
	}
	return value
}

const SizeField = 8

func SaveDataToStream(stream []byte, actualIndex *int, data []byte) {
	saveUint64ToSlise(stream, *actualIndex, uint64(len(data)))
	*actualIndex += SizeField
	for _, val := range data {
		stream[*actualIndex] = val
		*actualIndex++
	}
}
func LoadDataFromStream(stream []byte, actualIndex *int) ([]byte, error) {
	size := int(loadUint64FromSlise(stream, *actualIndex))
	if len(stream) < size+SizeField {
		return []byte{}, fmt.Errorf("error lenth. actual:%d, expected:%d", len(stream), size+SizeField)
	}
	*actualIndex += SizeField
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = stream[*actualIndex]
		*actualIndex++
	}
	return data, nil
}
