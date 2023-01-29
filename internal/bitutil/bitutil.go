package bitutil

func CombineBytes(op1 byte, op2 byte) uint16 {
	op1Padded := uint16(op1) << 8
	return uint16(op1Padded) | uint16(op2)
}
