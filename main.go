package main

// ############ FILE SERVER TESTING CODE ############

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gofer/internal/file-server"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func sendFile(dir string) error {
	file, err := os.Open(dir)
	// TODO: verbosely log file open failure
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()
	// TODO: verbosely log file stat failure
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	// TODO: verbosely log connection failure
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(fileInfo.Size()))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("written %d bytes over the network\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		sendFile("gofer/assets/rave_test.wav")
	}()

	server := fileserver.FileServer{}
	server.Start()
}

// ############ TEXT SERVER TESTING CODE ############

// func main() {
// 	s := server.NewTextServer(":3000")

// 	go func() {
// 		for msg := range s.GetMsgCh() {
// 			fmt.Printf("received message from connection (%s):%s\n", msg.GetFrom(), string(msg.GetPayload()))
// 		}
// 	}()

// 	log.Fatal(s.Start())
// }
