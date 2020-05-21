package vm

import (
	"encoding/json"
	"fmt"
	"github.com/netcloth/netcloth-chain/hexutil"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types/time"
	"io/ioutil"
	"math/big"
	"reflect"
	"testing"

	"github.com/netcloth/netcloth-chain/app/v0/vm/common"
)

type TwoOperandTestcase struct {
	X        string
	Y        string
	Expected string
}

type twoOperandParams struct {
	x string
	y string
}

var commonParams []*twoOperandParams
var twoOpMethods map[string]executionFunc

type OneOperandTestcase struct {
	X        string
	Expected string
}

var oneOpParams []string
var oneOpMethods map[string]executionFunc

func init() {

	// Params is a list of common edgecases that should be used for some common tests
	params := []string{
		"0000000000000000000000000000000000000000000000000000000000000000", // 0
		"0000000000000000000000000000000000000000000000000000000000000001", // +1
		"0000000000000000000000000000000000000000000000000000000000000005", // +5
		"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", // + max -1
		"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // + max
		"8000000000000000000000000000000000000000000000000000000000000000", // - max
		"8000000000000000000000000000000000000000000000000000000000000001", // - max+1
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb", // - 5
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // - 1
	}
	// Params are combined so each param is used on each 'side'
	commonParams = make([]*twoOperandParams, len(params)*len(params))
	for i, x := range params {
		for j, y := range params {
			commonParams[i*len(params)+j] = &twoOperandParams{x, y}
		}
	}
	twoOpMethods = map[string]executionFunc{
		"add":     opAdd,
		"sub":     opSub,
		"mul":     opMul,
		"div":     opDiv,
		"sdiv":    opSdiv,
		"mod":     opMod,
		"smod":    opSmod,
		"exp":     opExp,
		"signext": opSignExtend,
		"lt":      opLt,
		"gt":      opGt,
		"slt":     opSlt,
		"sgt":     opSgt,
		"eq":      opEq,
		"and":     opAnd,
		"or":      opOr,
		"xor":     opXor,
		"byte":    opByte,
		"shl":     opSHL,
		"shr":     opSHR,
		"sar":     opSAR,
	}

	oneOpMethods = map[string]executionFunc{
		"iszero": opIszero,
	}
}

func testTwoOperandOp(t *testing.T, tests []TwoOperandTestcase, opFn executionFunc, name string) {

	var (
		env            = newEVM()
		stack          = newstack()
		pc             = uint64(0)
		evmInterpreter = env.interpreter.(*EVMInterpreter)
	)
	// Stuff a couple of nonzero bigints into pool, to ensure that ops do not rely on pooled integers to be zero
	evmInterpreter.intPool = poolOfIntPools.get()
	evmInterpreter.intPool.put(big.NewInt(-1337))
	evmInterpreter.intPool.put(big.NewInt(-1337))
	evmInterpreter.intPool.put(big.NewInt(-1337))

	for i, test := range tests {
		x := new(big.Int).SetBytes(common.Hex2Bytes(test.X))
		y := new(big.Int).SetBytes(common.Hex2Bytes(test.Y))
		expected := new(big.Int).SetBytes(common.Hex2Bytes(test.Expected))
		stack.push(x)
		stack.push(y)
		opFn(&pc, evmInterpreter, nil, nil, stack)
		actual := stack.pop()

		if actual.Cmp(expected) != 0 {
			t.Errorf("Testcase %v %d, %v(%x, %x): expected  %x, got %x", name, i, name, x, y, expected, actual)
		}
		// Check pool usage
		// 1.pool is not allowed to contain anything on the stack
		// 2.pool is not allowed to contain the same pointers twice
		if evmInterpreter.intPool.pool.len() > 0 {

			poolvals := make(map[*big.Int]struct{})
			poolvals[actual] = struct{}{}

			for evmInterpreter.intPool.pool.len() > 0 {
				key := evmInterpreter.intPool.get()
				if _, exist := poolvals[key]; exist {
					t.Errorf("Testcase %v %d, pool contains double-entry", name, i)
				}
				poolvals[key] = struct{}{}
			}
		}
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func testOneOperandOp(t *testing.T, tests []OneOperandTestcase, opFn executionFunc, name string) {

	var (
		env            = newEVM()
		stack          = newstack()
		pc             = uint64(0)
		evmInterpreter = env.interpreter.(*EVMInterpreter)
	)
	// Stuff a couple of nonzero bigints into pool, to ensure that ops do not rely on pooled integers to be zero
	evmInterpreter.intPool = poolOfIntPools.get()
	evmInterpreter.intPool.put(big.NewInt(-1337))
	evmInterpreter.intPool.put(big.NewInt(-1337))
	evmInterpreter.intPool.put(big.NewInt(-1337))

	for i, test := range tests {
		x := new(big.Int).SetBytes(common.Hex2Bytes(test.X))
		expected := new(big.Int).SetBytes(common.Hex2Bytes(test.Expected))
		stack.push(x)
		opFn(&pc, evmInterpreter, nil, nil, stack)
		actual := stack.pop()

		if actual.Cmp(expected) != 0 {
			t.Errorf("Testcase %v %d, %v(%x): expected  %x, got %x", name, i, name, x, expected, actual)
		}
		// Check pool usage
		// 1.pool is not allowed to contain anything on the stack
		// 2.pool is not allowed to contain the same pointers twice
		if evmInterpreter.intPool.pool.len() > 0 {

			poolvals := make(map[*big.Int]struct{})
			poolvals[actual] = struct{}{}

			for evmInterpreter.intPool.pool.len() > 0 {
				key := evmInterpreter.intPool.get()
				if _, exist := poolvals[key]; exist {
					t.Errorf("Testcase %v %d, pool contains double-entry", name, i)
				}
				poolvals[key] = struct{}{}
			}
		}
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

// test opByte instruction
func TestByteOp(t *testing.T) {
	tests := []TwoOperandTestcase{
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", "00", "AB"},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", "01", "CD"},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", "00", "00"},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", "01", "CD"},
		{"0000000000000000000000000000000000000000000000000000000000102030", "1F", "30"},
		{"0000000000000000000000000000000000000000000000000000000000102030", "1E", "20"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "20", "00"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "FFFFFFFFFFFFFFFF", "00"},
	}
	testTwoOperandOp(t, tests, opByte, "byte")
}

// test opSHL instruction
func TestSHL(t *testing.T) {
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000002"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "ff", "8000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "0101", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "8000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
	}
	testTwoOperandOp(t, tests, opSHL, "shl")
}

func TestSHR(t *testing.T) {
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "00", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "01", "4000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "ff", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0101", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
	}
	testTwoOperandOp(t, tests, opSHR, "shr")
}

