package main

import (
	"fmt"
	connector "secexp2022/connector"
)

func main() {
	yobikou, err := connector.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer yobikou.Close()

	buffer := make([]byte, 1024)
	len, err := yobikou.Receive(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buffer[:len]))

	yobikou.Send([]byte("pong"))

	table, err := yobikou.ReceiveTable()
	if err != nil {
		panic(err)
	}

	fmt.Println(table)

	// ここでtableを使って計算する

	// ここで計算結果をyobikouに送る
	err = yobikou.SendTable(table)
	if err != nil {
		panic(err)
	}
}
