package tesk

import "fmt"

// 无参数
func NoParamsMethod() {
	fmt.Println("无参数方法")
}

/*func ParamsMethod() {
	fmt.Println("单个参数")
}

// 多个参数 -一般参数为固定，模式
func MultipleParamsMethod() {
	fmt.Println("参数1")
	fmt.Println("参数2")
	fmt.Println("参数3")
	fmt.Println("参数4")
}*/

// 单个参数
func ParamsMethod(data string) {
	fmt.Println("单个参数", data)
}

// 多个参数 -一般参数为固定，模式
func MultipleParamsMethod(param1 string, param2 bool, param3 string, param4 string, param5 string) {
	fmt.Println("参数1", param1)
	fmt.Println("参数2", param2)
	fmt.Println("参数3", param3)
	fmt.Println("参数4", param4)
	fmt.Println("参数5", param5)
}
