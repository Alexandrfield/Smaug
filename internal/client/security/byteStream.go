package security

import (
	"fmt"

	"github.com/Alexandrfield/Smaug/internal/common"
)

const sizeField = 8

func saveDataToStream(stream []byte, actualIndex *int, data []byte) {
	common.SaveUint64ToSlise(stream, *actualIndex, uint64(len(data)))
	*actualIndex += sizeField
	for _, val := range data {
		stream[*actualIndex] = val
		*actualIndex++
	}
}
func loadDataFromStream(stream []byte, actualIndex *int) ([]byte, error) {
	size := int(common.LoadUint64FromSlise(stream, *actualIndex))
	if len(stream) < size+sizeField {
		return []byte{}, fmt.Errorf("error lenth. actual:%d, expected:%d", len(stream), size+sizeField)
	}
	*actualIndex += sizeField
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = stream[*actualIndex]
		*actualIndex++
	}
	return data, nil
}
