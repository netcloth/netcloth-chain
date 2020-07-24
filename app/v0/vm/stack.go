package vm

import (
	"fmt"
	"math/big"
)

type Stack struct {
	data []*big.Int
}

func newstack() *Stack {
	return &Stack{data: make([]*big.Int, 0, 1024)}
}

func (st *Stack) Data() []*big.Int {
	return st.data
}

func (st *Stack) push(d *big.Int) {
	st.data = append(st.data, d)
}

func (st *Stack) pushN(ds ...*big.Int) {
	st.data = append(st.data, ds...)
}

func (st *Stack) pop() (ret *big.Int) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

func (st *Stack) len() int {
	return len(st.data)
}

func (st *Stack) swap(n int) {
	st.data[st.len()-n], st.data[st.len()-1] = st.data[st.len()-1], st.data[st.len()-n]
}

func (st *Stack) dup(pool *intPool, n int) {
	st.push(pool.get().Set(st.data[st.len()-n]))
}

func (st *Stack) peek() *big.Int {
	return st.data[st.len()-1]
}

// returns the n'th item in stack
func (st *Stack) Back(n int) *big.Int {
	return st.data[st.len()-n-1]
}

func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if st.len() > 0 {
		// print stack from top
		for j := len(st.data) - 1; j >= 0; j-- {
			fmt.Printf("%-3d %064x\n", len(st.data)-j, (st.data[j]))
		}

		// print stack from bottom
		//for i, val := range st.data {
		//	fmt.Printf("%-3d %d\n", i, val)
		//}
	} else {
		fmt.Println("--- empty ---")
	}
	fmt.Println("#############")
}
