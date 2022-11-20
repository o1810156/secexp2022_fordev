package connector

import (
	"fmt"
	"net"
	"time"
)

const PORT = ":10000"
const PROTOCOL = "tcp"
const DEADLINE = time.Duration(10) * time.Second

type Communicator interface {
	GetConn() net.Conn
}

func rawSend(com Communicator, data []byte) error {
	conn := com.GetConn()
	conn.Write(data)

	return nil
}

func send(com Communicator, data []byte) error {
	err := rawSend(com, data)
	if err != nil {
		return err
	}

	_, err = rawReceive(com, []byte{})

	return err
}

func rawReceive(com Communicator, buffer []byte) (int, error) {
	conn := com.GetConn()
	len, err := conn.Read(buffer)

	return len, err
}

func receive(com Communicator, buffer []byte) (int, error) {
	len, err := rawReceive(com, buffer)
	if err != nil {
		return 0, err
	}

	err = rawSend(com, []byte{})
	if err != nil {
		return 0, err
	}

	return len, nil
}

// 予備校に該当するよい英単語がないためローマ字表記
type YobikouServer struct {
	conn net.Conn
}

func NewYobikouServer() (YobikouServer, error) {
	tcpAddr, err := net.ResolveTCPAddr(PROTOCOL, PORT)
	if err != nil {
		return YobikouServer{}, err
	}

	listener, err := net.ListenTCP(PROTOCOL, tcpAddr)
	if err != nil {
		return YobikouServer{}, err
	}

	fmt.Println("中学側に接続します……")

	conn, err := listener.Accept()
	if err != nil {
		return YobikouServer{}, err
	}

	return YobikouServer{conn}, nil
}

func (yobikou YobikouServer) GetConn() net.Conn {
	return yobikou.conn
}

func (yobikou YobikouServer) Close() error {
	return yobikou.conn.Close()
}

func (yobikou YobikouServer) RawSend(data []byte) error {
	return rawSend(yobikou, data)
}

func (yobikou YobikouServer) Send(data []byte) error {
	return send(yobikou, data)
}

func (yobikou YobikouServer) RawReceive(buffer []byte) (int, error) {
	return rawReceive(yobikou, buffer)
}

func (yobikou YobikouServer) Receive(buffer []byte) (int, error) {
	return receive(yobikou, buffer)
}

// 予備校に合わせ中学校もローマ字で
type TyuugakuClient struct {
	conn net.Conn
}

func NewTyuugakuClient(ServerAddr string) (TyuugakuClient, error) {
	fmt.Println("予備校に接続します……")

	conn, err := net.DialTimeout(PROTOCOL, ServerAddr+PORT, DEADLINE)
	if err != nil {
		return TyuugakuClient{}, err
	}

	return TyuugakuClient{conn}, nil
}

func (tyuugaku TyuugakuClient) GetConn() net.Conn {
	return tyuugaku.conn
}

func (tyuugaku TyuugakuClient) Close() error {
	return tyuugaku.conn.Close()
}

func (tyuugaku TyuugakuClient) RawSend(data []byte) error {
	return rawSend(tyuugaku, data)
}

func (tyuugaku TyuugakuClient) Send(data []byte) error {
	return send(tyuugaku, data)
}

func (tyuugaku TyuugakuClient) RawReceive(buffer []byte) (int, error) {
	return rawReceive(tyuugaku, buffer)
}

func (tyuugaku TyuugakuClient) Receive(buffer []byte) (int, error) {
	return receive(tyuugaku, buffer)
}