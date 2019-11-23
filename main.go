package main

import (
	"fmt"
	"os"
)

func main() {
	// filedata := `
	// // 单行注释
	// /* 多行注释
	// 	abc
	// 	*b/
	// */
	// include("./test1.config")  //包含另一个本地配置，路径为当前相对路径
	// `
	// f, err := os.Open("test.config")
	// if err != nil {
	// 	fmt.Println("Open error: ", err)
	// 	return
	// }
	// defer f.Close()
	// buffer, err1 := ioutil.ReadAll(f)
	// fmt.Println("======>>>>>>", err1, buffer)
	fmt.Println("---------", os.Args)
	LoadPathConfig(os.Args[1])
	// rd := bufio.NewReader(f)
	// stacklist := make([]string, 10)
	// for {
	// 	line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
	// 	if err != nil || io.EOF == err {
	// 		break
	// 	}
	// 	zsIndex1 := strings.Index(line, "//")
	// 	if zsIndex1 <= 2 {
	// 		continue
	// 	}
	// 	zsIndex2 := strings.Index(line, "/*")
	// 	stacklen := len(stacklist)
	// 	if zsIndex2 >= 0 {
	// 		if stacklen <= 0 || stacklist[stacklen-1] != "*/" {

	// 		}
	// 	}
	// }
}
