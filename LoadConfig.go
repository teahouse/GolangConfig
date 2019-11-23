package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
func lcheckValueReplace(keyvalMap map[string]string, targetval string) string {
	if len(targetval) <= 1 {
		return targetval
	}
	var realvalstr string = ""
	//allvaluelist := strings.Split(targetval, "+")
	allvaluelist := []string{}
	var issssstr bool = false
	var preissssstr bool = false
	var fsssstr string = ""

	var targetvalEx string = ""
	var togval string = ""
	for _, val := range targetval {
		toval := fmt.Sprintf("%c", val)
		targetvalEx += toval
		forstrlen := len(targetvalEx)
		laststr := ""
		if forstrlen > 1 {
			laststr = targetvalEx[forstrlen-2 : forstrlen]
		}
		for {
			if issssstr {
				if laststr != "\\\"" && toval == "\"" {
					issssstr = false
					preissssstr = true
					break
				}
				fsssstr += toval
			}
			if toval == "\"" {
				issssstr = true
				fsssstr = ""
				break
			}
			break
		}
		if preissssstr {
			allvaluelist = append(allvaluelist, "\""+fsssstr+"\"")
			continue
		}
		if toval == "+" {
			if len(togval) > 0 {
				allvaluelist = append(allvaluelist, togval)
			}
		} else {
			togval += toval
		}
	}
	//fmt.Println("+++++++++++", allvaluelist)
	if len(allvaluelist) <= 0 {
		allvaluelist = append(allvaluelist, targetval)
	}
	for index := 0; index < len(allvaluelist); index++ {
		val := allvaluelist[index]
		//fmt.Println("========2222222222222==11==", val)
		if val[0:1] == "\"" {
			val = val[1 : len(val)-1]
		} else if val != "true" && val != "false" {
			val = lgetValueFromMap(keyvalMap, val, targetval)
		}
		//fmt.Println("========2222222222222==22==", val)
		realvalstr += val
	}
	return realvalstr
}

func analysisConfig(ffpath string, keyvalMap map[string]string) {
	fmt.Println("=====analysisConfig===222=====", ffpath)
	f, err := os.Open(ffpath)
	if err != nil {
		fmt.Println("Open error: ", err)
		os.Exit(5)
		return
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ReadAll error: ", err)
		os.Exit(5)
		return
	}
	filedata := fmt.Sprintf("%s", buffer)
	filedataEx := []string{}
	var forstr string = ""
	var forstrEx string = ""
	var islinezs bool = false
	var isnumlinezs bool = false

	var issssstr bool = false
	//var preissssstr bool = false
	var issscont bool = false
	var fsssstr string = ""
	for i, val := range filedata {
		toval := fmt.Sprintf("%c", val) + ""
		forstr += toval
		forstrlen := len(forstr)
		if forstrlen <= 1 {
			continue
		}
		laststr := forstr[forstrlen-2 : forstrlen]
		if isnumlinezs {
			if laststr == "*/" {
				isnumlinezs = false
			}
			continue
		} else if laststr == "*/" {
			fmt.Printf("error: unknow */ in %d", i)
			os.Exit(1)
		}
		if islinezs {
			if toval != "\n" {
				continue
			}
			islinezs = false
		}
		if laststr == "//" {
			islinezs = true
			forstrEx = forstrEx[:len(forstrEx)-1]
			continue
		}
		if laststr == "/*" {
			isnumlinezs = true
			forstrEx = forstrEx[:len(forstrEx)-1]
			continue
		}
		/*
		 */
		for {
			issscont = false
			if issssstr {
				if laststr == "\\\"" {
					fsssstr += toval
					break
				}
				if toval == "\n" {
					if forstr[len(forstr)-3:len(forstr)-2] != "\\" {
						fmt.Println("error: string line end not find \\")
						os.Exit(4)
					}
					forstrEx = forstrEx[0 : len(forstrEx)-2]
					issscont = true
					break
				}
				if toval == "\"" {
					issssstr = false
					//preissssstr = true
					break
				}
				fsssstr += toval
			}
			if toval == "\"" {
				issssstr = true
				fsssstr = ""
				break
			}
			break
		}
		if issscont {
			continue
		}
		if issssstr {
			forstrEx += toval
			continue
		}
		// if preissssstr {
		// 	replaceCount++
		// 	repkey := replaceStr + fmt.Sprintf("%d", replaceCount)
		// 	keyvalMap[repkey] = fsssstr
		// 	forstrEx += repkey
		// }
		if toval == "\r" || toval == "\t" || toval == " " {
			continue
		}
		if toval == "\n" || toval == ";" {
			if len(forstrEx) > 0 {
				//fmt.Println("=====1============", forstrEx)
				filedataEx = append(filedataEx, forstrEx)
			}
			forstrEx = ""
		} else {
			//fmt.Println("=====22============", forstrEx)
			forstrEx += toval
		}
	}
	fmt.Printf(">>>>%#v\n", filedataEx)
	paths, _ := filepath.Split(ffpath)
	for index := 0; index < len(filedataEx); index++ {
		tline := filedataEx[index]
		cmd, tarval := lcheckOneline(tline)
		if cmd == "include" {
			tarval = filepath.Clean(tarval)
			if paths != "." && paths != "./" && paths != ".\\" {
				tarval = paths + tarval
			}
			fmt.Println("=====analysisConfig========", paths, tarval)
			analysisConfig(tarval, keyvalMap)
			continue
		}
		keyvalMap[cmd] = tarval
	}
	fmt.Println("-------------------------------------", keyvalMap)
}

// LoadPathConfig outside call load path config
func LoadPathConfig(ffpath string) map[string]string {
	keyvalMap := map[string]string{}
	analysisConfig(ffpath, keyvalMap)
	retmap := map[string]string{}
	for key, val := range keyvalMap {
		tarval2 := lcheckValueReplace(keyvalMap, val)
		retmap[key] = tarval2
		fmt.Printf("%s=%s\n", key, tarval2)
	}
	return retmap
}
