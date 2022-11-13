package main

import (
	"fmt"
	connector "secexp2022/connector"
)

func main() {
	tyuugaku, err := connector.NewTyuugakuClient("0.0.0.0")
	if err != nil {
		panic(err)
	}
	defer tyuugaku.Close()

	tyuugaku.Send([]byte("ping"))

	buffer := make([]byte, 1024)
	len, err := tyuugaku.Receive(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buffer[:len]))

	table := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	err = tyuugaku.SendTable(table)
	if err != nil {
		panic(err)
	}

	result, err := tyuugaku.ReceiveTable()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
