package main

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	daoTemplatePath = "script/daoTemplate"
	modelFilePath   = "model/%s.go"
	daoFilePath     = "dao/%s.go"

	placeHolder = "XXX"
)

func writeDaoFile(tableName string) {
	// 1.open file
	dao, err := os.Open(daoTemplatePath)
	defer func() {
		err = dao.Close()
		if err != nil {
			log.Fatal("File close Failed.", err)
		}
	}()
	// 2.read code
	data, err := ioutil.ReadAll(dao)
	if err != nil {
		log.Fatal("Write dao file error.", err)
	}
	template := string(data)
	daoCode := strings.ReplaceAll(template, placeHolder, getModelName(tableName))

	// 3.if dao file exists, not create
	targetFile := fmt.Sprintf(daoFilePath, getFileName(tableName))
	_, err = os.Stat(targetFile)
	if !os.IsNotExist(err) {
		log.Fatal("Dao File exists.")
	}
	// 4.create and write code
	daoFile := createFile(targetFile)
	wrote, err := daoFile.WriteString(daoCode)
	if err != nil {
		log.Fatal("Write Dao File error.", zap.Error(err))
	}
	log.Printf("Wrote %d characters Successfully.", wrote)
}

func writeModelFile(tableName string) {
	fileName := getFileName(tableName)
	// if model file exists, not create
	targetFile := fmt.Sprintf(modelFilePath, fileName)
	_, err := os.Stat(targetFile)
	if !os.IsNotExist(err) {
		log.Fatal("Model File exists.")
	}
	// create file and write code
	file := createFile(targetFile)
	defer closeFile(file)

	buffer := bytes.Buffer{}
	buffer.WriteString("package model\n")
	buffer.WriteString(fmt.Sprintf("type %sModel struct {\n", getModelName(tableName)))

	var tail string
	for _, field := range fields {
		tail = fmt.Sprintf("json:\"%s\"", field)
		buffer.WriteString(getFieldName(field))
		buffer.WriteString(" ")
		buffer.WriteString(fieldsTypeMap[field])
		buffer.WriteString(fmt.Sprintf(" `%s`\n", tail))
	}
	buffer.WriteString("}\n")

	byteSlice := buffer.Bytes()
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes into model file.\n", bytesWritten)
}

func getModelName(tableName string) string {
	words := strings.Split(tableName, "_")
	var modelName string
	length := len(words)
	for i := range words {
		//skip prefix 'sp_'
		//if i == 0 {
		//	continue
		//}
		//skip suffix '_tab'
		if length-1 == i {
			break
		}
		modelName += strings.Title(words[i])
	}
	return modelName
}

func getFileName(tableName string) string {
	i := strings.Index(tableName, "_")
	j := strings.LastIndex(tableName, "_")
	if i < j {
		return tableName[i+1 : j]
	}
	if j != -1 {
		return tableName[:j]
	}
	return tableName
}

func closeFile(file *os.File) {
	err = file.Close()
	if err != nil {
		log.Fatal("File close Failed.", err)
	}
}

func createFile(filePath string) *os.File {
	file, err := os.OpenFile(
		filePath,
		os.O_RDWR|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
