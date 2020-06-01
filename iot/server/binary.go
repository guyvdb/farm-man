package server


import (
	"encoding/binary"
)

// some useful binary functions


func bytes_uint8_decode(data []byte) uint8 {
	return uint8(data[0])
}

func bytes_uint16_decode(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func bytes_uint32_decode(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func bytes_uint8_encode(value uint8) []byte {
	return []byte{byte(value)}
}

func bytes_uint16_encode(value uint16) []byte {
	result := make([]byte,2)
	binary.BigEndian.PutUint16(result, value)
	return result 
}

func bytes_uint32_encode(value uint32) []byte {
	result := make([]byte,4)
	binary.BigEndian.PutUint32(result, value)
	return result 	
}


