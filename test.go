package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"unsafe"

	"./craw"
	"./excel"
)

type CCrawOpt struct {
	exp 		int8
}

// 定义结构体Person
type CPrice struct {
	//CCrawOpt
	Id         	string 	//宗地编号
	Location   	string 	//宗地位置
	Measure    	string 	//净用地面积
	InitPrice  	string 	//起始价
	FinalPrice 	string 	//成交价
	Price		string	//最终价格
	Time     	string	//成交时间
	Clincher   	string 	//竞得人
	PageId     	string 	//页面id
	Natrue		string
}

func InspectStruct(o interface{}) {
	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		address := "not-addressable"

		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}
		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()
		}
		if valueField.CanAddr() {
			address = fmt.Sprint(valueField.Addr().Pointer())
		}

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Address: %v\t, Field type: %v\t, Field kind: %v\n", typeField.Name,
			valueField.Interface(), address, typeField.Type, valueField.Kind())

		if valueField.Kind() == reflect.Struct {
			InspectStruct(valueField.Interface())
		}
	}
}


func Test_1(){
	var it interface{} =  &CPrice{Id:"sdfsdf",Location: "sfsf"}
	val := reflect.ValueOf(it)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	val.FieldByName("Location").Set(reflect.ValueOf("武侯"))
	fmt.Println(val.FieldByName("Location").String())
}

func Test_2() {
	it := CPrice{Id:"sdfsdf",Location: "sfsf"}
	val := reflect.ValueOf(&it).Elem()
	val.FieldByName("Location").Set(reflect.ValueOf("武侯"))
	fmt.Println(val.FieldByName("Location").String())
}

func Test_3(){
	var it interface{} =  CPrice{Id:"sdfsdf",Location: "sfsf"}
	//p := &it
	val := reflect.ValueOf(it.(interface{}))
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	val.FieldByName("Location").Set(reflect.ValueOf("武侯"))
	fmt.Println(val.FieldByName("Location").String())
}

func Test_4() {
	val := reflect.TypeOf(CPrice{})
	//val = val.Elem()

	cp := reflect.New(val)
	cp = cp.Elem()
	cp.FieldByName("Location").Set(reflect.ValueOf("武侯"))


	it := cp.Addr().Interface()
	var _it interface{} = &CPrice{}

	p := it.(*CPrice)

	fmt.Println(reflect.TypeOf(it).Kind())
	fmt.Println(reflect.TypeOf(_it).Kind())
	fmt.Println(reflect.TypeOf(p).Kind())
	fmt.Println(*p)
}

type Person struct {
	name   string
	age    int
	gender bool
}
func Test_5() {
	john := Person{"Johnssss", 30, true}
	pp := unsafe.Pointer(&john)                                                    // 结构体的起始地址
	pname := (*string)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(john.name)))   // 属性name的起始地址，转换为*string类型
	page := (*int)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(john.age)))        // 属性age的起始地址，转换为*int类型
	pgender := (*bool)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(john.gender))) // 属性gender的起始地址，转换为*bool类型

	fmt.Printf("%d/%d/%d",unsafe.Offsetof(john.name),unsafe.Offsetof(john.age),unsafe.Offsetof(john.gender))
	// 进行赋值
	*pname = "Alicesssss"
	*page = 28
	*pgender = false
	fmt.Println(john) // {Alice 28 false}
}

func Test_6()  {
	name := "12233"
	fmt.Printf("%d \n",unsafe.Pointer(&name))
	name = "12233"
	fmt.Printf("%d \n",unsafe.Pointer(&name))
	name = "45667"
	fmt.Printf("%d \n",unsafe.Pointer(&name))
	name = name + "dfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasddfsdfsdfsdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfadfasdfasd"
	fmt.Printf("%d \n",unsafe.Pointer(&name))
}

func Test_7(){
	m := make(map[string]string)
	mm := m
	m["name"] = "zxy"
	fmt.Println(mm)

	ar := []int{1,2}
	arr := ar
	ar[1] = 3
	ar = append(ar,3)
	fmt.Println(arr)
}

func Test_8(){
	reg := regexp.MustCompile(`^[\d]{4}-[\d]{2}`)
	fmt.Println(reg.MatchString("2017-02-03"))

	fmt.Println(reg.FindString("2017-02-03"))
	fmt.Println(reg.FindStringIndex("2017-02-03"),1)
}

func Test_9(){
	p := CPrice{}
	pp := p
	pp.Id = "122"
	fmt.Println(p)
	fmt.Println(pp)
}

//---------------------------------
type Base struct{
	num 	int
}

