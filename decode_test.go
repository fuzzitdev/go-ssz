package ssz

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	//forkExample := fork{
	//	PreviousVersion: [4]byte{2, 3, 4, 1},
	//	CurrentVersion:  [4]byte{},
	//	Epoch:           10923910294,
	//}
	tests := []struct {
		input interface{}
		ptr   interface{}
	}{
		// Bool test cases.
		{input: true, ptr: new(bool)},
		{input: false, ptr: new(bool)},
		//// Uint8 test cases.
		{input: byte(1), ptr: new(byte)},
		{input: byte(0), ptr: new(byte)},
		// Uint16 test cases.
		{input: uint16(100), ptr: new(uint16)},
		{input: uint16(232), ptr: new(uint16)},
		// Uint32 test cases.
		{input: uint32(1), ptr: new(uint32)},
		{input: uint32(1029391), ptr: new(uint32)},
		// Uint64 test cases.
		{input: uint64(5), ptr: new(uint64)},
		{input: uint64(23929309), ptr: new(uint64)},
		// Byte slice, byte array test cases.
		{input: [8]byte{1, 2, 3, 4, 5, 6, 7, 8}, ptr: new([8]byte)},
		{input: []byte{9, 8, 9, 8}, ptr: new([]byte)},
		// Basic type array test cases.
		{input: [12]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ptr: new([12]uint64)},
		{input: [100]bool{true, false, true, true}, ptr: new([100]bool)},
		{input: [20]uint16{3, 4, 5}, ptr: new([20]uint16)},
		{input: [20]uint32{4, 5}, ptr: new([20]uint32)},
		// Basic type slice test cases.
		{input: []uint64{1, 2, 3}, ptr: new([]uint64)},
		{input: []bool{true, false, true, true, true}, ptr: new([]bool)},
		{input: []uint32{0, 0, 0}, ptr: new([]uint32)},
		{input: []uint32{92939, 232, 222}, ptr: new([]uint32)},
		//// Struct decoding test cases.
		//{input: forkExample, ptr: &fork{}},
		// Non-basic type slice/array test cases.
		//{input: []fork{forkExample, forkExample}, ptr: new([]fork)},
		{input: [][]uint64{{4, 3, 2}, {1}, {0}}, ptr: new([][]uint64)},
		{input: [][][]uint64{{{1, 2}, {3}}, {{4, 5}}, {{0}}}, ptr: new([][][]uint64)},
		//{input: [][3]uint64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, ptr: new([][3]uint64)},
		//{input: [4]fork{forkExample, forkExample, forkExample}, ptr: new([4]fork)},
	}
	for _, tt := range tests {
		buffer := new(bytes.Buffer)
		if err := Encode(buffer, tt.input); err != nil {
			panic(err)
		}
		fmt.Printf("Encoded: %v\n", buffer.Bytes())
		if err := Decode(buffer.Bytes(), tt.ptr); err != nil {
			t.Fatal(err)
		}
		output := reflect.ValueOf(tt.ptr).Elem().Interface()
		fmt.Printf("Decoded: %v\n", output)
		if !reflect.DeepEqual(output, tt.input) {
			t.Errorf("Expected %d, received %d", tt.input, output)
		}
	}
}
