package file

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/pkg/errors"
	_log "github.com/wx11055/gone/logger"
	"io"
	"io/ioutil"
	"os"
)

func store(data string, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func load(data string, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

//判断文件或文件夹是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// 创建文件夹
func CreateDirs(path string) error {
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func ReadFileWithIoutil(fullName string, handleLine func(int, string)) error {
	defer func() {
		if r := recover(); r != nil {
			_log.Warn("ReadFile With ioutil Error Warning")
			_log.Warn(r)
		}
	}()
	exist := IsExist(fullName)
	if !exist {
		return errors.New("File not exist!!")
	}
	if body, err := ioutil.ReadFile(fullName); err == nil {
		if err != nil {
			return err
		}
		contents := bytes.Split(body, []byte("\n"))
		for line, content := range contents {
			value := string(bytes.TrimSpace(content))
			handleLine(line, value)
		}
		return nil
	}
	return nil
}
func ReadLineWithOs(fullName string, handleLine func(int, string)) error {
	if fileObj, err := os.Open(fullName); err == nil {
		//if fileObj,err := os.OpenFile(name,os.O_RDONLY,0644); err == nil {
		defer fileObj.Close()
		if body, err := ioutil.ReadAll(fileObj); err == nil {
			contents := bytes.Split(body, []byte("\n"))
			for line, content := range contents {
				value := string(bytes.TrimSpace(content))
				handleLine(line, value)
			}
			return nil
		}
	}
	return nil
}
func ReadLineWithFile(fullName string, handleLine func(int, string)) error {
	if fileObj, err := os.Open(fullName); err == nil {
		defer fileObj.Close()
		//在定义空的byte列表时尽量大一些，否则这种方式读取内容可能造成文件读取不完整
		buf := make([]byte, 1024)
		var body []byte
		for {
			n, err := fileObj.Read(buf[:])
			if err == io.EOF {
				_log.Debug("File read finished!!!")
				break
			}
			if err != nil {
				return errors.New("File read error!!!")
			}
			body = append(body, buf[:n]...)
		}
		contents := bytes.Split(body, []byte("\n"))
		for line, content := range contents {
			value := string(bytes.TrimSpace(content))
			handleLine(line, value)
		}
	}
	return nil
}
func ReadLineWithBufio(fullName string, handleLine func(int, string)) error {
	if fileObj, err := os.Open(fullName); err == nil {
		defer fileObj.Close()
		l := bufio.NewReader(fileObj)
		i := 0
		for {
			i++
			line, isPrefix, err := l.ReadLine()
			if len(line) > 0 && err != nil {
				_log.Errorf("ReadLine returned both data and error: %s", err)
			}
			if isPrefix {
				_log.Errorf("ReadLine returned prefix")
			}
			if err != nil {
				if err != io.EOF {
					_log.Fatalf("Got unknown error: %s", err)
				}
				break
			}
			handleLine(i, string(bytes.TrimSpace(line)))
		}
	}
	return nil
}

func WriteFileWithIoutil(fullName string, body []byte) error {
	if err := ioutil.WriteFile(fullName, body, 0644); err != nil {
		return err
	}
	return nil
}

func WriteFileWithOS(fullName string, body []byte) error {
	/*
		const (
	        O_RDONLY int = syscall.O_RDONLY // 只读打开文件和os.Open()同义
	        O_WRONLY int = syscall.O_WRONLY // 只写打开文件
	        O_RDWR   int = syscall.O_RDWR   // 读写方式打开文件
	        O_APPEND int = syscall.O_APPEND // 当写的时候使用追加模式到文件末尾
	        O_CREATE int = syscall.O_CREAT  // 如果文件不存在，此案创建
	        O_EXCL   int = syscall.O_EXCL   // 和O_CREATE一起使用, 只有当文件不存在时才创建
	        O_SYNC   int = syscall.O_SYNC   // 以同步I/O方式打开文件，直接写入硬盘.
	        O_TRUNC  int = syscall.O_TRUNC  // 如果可以的话，当打开文件时先清空文件
			)
	*/
	fileObj, err := os.OpenFile(fullName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		_log.Error("Failed to open the file", err.Error())
		return err
	}
	defer fileObj.Close()
	// 或者使用io 写操作
	//	if  _,err := io.WriteString(fileObj,string(body));err != nil {
	if _, err := fileObj.Write(body); err != nil {
		return err
	}
	return nil
}

func WriteFileBufio(fullName string, flag int, body []byte) error {
	//flag = os.O_RDWR|os.O_CREATE|os.O_APPEND
	if fileObj, err := os.OpenFile(fullName, flag, 0644); err == nil {
		defer fileObj.Close()
		writeObj := bufio.NewWriter(fileObj)
		//使用Write方法,需要使用Writer对象的Flush方法将buffer中的数据刷到磁盘
		if _, err := writeObj.Write(body); err != nil {
			return err
		}
		if err := writeObj.Flush(); err != nil {
			return err
		}
	}
	return nil
}
