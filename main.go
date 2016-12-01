package main

import (
	"fmt"
	"goini"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	//PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			//files = append(files, dirPth+PthSep+fi.Name())
			files = append(files, fi.Name())
		}
	}
	return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func main() {

	os.Remove("flatc_gen.cmd")
	flog, err := os.OpenFile("flatc_gen.cmd", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer flog.Close()

	config := goini.SetConfig("conf.ini")

	srcPath := config.GetValue("path", "src")
	netPath := config.GetValue("path", "net")
	goPath := config.GetValue("path", "go")

	files, err := ListDir(srcPath, ".idl")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 日志
	l := log.New(flog, "", os.O_APPEND)
	for _, file := range files {
		arg := "flatc -n -o " + netPath + " " + srcPath + file
		cmd := exec.Command("flatc", "-n", "-o", netPath, srcPath+file)
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(arg)
			l.Println(arg)
		}
		arg2 := "flatc -g -o " + goPath + " " + srcPath + file
		cmd2 := exec.Command("flatc", "-g", "-o", goPath, srcPath+file)
		err2 := cmd2.Start()
		if err2 != nil {
			fmt.Println(err2)
		} else {
			fmt.Println(arg2)
			l.Println(arg2)
		}
	}
	fmt.Println("done!")

	/*
		files, err := ListDir("D:\\Go", ".txt")
		fmt.Println(files, err)
		files, err = WalkDir("E:\\Study", ".pdf")
		fmt.Println(files, err)
	*/
}
