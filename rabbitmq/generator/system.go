// Copyright © 2022 Aleksey Barabanov
// Licensed under the Apache License, Version 2.0

package main

import (
	"log"
	"os"
	"strconv"
)

// GetEnv Получение данных из переменных окружения, также передаем дефолт если переменной не найдется
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// CheckErr Проверка ошибки
func CheckErr(err error) { // Check error
	if err != nil {
		log.Fatal(err)
	}
}

// CheckErrLog Проверка ошибки со строковым префиксом (комментом)
func CheckErrLog(err error, comment string) { // Check error
	if err != nil {

		log.Fatalf("%v : %v", comment, err)
	}
}

func convertTextInt(str string) int {
	result, err := strconv.Atoi(str)
	CheckErrLog(err, "error converting")
	return result
}
