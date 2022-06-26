package main

import (
	"math/rand"
	"reflect"
)

func createRandomValueArray[T Number](n int64) []T {
	arr := []T{}

	v := reflect.ValueOf(new(T))

	for i := 0; int64(i) < n; i++ {
		switch v.Type().Elem().Kind() {
		case reflect.Float64:
			arr = append(arr, T(rand.Float64()))
		case reflect.Int64:
			arr = append(arr, T(rand.Int63()))
		}
	}
	return arr
}

func createRandomMatrix[T Number](n int64) [][]T {
	matrix := [][]T{}
	for i := 0; int64(i) < n; i++ {
		matrix = append(matrix, createRandomValueArray[T](n))
	}
	return matrix
}
