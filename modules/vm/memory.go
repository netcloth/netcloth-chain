package vm

import (
	"fmt"
	"math/big"

	"github.com/netcloth/netcloth-chain/modules/vm/common/math"
)

type Memory struct {
	store       []byte
	lastGasCost uint64
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Set(offset, size uint64, value []byte) {
	if size > 0 {
		if offset+size > uint64(len(m.store)) {
			panic("invalid memory: store empty")
		}
		copy(m.store[offset:offset+size], value)
	}
}

// Set32 sets the 32 bytes starting at offset to the value of val, left-padded with zeros to 32 bytes
func (m *Memory) Set32(offset uint64, val *big.Int) {
	if offset+32 > uint64(len(m.store)) {
		panic("invalid memory: store empty")
	}

	// Zero the memory ares
	copy(m.store[offset:offset+32], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	// Fill in relevant bits
	math.ReadBits(val, m.store[offset:offset+32])
}

//GetCopy
func (m *Memory) GetCopy(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.store[offset:offset+size])

		return
	}
	return
}

// Resize resizes the memory to size
func (m *Memory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	}
}

// Get returns offset + size as new slice
func (m *Memory) Get(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.store[offset:offset+size])
		return
	}
	return
}

// GetPtr returns the offset +size
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}
	if len(m.store) > int(offset) {
		return m.store[offset : offset+size]
	}

	return nil
}

// Len returns the length of the backing slice
func (m *Memory) Len() int {
	return len(m.store)
}

// Data returns the backing slice
func (m *Memory) Data() []byte {
	return m.store
}

func (m *Memory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.store)
	if len(m.store) > 0 {
		addr := 0
		for i := 0; i +32 <= len(m.store); i +=32 {
			fmt.Printf("%03d: %x\n", addr, m.store[i:i+32])
			addr++
		} else {
				fmt.Println("--- empty ----")
		}
	}
	fmt.Println("####################")
}
