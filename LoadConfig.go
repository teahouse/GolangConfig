package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 对一行配置进行解析，只有include或者key=value格式
func lcheckOneline(line string) (command string, targetval string) {
	index1 := strings.Index(line, "include")
	index2 := strings.Index(line, "=")
	if index1 >= 0 {
		command = "include"
		targetval = ""
		var isinstr bool = false
		for _, val := range line {
			toval := fmt.Sprintf("%c", val)
			if toval == "\"" {
				if !isinstr {
					isinstr = true
				} else {
					isinstr = false
					break
				}
				continue
			}
			if isinstr {
				targetval += toval
			}
		}
	} else if index2 >= 0 {
		var isinstr bool = false
		for _, val := range line {
			toval := fmt.Sprintf("%c", val)
			if isinstr {
				targetval += toval
			}
			if !isinstr {
				if toval != "=" {
					command += toval
				} else {
					if !isinstr {
						isinstr = true
					}
				}
			}
		}
		// vallen := len(targetval)
		// if vallen > 1 {
		// 	if strings.Index(targetval, "[[") >= 0 {
		// 		targetval = targetval[2 : vallen-2]
		// 	}
		// }
	} else {
		fmt.Printf("error: not find cmd=[%s]\n", line)
		os.Exit(2)
	}
	return
}

// 从key-value的Map里面获取对应的值
func lgetValueFromMap(keyvalMap map[string]string, tokey string, targetval string) string {
	val, ok := keyvalMap[tokey]
	if !ok {
		fmt.Printf("error: not find cmd=[%s] in [%s]\n", tokey, targetval)
		os.Exit(3)
	}
	if val[0:1] == "\"" {
		val = val[1 : len(val)-1]
	}
	return val
}

// 检测每个key-value的value值里面的替换关系
func lcheckValueReplace(keyvalMap map[string]string, targetval string) string {
	if len(targetval) <= 1 {
		return targetval
	}
	//valSplitList := strings.Split(targetval, "+")
	valSplitList := []string{}
	var isNumLine bool = false
	var preisNumLine bool = false
	var NumLineVal string = ""

	var targetvalEx string = ""
	var toSubStr string = ""
	for _, val := range targetval {
		toval := fmt.Sprintf("%c", val)
		targetvalEx += toval
		forstrlen := len(targetvalEx)
		lastTwoChar := ""
		if forstrlen > 1 {
			lastTwoChar = targetvalEx[forstrlen-2 : forstrlen]
		}
		for {
			if isNumLine {
				if lastTwoChar != "\\\"" && toval == "\"" {
					isNumLine = false
					preisNumLine = true
					break
				}
				NumLineVal += toval
			}
			if toval == "\"" {
				isNumLine = true
				NumLineVal = ""
				break
			}
			break
		}
		if preisNumLine {
			valSplitList = append(valSplitList, "\""+NumLineVal+"\"")
			continue
		}
		if toval == "+" {
			if len(toSubStr) > 0 {
				valSplitList = append(valSplitList, toSubStr)
			}
		} else {
			toSubStr += toval
		}
	}
	////fmt.Println("+++++++++++", valSplitList)
	if len(valSplitList) <= 0 {
		valSplitList = append(valSplitList, targetval)
	}
	var realvalstr string = ""
	for index := 0; index < len(valSplitList); index++ {
		val := valSplitList[index]
		////fmt.Println("========2222222222222==11==", val)
		if val[0:1] == "\"" {
			val = val[1 : len(val)-1]
		} else if val != "true" && val != "false" {
			val = lgetValueFromMap(keyvalMap, val, targetval)
		}
		////fmt.Println("========2222222222222==22==", val)
		realvalstr += val
	}
	return realvalstr
}

