package connector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnector(t *testing.T) {
	yc := make(chan YobikouServer)
	tc := make(chan ChugakuClient)

	go prepareServer(yc)
	time.Sleep(time.Duration(10) * time.Millisecond)
	go prepareClient(tc)

	yobikou := <-yc
	defer yobikou.Close()
	chugaku := <-tc
	defer chugaku.Close()

	t.Run("ping", testPing(yobikou, chugaku))
	t.Run("pong", testPong(yobikou, chugaku))
	t.Run("matrix", testMatrix(yobikou, chugaku))
}

func prepareServer(ch chan YobikouServer) {
	yobikou, err := NewYobikouServer()
	if err != nil {
		panic(err)
	}

	ch <- yobikou
}

func prepareClient(ch chan ChugakuClient) {
	chugaku, err := NewChugakuClient("0.0.0.0")
	if err != nil {
		panic(err)
	}

	ch <- chugaku
}

func testPing(yobikou YobikouServer, chugaku ChugakuClient) func(t *testing.T) {
	return func(t *testing.T) {
		go chugaku.Send([]byte("ping"))

		buffer := make([]byte, 1024)
		len, err := yobikou.Receive(buffer)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "ping", string(buffer[:len]))
	}
}

func testPong(yobikou YobikouServer, chugaku ChugakuClient) func(t *testing.T) {
	return func(t *testing.T) {
		go yobikou.Send([]byte("pong"))

		buffer := make([]byte, 1024)
		len, err := chugaku.Receive(buffer)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "pong", string(buffer[:len]))
	}
}
