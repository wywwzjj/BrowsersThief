package utils

import (
	"io"
	"log"
	"os"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CopyFile(source, dest string) bool {
	if source == "" || dest == "" {
		log.Println("source or dest is null")
		return false
	}

	sourceOpen, err := os.Open(source)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer sourceOpen.Close()

	destOpen, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer destOpen.Close()

	_, err = io.Copy(destOpen, sourceOpen)
	if err != nil {
		log.Println(err.Error())
		return false
	} else {
		return true
	}
}

func PrintSeparator(funcName string) {
	equalOp := strings.Repeat("=", 20)
	log.Println(equalOp + funcName + equalOp)
}