func (this *Base) Get() float32 {
	fmt.Printf("Base Get %d \n",this.num)
	return 0
}

type Base0 struct{
	ff 		int
	num 	int
}

func (this *Base0) Get() float32 {
	fmt.Printf("Base0 Get %d \n",this.num)
	return 0
}

type Baser interface {
	Get() float32
}

type TypeOne struct {
	Base0
	Base
	value float32
	//Base0
}

type TypeTwo struct {
	value float32
	*Base
}

type TypeThree struct {
	value float32
	Base
}

func (t *TypeOne) Get() float32 {
	fmt.Println("TypeOne Get")
	return t.value
}

func (t *TypeTwo) Get() float32 {
	fmt.Println("TypeTwo Get")
	return t.value * t.value
}

func (t *TypeThree) Get() float32 {
	fmt.Println("TypeThree Get")
	return t.value + t.value
}

func Test_10() {
	base := Base{1}
	base0 := Base0{2,5}
	t1 := &TypeOne{base0,base,1}
	t2 := &TypeTwo{2, &base}
	t3 := &TypeThree{4, base}

	//var b *Base
	b := (*Base)(unsafe.Pointer(t1))
	b.Get()

	b0 := (*Base0)(unsafe.Pointer(t1))
	b0.Get()

	//f := (*float32)(unsafe.Pointer(t1))
	//fmt.Printf("%d \n", f)

	b1 := (*Base)(unsafe.Pointer(t2))
	b1.Get()

	bases := []Baser{Baser(t1), Baser(t2), Baser(t3)}

	for s, _ := range bases {
		switch bases[s].(type) {
		case *TypeOne:
			fmt.Println("TypeOne")
		case *TypeTwo:
			fmt.Println("TypeTwo")
		case *TypeThree:
			fmt.Println("TypeThree")
		}

		fmt.Printf("The value is:  %f\n", bases[s].Get())
	}
}

type Parent struct {

}

func (this *Parent) Say(){
	fmt.Println("This is parent say~")
}

func (this *Parent)Call(method string,vals []reflect.Value){
	value := reflect.ValueOf(this)
	f := value.MethodByName(method)
	f.Call(vals)
}

type Child struct {
	Parent
}

func (this *Child) Say(){
	fmt.Println("This is child say~")
}

func Test_11(){
	c := &Child{}
	c.Call("Say",[]reflect.Value{})
	c.Say()
}

func Test_12(){
	var it interface{} = craw.NewCPriceCraw()
	if factor,ok := it.(excel.IExcelFactor);ok {
		fmt.Println(factor)
	}
}

type Car struct {
	model string
}
func (c Car) PrintModel() {
	fmt.Println(c.model)
}
func Test_13() {
	c := Car{model: "DeLorean DMC-12"}
	defer c.PrintModel()
	c.model = "Chevrolet Impala"
}

func Test_14(){
	buf := make([]byte,0,100)
	buf1 := buf[5:10]
	buf2 := buf[:11]
	buf1 = append(buf1,1)


	fmt.Println(buf1)
	fmt.Println(len(buf1))
	fmt.Println(cap(buf1))


	fmt.Println(buf2)
	fmt.Println(len(buf2))
	fmt.Println(cap(buf2))


	fmt.Println(string([]byte{ 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64}))

}

func Test_StringHeader(){
	bytes := []byte{ 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64}
	byteStr := (*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(&bytes[0])), Len: len(bytes)}))
	fmt.Println(*byteStr)

	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(byteStr))
	stringHeader.Len = 8
	fmt.Println(*byteStr)
}

func Test_String() {
	str1 := "1121"
	str2 := str1

	str1Header := (*reflect.StringHeader)(unsafe.Pointer(&str1))
	str2Header := (*reflect.StringHeader)(unsafe.Pointer(&str2))
	fmt.Printf("str1Header.Data = %v,str1Header.Len = %v \n",str1Header.Data,str1Header.Len)
	fmt.Printf("str2Header.Data = %v,str2Header.Len = %v \n",str2Header.Data,str2Header.Len)

	str1 = str1 + "sdfsd"

	fmt.Printf("str1Header.Data = %v,str1Header.Len = %v \n",str1Header.Data,str1Header.Len)
	fmt.Printf("str2Header.Data = %v,str2Header.Len = %v \n",str2Header.Data,str2Header.Len)


}

