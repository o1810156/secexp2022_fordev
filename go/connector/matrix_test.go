package connector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testMatrix(yobikou YobikouServer, tyuugaku TyuugakuClient) func(t *testing.T) {
	return func(t *testing.T) {
		table := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		t.Run("ping matrix", testPingMatrix(yobikou, tyuugaku, table))
		t.Run("pong matrix", testPongMatrix(yobikou, tyuugaku, table))
	}
}

func testPingMatrix(yobikou YobikouServer, tyuugaku TyuugakuClient, matrix [][]float64) func(t *testing.T) {
	return func(t *testing.T) {
		go tyuugaku.SendTable(matrix)

		result, err := yobikou.ReceiveTable()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, matrix, result)
	}
}

func testPongMatrix(yobikou YobikouServer, tyuugaku TyuugakuClient, matrix [][]float64) func(t *testing.T) {
	return func(t *testing.T) {
		go yobikou.SendTable(matrix)

		result, err := tyuugaku.ReceiveTable()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, matrix, result)
	}
}
