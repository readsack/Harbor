package main

import "fmt"


func main(){
	fmt.Println("Hello World!")
	var arr []int
	arr = append(arr, 1)
	arr = append(arr, 2)
	arr = append(arr, 3)

	fmt.Println(arr[1:2])

}