func analysisConfig(inputFilePath string, keyvalMap map[string]string) {
	//fmt.Println("=====analysisConfig===222=====", inputFilePath)
	f, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("Open error: %v", err)
		os.Exit(5)
		return
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("ReadAll error: %v", err)
		os.Exit(5)
		return
	}
	filedataStr := fmt.Sprintf("%s", buffer)
	fileLineList := []string{}
	var extendFileStr string = ""
	var oneLineCode string = ""
	var isNoteOneLine bool = false
	var isNoteNumLine bool = false

	var isNumLine bool = false
	//var preisNumLine bool = false
	var isLineJump bool = false
	var NumLineVal string = ""
	for i, val := range filedataStr {
		toval := fmt.Sprintf("%c", val) + ""
		extendFileStr += toval
		forstrlen := len(extendFileStr)
		if forstrlen <= 1 {
			continue
		}
		lastTwoChar := extendFileStr[forstrlen-2 : forstrlen]
		if isNoteNumLine {
			if lastTwoChar == "*/" {
				isNoteNumLine = false
			}
			continue
		} else if lastTwoChar == "*/" {
			fmt.Printf("error: unknow */ in %d", i)
			os.Exit(1)
		}
		if isNoteOneLine {
			if toval != "\n" {
				continue
			}
			isNoteOneLine = false
		}
		if lastTwoChar == "//" {
			isNoteOneLine = true
			oneLineCode = oneLineCode[:len(oneLineCode)-1]
			continue
		}
		if lastTwoChar == "/*" {
			isNoteNumLine = true
			oneLineCode = oneLineCode[:len(oneLineCode)-1]
			continue
		}
		/*
		 */
		for {
			isLineJump = false
			if isNumLine {
				if lastTwoChar == "\\\"" {
					NumLineVal += toval
					break
				}
				if toval == "\n" {
					if extendFileStr[len(extendFileStr)-3:len(extendFileStr)-2] != "\\" {
						fmt.Printf("error: string line end[%s] not find \\", extendFileStr)
						os.Exit(4)
					}
					oneLineCode = oneLineCode[0 : len(oneLineCode)-2]
					isLineJump = true
					break
				}
				if toval == "\"" {
					isNumLine = false
					//preisNumLine = true
					break
				}
				NumLineVal += toval
			}
			if toval == "\"" {
				isNumLine = true
				NumLineVal = ""
				break
			}
			break
		}
		if isLineJump {
			continue
		}
		if isNumLine {
			oneLineCode += toval
			continue
		}
		// if preisNumLine {
		// 	replaceCount++
		// 	repkey := replaceStr + fmt.Sprintf("%d", replaceCount)
		// 	keyvalMap[repkey] = NumLineVal
		// 	oneLineCode += repkey
		// }
		if toval == "\r" || toval == "\t" || toval == " " {
			continue
		}
		if toval == "\n" || toval == ";" {
			if len(oneLineCode) > 0 {
				////fmt.Println("=====1============", oneLineCode)
				fileLineList = append(fileLineList, oneLineCode)
			}
			oneLineCode = ""
		} else {
			////fmt.Println("=====22============", oneLineCode)
			oneLineCode += toval
		}
	}
	//fmt.Printf(">>>>%#v\n", fileLineList)
	paths, _ := filepath.Split(inputFilePath)
	// 解析每一行配置信息
	for index := 0; index < len(fileLineList); index++ {
		tline := fileLineList[index]
		tcmd, tvalue := lcheckOneline(tline)
		if tcmd == "include" { // 解析include嵌套包含配置文件
			tvalue = filepath.Clean(tvalue)
			if paths != "." && paths != "./" && paths != ".\\" {
				tvalue = paths + tvalue
			}
			//fmt.Println("=====analysisConfig========", paths, tvalue)
			analysisConfig(tvalue, keyvalMap)
			continue
		}
		keyvalMap[tcmd] = tvalue
	}
	//fmt.Println("-------------------------------------", keyvalMap)
}

// LoadPathConfig call load path config for Export
func LoadPathConfig(inputFilePath string) map[string]string {
	keyvalMap := map[string]string{}
	analysisConfig(inputFilePath, keyvalMap)
	retmap := map[string]string{}
	for key, val := range keyvalMap {
		tarval2 := lcheckValueReplace(keyvalMap, val)
		retmap[key] = tarval2
		fmt.Printf("%s=%s\n", key, tarval2)
	}
	return retmap
}