func TestSAR(t *testing.T) {
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "00", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "01", "c000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "ff", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0100", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0101", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000000000000000000000", "fe", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "f8", "000000000000000000000000000000000000000000000000000000000000007f"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "fe", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
	}

	testTwoOperandOp(t, tests, opSAR, "sar")
}

// getResult is a convenience function to generate the expected values
func getResult(args []*twoOperandParams, opFn executionFunc) []TwoOperandTestcase {
	var (
		env         = newEVM()
		stack       = newstack()
		pc          = uint64(0)
		interpreter = env.interpreter.(*EVMInterpreter)
	)
	interpreter.intPool = poolOfIntPools.get()
	result := make([]TwoOperandTestcase, len(args))
	for i, param := range args {
		x := new(big.Int).SetBytes(common.Hex2Bytes(param.x))
		y := new(big.Int).SetBytes(common.Hex2Bytes(param.y))
		stack.push(x)
		stack.push(y)
		opFn(&pc, interpreter, nil, nil, stack)
		actual := stack.pop()
		result[i] = TwoOperandTestcase{param.x, param.y, fmt.Sprintf("%064x", actual)}
	}
	return result
}

// utility function to fill the json-file with testcases
// Enable this test to generate the 'testcases_xx.json' files
func TestWriteExpectedValues(t *testing.T) {
	t.Skip("Enable this test to create json test cases.")

	for name, method := range twoOpMethods {
		data, err := json.Marshal(getResult(commonParams, method))
		if err != nil {
			t.Fatal(err)
		}
		_ = ioutil.WriteFile(fmt.Sprintf("testdata/testcases_%v.json", name), data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestJsonTestcases runs through all the testcases defined as json-files
func TestJsonTestcases(t *testing.T) {
	for name := range twoOpMethods {
		data, err := ioutil.ReadFile(fmt.Sprintf("testdata/testcases_%v.json", name))
		if err != nil {
			t.Fatal("Failed to read file", err)
		}
		var testcases []TwoOperandTestcase
		json.Unmarshal(data, &testcases)
		testTwoOperandOp(t, testcases, twoOpMethods[name], name)
	}
}

func opBenchmark(bench *testing.B, op func(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error), args ...string) {
	var (
		env            = newEVM()
		stack          = newstack()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	// convert args
	byteArgs := make([][]byte, len(args))
	for i, arg := range args {
		byteArgs[i] = common.Hex2Bytes(arg)
	}
	pc := uint64(0)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		for _, arg := range byteArgs {
			a := new(big.Int).SetBytes(arg)
			stack.push(a)
		}
		op(&pc, evmInterpreter, nil, nil, stack)
		stack.pop()
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpAdd64(b *testing.B) {
	x := "ffffffff"
	y := "fd37f3e2bba2c4f"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpAdd128(b *testing.B) {
	x := "ffffffffffffffff"
	y := "f5470b43c6549b016288e9a65629687"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpAdd256(b *testing.B) {
	x := "0802431afcbce1fc194c9eaa417b2fb67dc75a95db0bc7ec6b1c8af11df6a1da9"
	y := "a1f5aac137876480252e5dcac62c354ec0d42b76b0642b6181ed099849ea1d57"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpSub64(b *testing.B) {
	x := "51022b6317003a9d"
	y := "a20456c62e00753a"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpSub128(b *testing.B) {
	x := "4dde30faaacdc14d00327aac314e915d"
	y := "9bbc61f5559b829a0064f558629d22ba"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpSub256(b *testing.B) {
	x := "4bfcd8bb2ac462735b48a17580690283980aa2d679f091c64364594df113ea37"
	y := "97f9b1765588c4e6b69142eb00d20507301545acf3e1238c86c8b29be227d46e"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpMul(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMul, x, y)
}

func BenchmarkOpDiv256(b *testing.B) {
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpDiv128(b *testing.B) {
	x := "fdedc7f10142ff97"
	y := "fbdfda0e2ce356173d1993d5f70a2b11"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpDiv64(b *testing.B) {
	x := "fcb34eb3"
	y := "f97180878e839129"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpSdiv(b *testing.B) {
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"

	opBenchmark(b, opSdiv, x, y)
}

func BenchmarkOpMod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMod, x, y)
}

func BenchmarkOpSmod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSmod, x, y)
}

func BenchmarkOpExp(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opExp, x, y)
}

func BenchmarkOpSignExtend(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSignExtend, x, y)
}

func BenchmarkOpLt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opLt, x, y)
}

func BenchmarkOpGt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opGt, x, y)
}

func BenchmarkOpSlt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSlt, x, y)
}

func BenchmarkOpSgt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSgt, x, y)
}

