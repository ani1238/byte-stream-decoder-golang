package main

import (
	"fmt"
)

type DecodedStruct struct {
	Short1 int16
	Chars1 string
	Byte1  byte
	Chars2 string
	Short2 int16
	Chars3 string
	Long1  int32
}

func decodePacket(packet []byte, resultChan chan DecodedStruct) {
	if len(packet) != 44 {
		resultChan <- DecodedStruct{}
		return
	}

	decoded := DecodedStruct{}

	decoded.Short1 = int16(packet[0])<<8 | int16(packet[1])
	decoded.Chars1 = string(packet[2:14])
	decoded.Byte1 = packet[14]
	decoded.Chars2 = string(packet[15:23])
	decoded.Short2 = int16(packet[23])<<8 | int16(packet[24])
	decoded.Chars3 = string(packet[25:40])
	decoded.Long1 = int32(packet[40])<<24 | int32(packet[41])<<16 | int32(packet[42])<<8 | int32(packet[43])

	resultChan <- decoded
}

func main() {
	//ideally this should be array of packets with packets coming in every nanosecond into this system
	packet := []byte{0x04, 0xD2, 0x6B, 0x65, 0x65, 0x70, 0x64, 0x65, 0x63, 0x6F, 0x64, 0x69, 0x6E, 0x67, 0x38, 0x64, 0x6F, 0x6E, 0x74, 0x73, 0x74, 0x6F, 0x70, 0x03, 0x15, 0x63, 0x6F, 0x6E, 0x67, 0x72, 0x61, 0x74, 0x75, 0x6C, 0x61, 0x74, 0x69, 0x6F, 0x6E, 0x73, 0x07, 0x5B, 0xCD, 0x15}

	numPackets := 1000 // Example: Process 1000 packets concurrently

	resultChan := make(chan DecodedStruct, numPackets)

	for i := 0; i < numPackets; i++ {
		go decodePacket(packet, resultChan) //making worker pool
	}

	for i := 0; i < numPackets; i++ {
		decoded := <-resultChan
		if decoded.Short1 != 0 {
			fmt.Printf("Decoded struct: %+v\n", decoded)
		}
	}
	close(resultChan)
}
