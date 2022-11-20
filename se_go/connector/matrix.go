package connector

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type matrixSize struct {
	Rows int
	Cols int
}

type matrix struct {
	Data [][]string
}

func newMatrix(table [][]float64) (matrixSize, matrix, error) {
	rows := len(table)
	if rows < 1 {
		return matrixSize{}, matrix{}, fmt.Errorf("table must have at least one row")
	}
	cols := len(table[0])

	str_table := make([][]string, rows)
	for i := 0; i < rows; i++ {
		str_table[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			if len(table[i]) != cols {
				return matrixSize{}, matrix{}, fmt.Errorf("table must be a rectangle")
			}

			str_table[i][j] = strconv.FormatFloat(table[i][j], 'E', -1, 64)
		}
	}

	return matrixSize{rows, cols}, matrix{str_table}, nil
}

func strMatrix2FloatMatrix(table [][]string) ([][]float64, error) {
	rows := len(table)
	if rows < 1 {
		return nil, fmt.Errorf("table must have at least one row")
	}
	cols := len(table[0])

	float_table := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		float_table[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			if len(table[i]) != cols {
				return nil, fmt.Errorf("table must be a rectangle")
			}

			n, err := strconv.ParseFloat(table[i][j], 64)
			if err != nil {
				return nil, err
			}
			float_table[i][j] = n
		}
	}

	return float_table, nil
}

func sendTable(communicator Communicator, table [][]float64) error {
	size, matrix, err := newMatrix(table)
	if err != nil {
		return err
	}

	err = sendSize(communicator, size)
	if err != nil {
		return err
	}

	return sendMatrix(communicator, matrix)
}

func sendSize(communicator Communicator, size matrixSize) error {
	data, err := json.Marshal(size)
	if err != nil {
		return err
	}

	err = send(communicator, data)
	if err != nil {
		return err
	}

	// wait reply size
	buffer := make([]byte, 1024)
	_, err = receive(communicator, buffer)

	return err
}

func sendMatrix(communicator Communicator, mtrx matrix) error {
	data, err := json.Marshal(mtrx)
	if err != nil {
		return err
	}

	err = send(communicator, data)
	if err != nil {
		return err
	}

	// wait reply matrix
	buffer := make([]byte, 1024)
	_, err = receive(communicator, buffer)

	return err
}

const DOUBLE_BUFFER = 20

func receiveTable(communicator Communicator) ([][]float64, error) {
	size, err := receiveSize(communicator)
	if err != nil {
		return nil, err
	}

	return receiveMatrix(communicator, size)
}

func receiveSize(communicator Communicator) (matrixSize, error) {
	buffer := make([]byte, 1024)
	len, err := receive(communicator, buffer)
	if err != nil {
		return matrixSize{}, err
	}

	// reply size
	err = send(communicator, []byte("size"))
	if err != nil {
		return matrixSize{}, err
	}

	var size matrixSize
	err = json.Unmarshal(buffer[:len], &size)
	if err != nil {
		return matrixSize{}, err
	}

	return size, nil
}

func receiveMatrix(communicator Communicator, size matrixSize) ([][]float64, error) {
	buffer := make([]byte, size.Rows*size.Cols*DOUBLE_BUFFER+DOUBLE_BUFFER)
	len, err := receive(communicator, buffer)
	if err != nil {
		return nil, err
	}

	// reply matrix
	err = send(communicator, []byte("matrix"))
	if err != nil {
		return nil, err
	}

	var matrix matrix
	err = json.Unmarshal(buffer[:len], &matrix)
	if err != nil {
		return nil, err
	}

	return strMatrix2FloatMatrix(matrix.Data)
}

func (chugaku ChugakuClient) SendTable(table [][]float64) error {
	return sendTable(chugaku, table)
}

func (chugaku ChugakuClient) ReceiveTable() ([][]float64, error) {
	return receiveTable(chugaku)
}

func (yobikou YobikouServer) SendTable(table [][]float64) error {
	return sendTable(yobikou, table)
}

func (yobikou YobikouServer) ReceiveTable() ([][]float64, error) {
	return receiveTable(yobikou)
}