func BenchmarkOpEq(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opEq, x, y)
}

func BenchmarkOpEq2(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201fffffffe"
	opBenchmark(b, opEq, x, y)
}

func BenchmarkOpAnd(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opAnd, x, y)
}

func BenchmarkOpOr(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opOr, x, y)
}

func BenchmarkOpXor(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opXor, x, y)
}

func BenchmarkOpByte(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opByte, x, y)
}

func BenchmarkOpAddmod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	z := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opAddmod, x, y, z)
}

func BenchmarkOpMulmod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	z := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMulmod, x, y, z)
}

func BenchmarkOpSHL(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSHL, x, y)
}

func BenchmarkOpSHR(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSHR, x, y)
}

func BenchmarkOpSAR(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSAR, x, y)
}

func BenchmarkOpIsZero(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	opBenchmark(b, opIszero, x)
}

func TestOpMstore(t *testing.T) {
	var (
		env            = newEVM()
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	mem.Resize(64)
	pc := uint64(0)
	v := "abcdef00000000000000abba000000000deaf000000c0de00100000000133700"
	stack.pushN(new(big.Int).SetBytes(common.Hex2Bytes(v)), big.NewInt(0))
	opMstore(&pc, evmInterpreter, nil, mem, stack)
	if got := common.Bytes2Hex(mem.GetCopy(0, 32)); got != v {
		t.Fatalf("Mstore fail, got %v, expected %v", got, v)
	}
	stack.pushN(big.NewInt(0x1), big.NewInt(0))
	opMstore(&pc, evmInterpreter, nil, mem, stack)
	if common.Bytes2Hex(mem.GetCopy(0, 32)) != "0000000000000000000000000000000000000000000000000000000000000001" {
		t.Fatalf("Mstore failed to overwrite previous value")
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpMstore(bench *testing.B) {
	bench.SetParallelism(1)
	var (
		env            = newEVM()
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	mem.Resize(64)
	pc := uint64(0)
	memStart := big.NewInt(0)
	value := big.NewInt(0x1337)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.pushN(value, memStart)
		opMstore(&pc, evmInterpreter, nil, mem, stack)
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpSHA3(bench *testing.B) {
	var (
		env            = newEVM()
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)
	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	mem.Resize(32)
	pc := uint64(0)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.pushN(big.NewInt(32), big.NewInt(0))
		opSha3(&pc, evmInterpreter, nil, mem, stack)
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

//OneOp test

func TestIsZero(t *testing.T) {
	tests := []OneOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000000", "01"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "00"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "00"},
	}
	testOneOperandOp(t, tests, opIszero, "opIszero")
}

func TestNot(t *testing.T) {
	tests := []OneOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000000", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"1000000000000000000000000000000000000000000000000000000000000002", "effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd"},
	}
	testOneOperandOp(t, tests, opNot, "opNot")
}

func TestOpAddress(t *testing.T) {
	addr := sdk.AccAddress{0xab}
	var (
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{}, &dummyContractRef{address: addr}, new(big.Int), 0)
	)

	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()
	opAddress(&pc, interpreter, contract, mem, stack)

	actualAddr := sdk.AccAddress(stack.pop().Bytes())
	if !actualAddr.Equals(addr) {
		t.Errorf("Address fail, got %x, expected %x", actualAddr, addr)
	}
}

func TestOpBalance(t *testing.T) {
	var (
		addr    = sdk.AccAddress{0xab}
		balance = big.NewInt(100)

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{}, &dummyContractRef{address: addr}, new(big.Int), 0)
	)

	env.StateDB.SetBalance(addr, balance)

	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()
	stack.push(big.NewInt(0).SetBytes(addr))

	opBalance(&pc, interpreter, contract, mem, stack)

	actualBalance := stack.pop()
	if actualBalance.Cmp(balance) != 0 {
		t.Errorf("Balance fail, got %d, expected %d", actualBalance, balance)
	}
}

func TestOpCaller(t *testing.T) {
	var (
		expectedCallerAddress = sdk.AccAddress{0xab}

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: expectedCallerAddress}, &dummyContractRef{address: expectedCallerAddress}, new(big.Int), 0)
	)

	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	opCaller(&pc, interpreter, contract, mem, stack)

	actualCallerAddress := sdk.AccAddress(stack.pop().Bytes())
	require.True(t, actualCallerAddress.Equals(expectedCallerAddress))
}

