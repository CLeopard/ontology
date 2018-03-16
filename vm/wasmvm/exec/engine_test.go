package exec

import (
	"encoding/binary"
	"fmt"
	"github.com/Ontology/common"
	"io/ioutil"
	"testing"
	"math"
	"github.com/Ontology/vm/wasmvm/util"

)

var service = NewInteropService()

func TestAdd(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)

	code, err := ioutil.ReadFile("./testdata2/math.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	fmt.Printf("code bytes:%v\n", code)
	codestring := common.ToHexString(code)

	fmt.Println(codestring)

	b, _ := common.HexToBytes(codestring)
	fmt.Println(b)
	method2 := "add"
	input2 := make([]byte, 9)
	input2[0] = byte(len(method2))
	copy(input2[1:len(method2)+1], []byte(method2))
	input2[len(method2)+1] = byte(2) //param count
	input2[len(method2)+2] = byte(1) //param1 length
	input2[len(method2)+3] = byte(1) //param2 length
	input2[len(method2)+4] = byte(5) //param1
	input2[len(method2)+5] = byte(9) //param2

	fmt.Println(input2)
	res2, err := engine.Call(common.Uint160{}, code, input2)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res2)
	if binary.LittleEndian.Uint32(res2) != uint32(14) {
		t.Error("the result should be 14")
	}

}

func TestSquare(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)

	code, err := ioutil.ReadFile("./testdata2/math.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "square"

	input := make([]byte, 10)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(1) //param count
	input[len(method)+2] = byte(1) //param1 length
	input[len(method)+3] = byte(5) //param1

	fmt.Println(input)
	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	if binary.LittleEndian.Uint32(res) != uint32(25) {
		t.Error("the result should be 25")
	}

}

