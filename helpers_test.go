package ssz

import (
	"reflect"
	"testing"
)

func TestPack_NoItems(t *testing.T) {
	output, err := pack([][]byte{})
	if err != nil {
		t.Fatalf("pack() error = %v", err)
	}
	if len(output[0]) != BytesPerChunk {
		t.Errorf("Expected empty input to return an empty chunk, received %v", output)
	}
}

func TestPack_ExactBytePerChunkLength(t *testing.T) {
	input := [][]byte{}
	for i := 0; i < 10; i++ {
		item := make([]byte, BytesPerChunk)
		input = append(input, item)
	}
	output, err := pack(input)
	if err != nil {
		t.Fatalf("pack() error = %v", err)
	}
	if len(output) != 10 {
		t.Errorf("Expected empty input to return an empty chunk, received %v", output)
	}
	if !reflect.DeepEqual(output, input) {
		t.Errorf("pack() = %v, want %v", output, input)
	}
}

func TestPack_OK(t *testing.T) {
	tests := []struct {
		name   string
		input  [][]byte
		output [][]byte
	}{
		{
			name:   "an item having less than BytesPerChunk should return a padded chunk",
			input:  [][]byte{make([]byte, BytesPerChunk-4)},
			output: [][]byte{make([]byte, BytesPerChunk)},
		},
		{
			name:   "two items having less than BytesPerChunk should return two chunks",
			input:  [][]byte{make([]byte, BytesPerChunk-5), make([]byte, BytesPerChunk-5)},
			output: [][]byte{make([]byte, BytesPerChunk), make([]byte, BytesPerChunk)},
		},
		{
			name:   "two items with length BytesPerChunk/2 should return one chunk",
			input:  [][]byte{make([]byte, BytesPerChunk/2), make([]byte, BytesPerChunk/2)},
			output: [][]byte{make([]byte, BytesPerChunk)},
		},
		{
			name:   "an item with length BytesPerChunk*2 should return two chunks",
			input:  [][]byte{make([]byte, BytesPerChunk*2)},
			output: [][]byte{make([]byte, BytesPerChunk), make([]byte, BytesPerChunk)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pack(tt.input)
			if err != nil {
				t.Fatalf("pack() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("pack() = %v, want %v", got, tt.output)
			}
		})
	}
}

func TestMerkleize_Identity(t *testing.T) {
	input := [][]byte{make([]byte, BytesPerChunk)}
	output := merkleize(input)
	if !reflect.DeepEqual(output[:], input[0]) {
		t.Errorf("merkleize() = %v, want %v", output, input)
	}
}

func TestMerkleize_OK(t *testing.T) {
	chunk := make([]byte, BytesPerChunk)
	secondLayerRoot := Hash(append(chunk, chunk...))
	thirdLayerRoot := Hash(append(secondLayerRoot[:], secondLayerRoot[:]...))
	tests := []struct {
		name   string
		input  [][]byte
		output [32]byte
	}{
		{
			name:   "two elements should return the hash of their concatenation",
			input:  [][]byte{make([]byte, BytesPerChunk), make([]byte, BytesPerChunk)},
			output: Hash(make([]byte, BytesPerChunk*2)),
		},
		{
			name:   "four chunks should return the Merkle root of a three layer trie",
			input:  [][]byte{chunk, chunk, chunk, chunk},
			output: thirdLayerRoot,
		},
		{
			name:   "three chunks should pad until there are four chunks",
			input:  [][]byte{chunk, chunk, chunk},
			output: thirdLayerRoot,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := merkleize(tt.input)
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("merkleize() = %v, want %v", got, tt.output)
			}
		})
	}
}

func TestIsPowerTwo(t *testing.T) {
	tests := []struct {
		input  int
		output bool
	}{
		{input: 4, output: true},
		{input: 5, output: false},
		{input: 1, output: true},
		{input: 0, output: false},
		{input: 2, output: true},
		{input: 256, output: true},
		{input: 1024, output: true},
		{input: 1000000, output: false},
	}
	for _, tt := range tests {
		got := isPowerTwo(tt.input)
		if got != tt.output {
			t.Errorf("isPowerTwo() = %v, want %v", got, tt.output)
		}
	}
}

func BenchmarkPack(b *testing.B) {
	input := [][]byte{make([]byte, BytesPerChunk*8000)}
	for n := 0; n < b.N; n++ {
		if _, err := pack(input); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMerkleize(b *testing.B) {
	input := make([][]byte, 1000)
	for i := 0; i < len(input); i++ {
		input[i] = make([]byte, BytesPerChunk)
	}
	for n := 0; n < b.N; n++ {
		merkleize(input)
	}
}

func BenchmarkIsPowerTwo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		isPowerTwo(1 << 36)
	}
}
