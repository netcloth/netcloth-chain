package vm

import (
	"encoding/json"
	"fmt"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
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
