package connector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testMatrix(yobikou YobikouServer, chugaku ChugakuClient) func(t *testing.T) {
	return func(t *testing.T) {
		table := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		t.Run("ping matrix", testPingMatrix(yobikou, chugaku, table))
		t.Run("pong matrix", testPongMatrix(yobikou, chugaku, table))
	}
}

func testPingMatrix(yobikou YobikouServer, chugaku ChugakuClient, matrix [][]float64) func(t *testing.T) {
	return func(t *testing.T) {
		go chugaku.SendTable(matrix)

		result, err := yobikou.ReceiveTable()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, matrix, result)
	}
}

func testPongMatrix(yobikou YobikouServer, chugaku ChugakuClient, matrix [][]float64) func(t *testing.T) {
	return func(t *testing.T) {
		go yobikou.SendTable(matrix)

		result, err := chugaku.ReceiveTable()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, matrix, result)
	}
}