func TestEnvAddTwo(t *testing.T) {

	service.Register("addOne", func(engine *ExecutionEngine) (bool, error) {
		fmt.Println(engine)
		param:= engine.vm.envCall.envParams[0]
		engine.vm.ctx = engine.vm.envCall.envPreCtx
		engine.vm.pushUint64(param+1)
		return true,nil
	})

	engine := NewExecutionEngine(nil,nil,nil,service)

	code, err := ioutil.ReadFile("./testdata2/testenv.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "addTwo"

	input := make([]byte, 8)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(0)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	if binary.LittleEndian.Uint32(res) != uint32(3) {
		t.Error("the result should be 3")
	}
}

func TestBlockHeight(t *testing.T) {

	service.Register("getBlockHeight", func(engine *ExecutionEngine) (bool, error) {

		engine.vm.ctx = engine.vm.envCall.envPreCtx
		engine.vm.pushUint64(uint64(25535))
		return true,nil
	})


	engine := NewExecutionEngine(nil,nil,nil,service)

	code, err := ioutil.ReadFile("./testdata2/testBlockHeight.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "blockHeight"

	input := make([]byte, 13)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(0)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	if binary.LittleEndian.Uint32(res) != uint32(25535) {
		t.Error("the result should be 25535")
	}
}

func TestMem(t *testing.T) {

	service.Register("getString", func(engine *ExecutionEngine) (bool, error) {

		mem := engine.vm.memory.Memory
		param := engine.vm.envCall.envParams
		start := int(param[0])
		length := int(param[1])
		fmt.Printf("start is %d,length is %d\n", start, length)
		str := string(mem[start : start+length])
		engine.vm.ctx = engine.vm.envCall.envPreCtx
		engine.vm.pushUint64(uint64(len(str)))
		return true,nil

	})

	engine := NewExecutionEngine(nil,nil,nil,service)
	//test
	code, err := ioutil.ReadFile("./testdata2/TestMemory.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "getStrLen"

	input := make([]byte, 11)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(0)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(engine.vm.memory)
	if binary.LittleEndian.Uint32(res) != uint32(11) {
		t.Error("the result should be 11")
	}
}

func TestGlobal(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/str.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "_getStr"

	input := make([]byte, 9)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(0)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(engine.vm.memory)

}


func TestIf(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/ifTest.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "testif"
	p := 5

	input := make([]byte, 100)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(1)
	input[len(method)+2] = byte(1)
	input[len(method)+3] = byte(p)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	//fmt.Println(engine.memory)

	var shouldR int
	if p < 5 {
		shouldR = 10 + p
	} else {
		shouldR = 20 + p
	}

	if binary.LittleEndian.Uint32(res) != uint32(shouldR) {
		t.Fatal("result should be", shouldR)
	}

}

func TestLoop(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/ifTest.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "testfor"
	p := 10

	input := make([]byte, 100)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(1)
	input[len(method)+2] = byte(1)
	input[len(method)+3] = byte(p)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	//fmt.Println(engine.memory)

	if binary.LittleEndian.Uint32(res) != uint32(2*(p+1)) {
		t.Fatal("result should be", 2*(p+1))
	}

}

func TestWhileLoop(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/ifTest.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "testwhile"
	p := 10

	input := make([]byte, 100)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(1)
	input[len(method)+2] = byte(1)
	input[len(method)+3] = byte(p)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	//fmt.Println(engine.memory)

	if binary.LittleEndian.Uint32(res) != uint32(p+10) {
		t.Fatal("result should be", p+10)
	}

}

func TestIfII(t *testing.T) {

	s := "b456c4862902525e17ace6a2607f0806f51df0a98c3629c27f00efcf87ee8784"
	fmt.Println([]byte(s))
	u := binary.LittleEndian.Uint64([]byte(s))

	b := make([]byte, 64)
	binary.LittleEndian.PutUint64(b, u)
	fmt.Println(b)

	fmt.Println(u)
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/ifTest.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "testifII"
	p := 10

	input := make([]byte, 100)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(1)
	input[len(method)+2] = byte(1)
	input[len(method)+3] = byte(p)

	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	//fmt.Println(engine.memory)

	if binary.LittleEndian.Uint32(res) != uint32(60) {
		t.Fatal("result should be", 60)
	}

}

func TestUI256Addr(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/ui256Test.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "getAddress"

	input := make([]byte, 100)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(4)
	input[len(method)+2] = byte(8)
	input[len(method)+3] = byte(8)
	input[len(method)+4] = byte(8)
	input[len(method)+5] = byte(8)

	s := "abcdefghijklmnopqrstuvwxyz123456"
	bytes := []byte(s)
	u256, _ := common.Uint256ParseFromBytes(bytes)

	fmt.Println(u256.ToString())

	copy(input[len(method)+6:], bytes)

	fmt.Printf("bytes is %v\n", bytes)
	//fmt.Printf("before call ,vm.mem len is %d\n", len(engine.vm.memory.Memory))
	res, err := engine.Call(common.Uint160{}, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(engine.vm.memory)
	addr, _ := common.Uint256ParseFromBytes(engine.vm.memory.Memory[:32])
	fmt.Printf("after bytes:%v\n", engine.vm.memory.Memory[:32])
	fmt.Println("afterf:", addr.ToString())
	if addr.ToString() != u256.ToString() {
		t.Fatal("the address before is not same as after!")
	}

	if binary.LittleEndian.Uint32(res) != uint32(0) {
		t.Error("the result should be 0")
	}
}


func TestStrings(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/strings.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	input := make([]interface{}, 2)
	input[0] = "getStringlen"
	input[1] = "abcdefghijklmnopqrstuvwxyz"
	//input[2] = 3

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(string(res))
	if binary.LittleEndian.Uint32(res) != 26{
		t.Fatal("the res should be 26")
	}
	//fmt.Println(engine.memory)

}

func TestIntArraySum(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/intarray.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}


	input := make([]interface{}, 3)
	input[0] = "_sum"
	input[1] = []int{1, 2, 3, 4}
	input[2] = 4

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	idx := int(binary.LittleEndian.Uint32(res))
	fmt.Println(engine.GetMemory().Memory[idx:idx+10])
	if binary.LittleEndian.Uint32(res) != 10 {
		t.Fatal("the res should be 10")
	}
	//fmt.Println(engine.memory)

}


func TestSimplestruct(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/simplestruct.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	type student struct{
		name string
		math int
		eng int64

	}

	s := student{name:"jack",math:90,eng:95}

	input := make([]interface{}, 2)
	input[0] = "getSum"
	input[1] = s

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(engine.vm.memory.Memory[:20])
	if binary.LittleEndian.Uint32(res) != 185 {
		t.Fatal("the res should be 185")
	}


}


func TestSimplestruct2(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/simplestruct.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	type student struct{
		name string
		math int
		eng int64

	}

	s := student{name:"jack",math:90,eng:95}

	input := make([]interface{}, 2)
	input[0] = "getName"
	input[1] = s

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	idx:=int(binary.LittleEndian.Uint32(res))
	fmt.Println(engine.vm.memory.Memory[idx:idx+20])
	var length int
	for i := idx;engine.vm.memory.Memory[i]!= 0;i++{
		length += 1
	}
	if string(engine.vm.memory.Memory[idx:idx+length]) != "jack"{
		t.Fatal("the res should be jack")
	}
/*	if binary.LittleEndian.Uint32(res) != 185 {
		t.Fatal("the res should be 185")
	}*/


}

func TestComplexstruct(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/complexStruct.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	type acct struct{
		Acct1 []int
		Acct2 []int

	}

	s := acct{Acct1:[]int{10,20,30,40}}

	input := make([]interface{}, 3)
	input[0] = "sumAcct1"
	input[1] = s
	input[2] = 4

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(engine.vm.memory.Memory[:20])
	tmp:= binary.LittleEndian.Uint32(engine.vm.memory.Memory[0:4])
	fmt.Println(engine.vm.memory.Memory[tmp:tmp+20])
	if binary.LittleEndian.Uint32(res) != 100 {
		t.Fatal("the res should be 100")
	}
}

func TestFloatSum(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/float.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}


	input := make([]interface{}, 3)
	input[0] = "sum"
	input[1] = float32(1.1)
	input[2] = float32(0.5)

	fmt.Printf("input is %v\n", input)


	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(util.ByteToFloat32(res))
	fmt.Println(math.Float32frombits(binary.LittleEndian.Uint32(res)))
	if math.Float32frombits(binary.LittleEndian.Uint32(res))  != 1.6 {
		t.Fatal("the res should be  1.6 ")
	}
	//fmt.Println(engine.memory)

}
func TestDoubleSum(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/float.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}


	input := make([]interface{}, 3)
	input[0] = "sumDouble"
	input[1] = 1.1
	input[2] = 0.5

	fmt.Printf("input is %v\n", input)


	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	fmt.Println(util.ByteToFloat64(res))
	fmt.Println(math.Float64frombits(binary.LittleEndian.Uint64(res)))
	if math.Float64frombits(binary.LittleEndian.Uint64(res))  != 1.6 {
		t.Fatal("the res should be  1.6 ")
	}
}


func TestCalloc(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/calloc.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}


	input := make([]interface{}, 1)
	input[0] = "retArray"


	fmt.Printf("input is %v\n", input)


	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	fmt.Println(engine.vm.memory.MemPoints)
	offset := binary.LittleEndian.Uint32(res)
	bytes,_ := engine.vm.memory.GetPointerMemory(uint64(offset))
	fmt.Println(bytes)

}