func TestOpCallValue(t *testing.T) {
	var (
		addr          = sdk.AccAddress{0xab}
		expectedValue = big.NewInt(1000)

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, expectedValue, 0)
	)

	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	opCallValue(&pc, interpreter, contract, mem, stack)

	actualValue := big.NewInt(0).SetBytes(stack.pop().Bytes())
	require.True(t, actualValue.Cmp(expectedValue) == 0)
}

func TestOpCodeSize(t *testing.T) {
	var (
		addr             = sdk.AccAddress{0xab}
		value            = big.NewInt(1000)
		code             = []byte("abc")
		expectedCodeSize = big.NewInt(int64(len(code)))

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	contract.SetCallCode(&addr, sdk.Hash{}, code)
	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	opCodeSize(&pc, interpreter, contract, mem, stack)

	actualCodeSize := big.NewInt(0).SetBytes(stack.pop().Bytes())
	require.True(t, actualCodeSize.Cmp(expectedCodeSize) == 0)
}

func TestOpCodeCopy(t *testing.T) {
	var (
		addr  = sdk.AccAddress{0xab}
		value = big.NewInt(1000)
		code  = []byte("abc")

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	contract.SetCallCode(&addr, sdk.Hash{}, code)
	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()
	mem.Resize(64)

	stack.push(big.NewInt(int64(len(code))))
	stack.push(big.NewInt(0))
	stack.push(big.NewInt(0))
	opCodeCopy(&pc, interpreter, contract, mem, stack)
	require.True(t, reflect.DeepEqual(mem.Data()[:3], code))

	stack.push(big.NewInt(int64(len(code))))
	stack.push(big.NewInt(0))
	stack.push(big.NewInt(10))
	opCodeCopy(&pc, interpreter, contract, mem, stack)
	require.True(t, reflect.DeepEqual(mem.Data()[10:13], code))
}

func TestOpGasPrice(t *testing.T) {
	var (
		addr     = sdk.AccAddress{0xab}
		value    = big.NewInt(1000)
		gasPrice = big.NewInt(1000)

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	env.Context.GasPrice = gasPrice
	opGasprice(&pc, interpreter, contract, mem, stack)

	actualGasPrice := big.NewInt(0).SetBytes(stack.pop().Bytes())
	require.True(t, actualGasPrice.Cmp(gasPrice) == 0)
}

func TestOpGasLimit(t *testing.T) {
	var (
		addr     = sdk.AccAddress{0xab}
		value    = big.NewInt(1000)
		gasLimit = uint64(1000)

		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	env.Context.GasLimit = gasLimit
	opGasLimit(&pc, interpreter, contract, mem, stack)

	actualGasLimit := big.NewInt(0).SetBytes(stack.pop().Bytes())
	require.True(t, actualGasLimit.Uint64() == gasLimit)
}

func TestOpPush1(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()

	code, _ := hexutil.Decode("010203")
	contract.SetCallCode(&addr, sdk.Hash{}, code)
	pc := uint64(0)

	opPush1(&pc, interpreter, contract, mem, stack)

	require.Equal(t, pc, uint64(1))
	v := stack.pop().Uint64()
	require.Equal(t, uint64(2), v)

	//
	pc = uint64(100)

	opPush1(&pc, interpreter, contract, mem, stack)

	require.Equal(t, pc, uint64(101))
	v = stack.pop().Uint64()
	require.Equal(t, uint64(0), v)
}

func TestOpPushN(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()

	code, _ := hexutil.Decode("0102030405060708090a0b0c0d0f02030405060708090a0b0c0d0f02030405060708090a0b0c0d0f02030405060708090a0b0c0d0f")
	contract.SetCallCode(&addr, sdk.Hash{}, code)

	// test push2
	pc := uint64(0)
	makePush(2, 2)(&pc, interpreter, contract, mem, stack)

	require.Equal(t, pc, uint64(2))
	v := stack.pop()
	require.True(t, v.Cmp(big.NewInt(0).SetBytes(code[1:3])) == 0)

	// test push3
	pc = uint64(0)
	makePush(3, 3)(&pc, interpreter, contract, mem, stack)

	require.Equal(t, pc, uint64(3))
	v = stack.pop()
	require.True(t, v.Cmp(big.NewInt(0).SetBytes(code[1:4])) == 0)

	// test push32
	pc = uint64(0)
	makePush(32, 32)(&pc, interpreter, contract, mem, stack)

	require.Equal(t, pc, uint64(32))
	v = stack.pop()
	require.True(t, v.Cmp(big.NewInt(0).SetBytes(code[1:33])) == 0)

	// test push32
	code, _ = hexutil.Decode("0102030405060708090a0b0c0d0f")
	contract.SetCallCode(&addr, sdk.Hash{}, code)

	pc = uint64(0)
	makePush(32, 32)(&pc, interpreter, contract, mem, stack)

	require.Equal(t, pc, uint64(32))
	v = stack.pop()
	require.True(t, v.Cmp(big.NewInt(0).SetBytes(common.RightPadBytes(code[1:], 32))) == 0)
}

func TestOpPush2ToOpPush32(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()

	code, _ := hexutil.Decode("0102030405060708090a0b0c0d0f02030405060708090a0b0c0d0f020304050607")
	contract.SetCallCode(&addr, sdk.Hash{}, code)

	type testCase struct {
		pushN         int
		expectedValue *big.Int
	}

	var testCases []testCase
	for i := 2; i < 33; i++ {
		testCases = append(testCases, testCase{i, big.NewInt(0).SetBytes(common.RightPadBytes(code[1:1+i], i))})
	}

	for _, tc := range testCases {
		pc := uint64(0)
		makePush(uint64(tc.pushN), tc.pushN)(&pc, interpreter, contract, mem, stack)

		require.Equal(t, pc, uint64(tc.pushN))
		v := stack.pop()
		require.True(t, v.Cmp(tc.expectedValue) == 0)
	}
}

func TestOpDup1ToOpDup16(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	for i := 16; i != 0; i-- {
		stack.push(big.NewInt(int64(i)))
	}

	for i := 1; i < 17; i++ {
		makeDup(int64(i))(&pc, interpreter, contract, nil, stack)

		v := stack.pop().Int64()
		require.Equal(t, v, int64(i))
	}
}

func TestOpSwap1ToOpSwap16(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	for i := 16; i >= 0; i-- {
		stack.push(big.NewInt(int64(i)))
	}

	for i := 16; i != 0; i-- {
		makeSwap(int64(i))(&pc, interpreter, contract, nil, stack)

		v := stack.peek().Int64()
		require.Equal(t, v, int64(i))
	}
}

func TestOpLog0(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	logData := []byte("ab")
	blockNumber := big.NewInt(100)
	mem.Resize(64)
	mem.Set(0, 2, logData)
	interpreter.intPool = poolOfIntPools.get()
	interpreter.evm.BlockNumber = blockNumber
	pc := uint64(0)

	expectedLog := Log{
		Address:     addr,
		Topics:      make([]sdk.Hash, 0),
		Data:        logData,
		BlockNumber: blockNumber.Uint64(),
	}

	stack.push(big.NewInt(2))
	stack.push(big.NewInt(0))
	makeLog(0)(&pc, interpreter, contract, mem, stack)
	logs := interpreter.evm.StateDB.Logs()
	require.True(t, len(logs) == 1)
	require.True(t, reflect.DeepEqual(*logs[0], expectedLog))
}

func TestOpLog1(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	logData := []byte("ab")
	blockNumber := big.NewInt(100)
	mem.Resize(64)
	mem.Set(0, 2, logData)
	interpreter.intPool = poolOfIntPools.get()
	interpreter.evm.BlockNumber = blockNumber
	pc := uint64(0)

	expectedLog := Log{
		Address:     addr,
		Topics:      make([]sdk.Hash, 0),
		Data:        logData,
		BlockNumber: blockNumber.Uint64(),
	}
	topic1 := sdk.BytesToHash([]byte("1"))
	expectedLog.Topics = append(expectedLog.Topics, topic1)

	stack.push(big.NewInt(0).SetBytes(topic1.Bytes()))
	stack.push(big.NewInt(2))
	stack.push(big.NewInt(0))
	makeLog(1)(&pc, interpreter, contract, mem, stack)
	logs := interpreter.evm.StateDB.Logs()
	require.True(t, len(logs) == 1)
	require.True(t, reflect.DeepEqual(*logs[0], expectedLog))
}

func TestOpLog2(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	logData := []byte("ab")
	blockNumber := big.NewInt(100)
	mem.Resize(64)
	mem.Set(0, 2, logData)
	interpreter.intPool = poolOfIntPools.get()
	interpreter.evm.BlockNumber = blockNumber
	pc := uint64(0)

	expectedLog := Log{
		Address:     addr,
		Topics:      make([]sdk.Hash, 0),
		Data:        logData,
		BlockNumber: blockNumber.Uint64(),
	}

	topic1 := sdk.BytesToHash([]byte("1"))
	topic2 := sdk.BytesToHash([]byte("2"))
	expectedLog.Topics = append(expectedLog.Topics, topic1)
	expectedLog.Topics = append(expectedLog.Topics, topic2)

	stack.push(big.NewInt(0).SetBytes(topic2.Bytes()))
	stack.push(big.NewInt(0).SetBytes(topic1.Bytes()))
	stack.push(big.NewInt(2))
	stack.push(big.NewInt(0))
	makeLog(2)(&pc, interpreter, contract, mem, stack)
	logs := interpreter.evm.StateDB.Logs()
	require.True(t, len(logs) == 1)
	require.True(t, reflect.DeepEqual(*logs[0], expectedLog))
}

func TestOpLog3(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	logData := []byte("ab")
	blockNumber := big.NewInt(100)
	mem.Resize(64)
	mem.Set(0, 2, logData)
	interpreter.intPool = poolOfIntPools.get()
	interpreter.evm.BlockNumber = blockNumber
	pc := uint64(0)

	expectedLog := Log{
		Address:     addr,
		Topics:      make([]sdk.Hash, 0),
		Data:        logData,
		BlockNumber: blockNumber.Uint64(),
	}

	topic1 := sdk.BytesToHash([]byte("1"))
	topic2 := sdk.BytesToHash([]byte("2"))
	topic3 := sdk.BytesToHash([]byte("3"))
	expectedLog.Topics = append(expectedLog.Topics, topic1)
	expectedLog.Topics = append(expectedLog.Topics, topic2)
	expectedLog.Topics = append(expectedLog.Topics, topic3)

	stack.push(big.NewInt(0).SetBytes(topic3.Bytes()))
	stack.push(big.NewInt(0).SetBytes(topic2.Bytes()))
	stack.push(big.NewInt(0).SetBytes(topic1.Bytes()))
	stack.push(big.NewInt(2))
	stack.push(big.NewInt(0))
	makeLog(3)(&pc, interpreter, contract, mem, stack)
	logs := interpreter.evm.StateDB.Logs()
	require.True(t, len(logs) == 1)
	require.True(t, reflect.DeepEqual(*logs[0], expectedLog))
}

func TestOpLog4(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	logData := []byte("ab")
	blockNumber := big.NewInt(100)
	mem.Resize(64)
	mem.Set(0, 2, logData)
	interpreter.intPool = poolOfIntPools.get()
	interpreter.evm.BlockNumber = blockNumber
	pc := uint64(0)

	expectedLog := Log{
		Address:     addr,
		Topics:      make([]sdk.Hash, 0),
		Data:        logData,
		BlockNumber: blockNumber.Uint64(),
	}

	topic1 := sdk.BytesToHash([]byte("1"))
	topic2 := sdk.BytesToHash([]byte("2"))
	topic3 := sdk.BytesToHash([]byte("3"))
	topic4 := sdk.BytesToHash([]byte("4"))
	expectedLog.Topics = append(expectedLog.Topics, topic1)
	expectedLog.Topics = append(expectedLog.Topics, topic2)
	expectedLog.Topics = append(expectedLog.Topics, topic3)
	expectedLog.Topics = append(expectedLog.Topics, topic4)

	stack.push(big.NewInt(0).SetBytes(topic4.Bytes()))
	stack.push(big.NewInt(0).SetBytes(topic3.Bytes()))
	stack.push(big.NewInt(0).SetBytes(topic2.Bytes()))
	stack.push(big.NewInt(0).SetBytes(topic1.Bytes()))
	stack.push(big.NewInt(2))
	stack.push(big.NewInt(0))
	makeLog(4)(&pc, interpreter, contract, mem, stack)
	logs := interpreter.evm.StateDB.Logs()
	require.True(t, len(logs) == 1)
	require.True(t, reflect.DeepEqual(*logs[0], expectedLog))
}

func TestOpCallDataLoad(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	inputHexData := "0000000000000000000000000000000000000000000000000000000000000020"
	inputData, _ := hexutil.Decode(inputHexData)
	contract.Input = inputData

	stack.push(big.NewInt(0))
	opCallDataLoad(&pc, interpreter, contract, mem, stack)

	v := fmt.Sprintf("%064x", stack.pop().Uint64())
	require.True(t, reflect.DeepEqual(inputHexData, v))
}

func TestOpCallDataSize(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	inputHexData := "0000000000000000000000000000000000000000000000000000000000000020"
	inputData, _ := hexutil.Decode(inputHexData)
	contract.Input = inputData

	opCallDataSize(&pc, interpreter, contract, mem, stack)

	v := stack.pop().Uint64()
	require.True(t, v == uint64(len(inputData)))
}

func TestOpCallDataCopy(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	inputHexData := "0000000000000000000000000000000000000000000000000000000000000020"
	inputData, _ := hexutil.Decode(inputHexData)
	contract.Input = inputData

	mem.Resize(10)

	stack.push(big.NewInt(10)) // data len
	stack.push(big.NewInt(22)) // data offset in input
	stack.push(big.NewInt(0))  // data offset in memory
	opCallDataCopy(&pc, interpreter, contract, mem, stack)

	require.True(t, reflect.DeepEqual(mem.Data(), inputData[22:]))
}

func TestOpAddmod(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	stack.push(big.NewInt(3))
	stack.push(big.NewInt(12))
	stack.push(big.NewInt(1))
	opAddmod(&pc, interpreter, contract, nil, stack)

	v := stack.pop()
	expected := big.NewInt(1)
	require.True(t, v.Cmp(expected) == 0)

	//
	stack.push(big.NewInt(8))
	stack.push(big.NewInt(96))
	stack.push(big.NewInt(4))
	opAddmod(&pc, interpreter, contract, nil, stack)

	v = stack.pop()
	expected = big.NewInt(4)
	require.True(t, v.Cmp(expected) == 0)

	//
	stack.push(big.NewInt(-3))
	stack.push(big.NewInt(12))
	stack.push(big.NewInt(1))
	opAddmod(&pc, interpreter, contract, nil, stack)

	v = stack.pop()
	expected = big.NewInt(0)
	require.True(t, v.Cmp(expected) == 0)
}

func TestOpMulmod(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	stack.push(big.NewInt(3))
	stack.push(big.NewInt(11))
	stack.push(big.NewInt(1))
	opMulmod(&pc, interpreter, contract, nil, stack)

	v := stack.pop()
	expected := big.NewInt(2)
	require.True(t, v.Cmp(expected) == 0)

	//
	stack.push(big.NewInt(8))
	stack.push(big.NewInt(7))
	stack.push(big.NewInt(3))
	opMulmod(&pc, interpreter, contract, nil, stack)

	v = stack.pop()
	expected = big.NewInt(5)
	require.True(t, v.Cmp(expected) == 0)

	//
	stack.push(big.NewInt(8))
	stack.push(big.NewInt(8))
	stack.push(big.NewInt(3))
	opMulmod(&pc, interpreter, contract, nil, stack)

	v = stack.pop()
	expected = big.NewInt(0)
	require.True(t, v.Cmp(expected) == 0)

	//
	stack.push(big.NewInt(-3))
	stack.push(big.NewInt(11))
	stack.push(big.NewInt(1))
	opMulmod(&pc, interpreter, contract, nil, stack)

	v = stack.pop()
	expected = big.NewInt(0)
	require.True(t, v.Cmp(expected) == 0)
}

func TestOpReturnDataSize(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	interpreter.SetReturnData([]byte("abc"))
	pc := uint64(0)

	opReturnDataSize(&pc, interpreter, contract, nil, stack)

	v := stack.pop()

	require.True(t, v.Cmp(big.NewInt(3)) == 0)
}

func TestOpReturnDataCopy(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	mem.Resize(64)

	interpreter.intPool = poolOfIntPools.get()
	interpreter.SetReturnData([]byte("abc"))
	pc := uint64(0)

	stack.push(big.NewInt(3))
	stack.push(big.NewInt(0))
	stack.push(big.NewInt(10))
	opReturnDataCopy(&pc, interpreter, contract, mem, stack)

	require.True(t, reflect.DeepEqual(mem.Data()[10:13], []byte("abc")))

	//
	stack.push(big.NewInt(4))
	stack.push(big.NewInt(0))
	stack.push(big.NewInt(10))
	_, err := opReturnDataCopy(&pc, interpreter, contract, mem, stack)

	require.Equal(t, ErrReturnDataOutOfBounds.Error(), err.Error())
}

func TestOpTimestamp(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	now := time.Now().Unix()
	interpreter.evm.Time = big.NewInt(now)
	pc := uint64(0)

	opTimestamp(&pc, interpreter, contract, nil, stack)

	v := stack.pop()

	require.Equal(t, now, v.Int64())
}

func TestOpNumber(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	blockNumber := int64(100)
	interpreter.evm.BlockNumber = big.NewInt(blockNumber)
	pc := uint64(0)

	opNumber(&pc, interpreter, contract, nil, stack)

	v := stack.pop()

	require.Equal(t, blockNumber, v.Int64())
}

func TestOpDifficulty(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	opDifficulty(&pc, interpreter, contract, nil, stack)

	v := stack.pop()

	require.Equal(t, int64(1), v.Int64())
}

func TestOpChainID(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	interpreter.evm.chainConfig.ChainID = "test"
	pc := uint64(0)

	opChainID(&pc, interpreter, contract, nil, stack)

	v := stack.pop()
	t.Log(string(v.Bytes()))

	require.True(t, string(v.Bytes()) == interpreter.evm.chainConfig.ChainID)
}

func TestOpPop(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	stack.push(big.NewInt(0))
	require.True(t, stack.len() == 1)
	opPop(&pc, interpreter, contract, nil, stack)

	require.True(t, stack.len() == 0)
}

func TestOpMload(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	offset := uint64(10)
	d := big.NewInt(100)

	mem.Resize(64)
	mem.Set32(offset, d)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	stack.push(big.NewInt(int64(offset)))
	opMload(&pc, interpreter, contract, mem, stack)

	v := stack.pop()

	require.True(t, v.Cmp(d) == 0)
}

func TestOpMstore8(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	offset := uint64(10)
	hexData, _ := hexutil.Decode("0xdededede")
	d := big.NewInt(0).SetBytes(hexData)

	mem.Resize(64)

	interpreter.intPool = poolOfIntPools.get()
	pc := uint64(0)

	stack.push(d)
	stack.push(big.NewInt(int64(offset)))
	opMstore8(&pc, interpreter, contract, mem, stack)

	memD := mem.Data()[offset : offset+4]
	require.True(t, memD[0] == hexData[0] && memD[1] == memD[2] && memD[2] == memD[3] && memD[3] == 0)
}

func TestOpJump(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	code, _ := hexutil.Decode("608060405261271060005534801561001657600080fd5b506040516107ec3803806107ec8339818101604052602081101561003957600080fd5b8101908080519060200190929190505050610059816100ad60201b60201c565b61005f57fe5b33600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600181905550506100bb565b600080548211159050919050565b610722806100ca6000396000f3fe60806040526004361061007b5760003560e01c80639c75ed6c1161004e5780639c75ed6c1461019c578063d9c90cff146101c7578063de8c50c81461021a578063e662bd25146102455761007b565b80632e1a7d4d146100805780633a6e3d98146100bb57806345596e2e1461010a5780638da5cb5b14610145575b600080fd5b34801561008c57600080fd5b506100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610289565b005b3480156100c757600080fd5b506100f4600480360360208110156100de57600080fd5b810190808035906020019092919050505061034c565b6040518082815260200191505060405180910390f35b34801561011657600080fd5b506101436004803603602081101561012d57600080fd5b810190808035906020019092919050505061037e565b005b34801561015157600080fd5b5061015a6103ee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101a857600080fd5b506101b1610414565b6040518082815260200191505060405180910390f35b3480156101d357600080fd5b50610200600480360360208110156101ea57600080fd5b810190808035906020019092919050505061041a565b604051808215151515815260200191505060405180910390f35b34801561022657600080fd5b5061022f610428565b6040518082815260200191505060405180910390f35b6102876004803603602081101561025b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061042e565b005b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102e057fe5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610348573d6000803e3d6000fd5b5050565b60006103776000546103696001548561053590919063ffffffff16565b6105bb90919063ffffffff16565b9050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103d557fe5b6103de8161041a565b6103e457fe5b8060018190555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b600080548211159050919050565b60005481565b60006104393461034c565b9050600081340390508273ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610488573d6000803e3d6000fd5b507f9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc4433843484604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200194505050505060405180910390a1505050565b60008083141561054857600090506105b5565b600082840290508284828161055957fe5b04146105b0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806106cc6021913960400191505060405180910390fd5b809150505b92915050565b60006105fd83836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f000000000000815250610605565b905092915050565b600080831182906106b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561067657808201518184015260208101905061065b565b50505050905090810190601f1680156106a35780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385816106bd57fe5b04905080915050939250505056fe536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a26469706673582212204222d9732198380684d25bb3c9e975cd4ae4a7664c6d2a4d81ac457a590d9d7e64736f6c63430006000033")
	contract.SetCallCode(&addr, sdk.Hash{}, code)
	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	stack.push(big.NewInt(22)) // code[22] == '6b': 91 is JUMPDEST
	_, err := opJump(&pc, interpreter, contract, nil, stack)
	require.Nil(t, err)
	require.Equal(t, uint64(22), pc)
}

func TestOpJumpi(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	code, _ := hexutil.Decode("608060405261271060005534801561001657600080fd5b506040516107ec3803806107ec8339818101604052602081101561003957600080fd5b8101908080519060200190929190505050610059816100ad60201b60201c565b61005f57fe5b33600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600181905550506100bb565b600080548211159050919050565b610722806100ca6000396000f3fe60806040526004361061007b5760003560e01c80639c75ed6c1161004e5780639c75ed6c1461019c578063d9c90cff146101c7578063de8c50c81461021a578063e662bd25146102455761007b565b80632e1a7d4d146100805780633a6e3d98146100bb57806345596e2e1461010a5780638da5cb5b14610145575b600080fd5b34801561008c57600080fd5b506100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610289565b005b3480156100c757600080fd5b506100f4600480360360208110156100de57600080fd5b810190808035906020019092919050505061034c565b6040518082815260200191505060405180910390f35b34801561011657600080fd5b506101436004803603602081101561012d57600080fd5b810190808035906020019092919050505061037e565b005b34801561015157600080fd5b5061015a6103ee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101a857600080fd5b506101b1610414565b6040518082815260200191505060405180910390f35b3480156101d357600080fd5b50610200600480360360208110156101ea57600080fd5b810190808035906020019092919050505061041a565b604051808215151515815260200191505060405180910390f35b34801561022657600080fd5b5061022f610428565b6040518082815260200191505060405180910390f35b6102876004803603602081101561025b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061042e565b005b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102e057fe5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610348573d6000803e3d6000fd5b5050565b60006103776000546103696001548561053590919063ffffffff16565b6105bb90919063ffffffff16565b9050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103d557fe5b6103de8161041a565b6103e457fe5b8060018190555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b600080548211159050919050565b60005481565b60006104393461034c565b9050600081340390508273ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610488573d6000803e3d6000fd5b507f9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc4433843484604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200194505050505060405180910390a1505050565b60008083141561054857600090506105b5565b600082840290508284828161055957fe5b04146105b0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806106cc6021913960400191505060405180910390fd5b809150505b92915050565b60006105fd83836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f000000000000815250610605565b905092915050565b600080831182906106b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561067657808201518184015260208101905061065b565b50505050905090810190601f1680156106a35780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385816106bd57fe5b04905080915050939250505056fe536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a26469706673582212204222d9732198380684d25bb3c9e975cd4ae4a7664c6d2a4d81ac457a590d9d7e64736f6c63430006000033")
	contract.SetCallCode(&addr, sdk.Hash{}, code)
	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	stack.push(big.NewInt(0))  // cond == 0 can not jump, and the pc will puls 1
	stack.push(big.NewInt(22)) // code[22] == '6b': 91 is JUMPDEST
	_, err := opJumpi(&pc, interpreter, contract, nil, stack)
	require.Equal(t, uint64(1), pc)

	stack.push(big.NewInt(1))  // cond != 0 can jump
	stack.push(big.NewInt(22)) // code[22] == '6b': 91 is JUMPDEST
	_, err = opJumpi(&pc, interpreter, contract, nil, stack)
	require.Nil(t, err)
	require.Equal(t, uint64(22), pc)
}

func TestOpPc(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	pc := uint64(100)
	interpreter.intPool = poolOfIntPools.get()

	opPc(&pc, interpreter, contract, nil, stack)

	v := stack.pop().Uint64()
	require.Equal(t, v, pc)
}

func TestOpMsize(t *testing.T) {
	var (
		addr        = sdk.AccAddress{0xab}
		value       = big.NewInt(1000)
		env         = newEVM()
		stack       = newstack()
		mem         = NewMemory()
		interpreter = NewEVMInterpreter(env, env.vmConfig)
		contract    = NewContract(&dummyContractRef{address: addr}, &dummyContractRef{address: addr}, value, 0)
	)

	mem.Resize(64)
	pc := uint64(0)
	interpreter.intPool = poolOfIntPools.get()

	opMsize(&pc, interpreter, contract, mem, stack)

	v := stack.pop().Uint64()
	require.Equal(t, v, uint64(64))
}
