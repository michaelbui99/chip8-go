package bitutil

func CombineBytes(msb byte, lsb byte) uint16 {
	msbShifted := uint16(msb) << 8
	return uint16(msbShifted) | uint16(lsb)
}