func TestMalloc(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/malloc.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}


	input := make([]interface{}, 4)
	input[0] = "initStu"
	input[1] = 100
	input[2] = 90
	input[3] = 85


	fmt.Printf("input is %v\n", input)


	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	fmt.Println(engine.vm.memory.MemPoints)
	offset := binary.LittleEndian.Uint32(res)
	bytes,_ := engine.vm.memory.GetPointerMemory(uint64(offset))
	fmt.Println(bytes)

}

//use 'arrayLen' instead of  'sizeof'
func TestArraylen(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/arraylen.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}


	input := make([]interface{}, 3)
	input[0] = "combine"
	input[1] = []int{1,2,3,4,5}
	input[2] = []int{6,7,8,9,10,11}



	fmt.Printf("input is %v\n", input)


	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	fmt.Println(engine.vm.memory.MemPoints)
	offset := binary.LittleEndian.Uint32(res)
	bytes,_ := engine.vm.memory.GetPointerMemory(uint64(offset))
	fmt.Println(bytes)
	cnt := len(bytes) / 4
	resarray := make([]int,cnt)
	for i:= 0 ;i < cnt;i++{
		resarray[i] = int(binary.LittleEndian.Uint32(bytes[i*4:(i+1)*4]))
	}
	fmt.Println(resarray)


}


func TestAddress(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/testGetAddress.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	input := make([]interface{}, 1)
	input[0] = "_getaddr"

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
}


func TestContract(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/contractTest.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	input := make([]interface{}, 3)
	input[0] = "apply"
	input[1] = 9999        //code ,address
	input[2] = 0			//action

	msg := make([]interface{},3)
	msg[0] = 9999
	msg[1] = 1000
	msg[2] = 50

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,msg)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	fmt.Println(engine.vm.memory.MemPoints)
	fmt.Println(engine.vm.memory.Memory[0:12])

}


func TestString(t *testing.T) {
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/stringtest.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	input := make([]interface{}, 2)
	input[0] = "greeting"
	input[1] = "may the force be with you"        //code ,address

	fmt.Printf("input is %v\n", input)


	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)


	fmt.Println(engine.vm.memory.MemPoints)
	offset := uint64(binary.LittleEndian.Uint32(res))
	length := engine.vm.memory.MemPoints[offset].Length
	fmt.Println(engine.vm.memory.Memory[offset:int(offset)+int(length)])
	fmt.Println(string(engine.vm.memory.Memory[offset:int(offset)+int(length)]))
	if input[1] != string(engine.vm.memory.Memory[offset:int(offset)+int(length)]){
		t.Fatal("the return should be :"+input[1].(string))
	}

}

func TestContractCall(t *testing.T){
	engine := NewExecutionEngine(nil,nil,nil,nil)
	//test
	code, err := ioutil.ReadFile("./testdata2/callContract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	input := make([]interface{}, 3)
	input[0] = "testCall"
	input[1] = 10
	input[2] = 29

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(common.Uint160{}, code, input,nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	if binary.LittleEndian.Uint32(res) != uint32(39){
		t.Fatal("the result should be 39!")
	}


}