func Test_Slice(){
	bytes := make([]byte,2,3)
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))

	bytess := (*[]byte)(unsafe.Pointer(header))  //指针
	bytesEx := bytes

	fmt.Println(uintptr(unsafe.Pointer(&bytes)))
	fmt.Println(uintptr(unsafe.Pointer(&bytesEx)))

	fmt.Println(header.Data)

	fmt.Println("========================")
	bytes = append(bytes,2)
	fmt.Println(header.Data)

	bytes[0] = 1
	bytes[1] = 1

	fmt.Println(*bytess)
	fmt.Println(bytes)
	fmt.Println(bytesEx)

	fmt.Println("========================")

	bytes = append(bytes,2)
	fmt.Println(header.Data)

	bytes[0] = 2
	bytes[1] = 2

	fmt.Println(*bytess)
	fmt.Println(bytes)
	fmt.Println(bytesEx)

}

func Test_Array(){
	var byteArr = [3][3]byte{
		[3]byte{1,2,3},
		[3]byte{4,5,6},
		[3]byte{7,8,9},
	}

	fmt.Println(byteArr[0])
	fmt.Printf("byteArr[0]中保存的地址%p\n", &byteArr[0])
	fmt.Printf("byteArr[0][0]的地址%p\n", &byteArr[0][0])
	fmt.Printf("byteArr[1]的地址%p\n", &byteArr[1])
	fmt.Printf("byteArr[0][0]的地址%p\n", &byteArr[1][0])
	fmt.Printf("byteArr[2]的地址%p\n", &byteArr[2])
	fmt.Printf("byteArr[2][0]的地址%p\n", &byteArr[2][0])

	byteArray := (*[9]byte)(unsafe.Pointer(&byteArr))
	fmt.Printf("*byteArray[0] = %v\n", (*byteArray)[0])
	fmt.Printf("*byteArray[0] = %v\n", (*byteArray)[3])
	fmt.Printf("*byteArray[0] = %v\n", (*byteArray)[6])

	temp := byteArr
	temp[0][0] = 10
	fmt.Println(temp[0][0])
	fmt.Println(byteArr[0][0])

}

func Test_SliceHeader(){
	bytes := []byte{ 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64}
	array := (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len: len(bytes),
		Cap: cap(bytes),
	}))

	fmt.Println(*array)

	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(array))
	sliceHeader.Len = 5

	fmt.Println(*array)

}

func Test_SliceHeaderEx() {
	type T struct {
		a int
		b int
		c int
	}

	type SliceHeader struct {
		addr uintptr
		len  int
		cap  int
	}

	t := &T{a: 1,b: 2,c: 3}
	p := unsafe.Sizeof(*t)


	sl := &SliceHeader{
		addr: uintptr(unsafe.Pointer(t)),len:  int(p),cap:  int(p),}

	b := (*[]byte)(unsafe.Pointer(sl))

	fmt.Println(*t)
	fmt.Println(*b)


	(*b)[0] = 7
	(*b)[8] = 5
	(*b)[16] = 8

	fmt.Println(*t)

}

type T struct {
	a int
}

func (t T)SetA(a int){
	t.a = a
}
func (t *T)SetAX(a int){
	t.a = a
}
func Test_T(){
	t := T{1}
	t.SetA(2)
	fmt.Printf("%v \n", t)
	t.SetAX(2)
	fmt.Printf("%v \n", t)
}

func HttpGet_1(url string,path string,values url.Values) ([]byte,error){
	url = url + path
	if len(values) > 0 {
		url = url + "?" + values.Encode()
	}

	var res,err = http.Get(url)
	fmt.Println("HttpGet:"+url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body,err1 = ioutil.ReadAll(res.Body)
	if err1 != nil {
		return nil,err1
	}
	return body,nil
}

func HttpGet_2(url string) ([]byte,error){
	//url = url + path
	//if len(values) > 0 {
	//	url = url + "?" + values.Encode()
	//}

	var res,err = http.Get(url)
	fmt.Println("HttpGet:"+url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body,err1 = ioutil.ReadAll(res.Body)
	if err1 != nil {
		return nil,err1
	}
	return body,nil
}

func main() {
	//Test_1()
	//Test_2()
	//Test_3()
	//Test_4()
	//Test_5()
	//Test_6()
	//Test_7()
	//Test_8()
	//Test_9()
	//Test_10()

	//Test_11()
	//Test_12()
	//Test_13()
	//Test_14()

	//Test_StringHeader()
	//Test_SliceHeader()
	//Test_SliceHeaderEx()
	//Test_Slice()
	//Test_Array()
	//Test_String()
	//Test_T()
	HttpGet_2("http://127.0.0.1:9003/enter_room?userid=12&name=%E6%AC%A7%E9%98%B3%E6%9C%89%E9%92%B1&roomid=727284&sign=3872536de28976532f24497bc4177b63")

}




