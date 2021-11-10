package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	err      error
	fileInfo os.FileInfo

	tableName     string
	fields        []string
	fieldsTypeMap = make(map[string]string)
)

const (
	sqlScriptPath = "script/tab.sql" // need create firstly, contains one table
	entityPath    = "dao/entity/%s.go"
)

// auto generate code by creating table SQL
func main() {
	fileInfo, err = os.Stat(sqlScriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
		}
	}
	log.Println("File does exist. File information:")
	// read sql, get table name, fields, fields Type
	sqlFile, err := os.Open(sqlScriptPath)
	defer func() {
		err = sqlFile.Close()
		if err != nil {
			log.Fatal("File close Failed.", err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(sqlFile)
	if err != nil {
		log.Fatal(err)
	}
	sql := string(data)
	i := strings.Index(sql, "(")
	tableName = getTableName(sql[:i])
	j := strings.LastIndex(sql, ")")
	scanFields(sql[i+1 : j])
	writeEntityFile(tableName)
	writeModelFile(tableName)
	writeDaoFile(tableName)
}

func writeEntityFile(tableName string) {
	// if file exists, not create
	targetFile := fmt.Sprintf(entityPath, tableName)
	fileInfo, err = os.Stat(targetFile)
	if !os.IsNotExist(err) {
		log.Fatalf("%s Entity File exists.", fileInfo.Name())
	}
	// create file and write code
	file := createFile(targetFile)
	defer closeFile(file)
	entityName := getEntityName(tableName)
	fmt.Print(entityName)

	buffer := bytes.Buffer{}
	buffer.WriteString("package entity\n")
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", entityName))
	var tail string
	for _, field := range fields {
		if field == "id" {
			tail = fmt.Sprintf("gorm:\"column:%s;primary_key;AUTO_INCREMENT\"", field)
		} else {
			tail = fmt.Sprintf("gorm:\"column:%s\"", field)
		}
		buffer.WriteString(getFieldName(field))
		buffer.WriteString(" ")
		buffer.WriteString(fieldsTypeMap[field])
		buffer.WriteString(fmt.Sprintf(" `%s`\n", tail))
	}
	buffer.WriteString("}\n")

	buffer.WriteString(fmt.Sprintf("func (*%s) TableName() string {\n", entityName))
	buffer.WriteString(fmt.Sprintf("return %q\n}", tableName))
	byteSlice := buffer.Bytes()
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AUTO Generate SUCCESS ! Wrote %d bytes.\n", bytesWritten)
}

func getTableName(create string) string {
	statement := strings.Fields(create)
	// ignore word case
	if strings.ToUpper(statement[0]) != "CREATE" {
		log.Fatal("create statement error")
	}
	if strings.ToUpper(statement[1]) != "TABLE" {
		log.Fatal("create statement error")
	}
	return strings.Trim(statement[2], "`")
}

func getEntityName(tableName string) string {
	words := strings.Split(tableName, "_")
	var entityName string
	for i := range words {
		//skip prefix 'sp_' //todo
		//if i == 0 {
		//	continue
		//}
		entityName += strings.Title(words[i])
	}
	return entityName
}

func getFieldName(fieldName string) string {
	words := strings.Split(fieldName, "_")
	var field string
	for i := range words {
		field += strings.Title(words[i])
	}
	return field
}

func scanFields(sql string) {
	sql = strings.Trim(sql, "\n")
	lines := strings.Split(sql, ",\n")
	for _, line := range lines {
		/*
			if !strings.Contains(line, "COMMENT") {
				continue
			}
		*/
		// depend on SQL statement style
		fieldDefine := strings.Fields(line)
		if len(fieldDefine) > 1 && strings.HasPrefix(fieldDefine[0], "`") {
			field := strings.Trim(fieldDefine[0], "`")
			fields = append(fields, field)
			fieldsTypeMap[field] = getFieldType(fieldDefine[1])
		}
	}
}

// must be complete todo
func getFieldType(define string) string {
	i := strings.Index(define, "(")
	if i == -1 {
		i = len(define)
	}
	switch define[:i] {
	case "tinyint":
		return "int8"
	case "int":
		return "int32"
	case "bigint":
		return "int64"
	case "varchar":
		return "string"
	case "text":
		return "string"
	case "json":
		return "[]byte"
	}
	return ""
}
