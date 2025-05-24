package common

func SaveUint64ToSlise(slice []byte, startPos int, value uint64) {
	// temp :=make([]byte, 8)
	for i := 0; i < 8; i++ {
		slice[startPos+i] = uint8((value >> (8 * (7 - i))) & 0xff)
	}
}
func LoadUint64FromSlise(slice []byte, startPos int) uint64 {
	var value uint64 = 0
	for i := 0; i < 8; i++ {
		value <<= 8
		value |= uint64(slice[startPos+i])
		// slice[startPos + i] = (value>>(8*(7-i)))&0xff
	}
	return value
}
