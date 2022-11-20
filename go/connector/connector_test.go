package connector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnector(t *testing.T) {
	yc := make(chan YobikouServer)
	tc := make(chan TyuugakuClient)

	go prepareServer(yc)
	time.Sleep(time.Duration(10) * time.Millisecond)
	go prepareClient(tc)

	yobikou := <-yc
	defer yobikou.Close()
	tyuugaku := <-tc
	defer tyuugaku.Close()

	t.Run("ping", testPing(yobikou, tyuugaku))
	t.Run("pong", testPong(yobikou, tyuugaku))
	t.Run("matrix", testMatrix(yobikou, tyuugaku))
}

func prepareServer(ch chan YobikouServer) {
	yobikou, err := NewYobikouServer()
	if err != nil {
		panic(err)
	}

	ch <- yobikou
}

func prepareClient(ch chan TyuugakuClient) {
	tyuugaku, err := NewTyuugakuClient("0.0.0.0")
	if err != nil {
		panic(err)
	}

	ch <- tyuugaku
}

func testPing(yobikou YobikouServer, tyuugaku TyuugakuClient) func(t *testing.T) {
	return func(t *testing.T) {
		go tyuugaku.Send([]byte("ping"))

		buffer := make([]byte, 1024)
		len, err := yobikou.Receive(buffer)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "ping", string(buffer[:len]))
	}
}

func testPong(yobikou YobikouServer, tyuugaku TyuugakuClient) func(t *testing.T) {
	return func(t *testing.T) {
		go yobikou.Send([]byte("pong"))

		buffer := make([]byte, 1024)
		len, err := tyuugaku.Receive(buffer)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "pong", string(buffer[:len]))
	}
}
