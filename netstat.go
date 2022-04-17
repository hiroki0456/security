package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

type NetCat struct {
	buffer string
}

type NetCater interface {
	run()
	send()
	listen()
	handle()
}

var (
	cmd_flag     bool
	execute_flag string
	listen_flag  bool
	port_flag    int
	target_flag  string
	upload_flag  string
)

func init() {
	flag.BoolVar(&cmd_flag, "c", false, "対話型シェルの初期化")
	flag.StringVar(&execute_flag, "e", "", "指定のコマンドの実行")
	flag.BoolVar(&listen_flag, "l", false, "通信待受モード")
	flag.IntVar(&port_flag, "p", 5555, "ポート番号の指定")
	flag.StringVar(&target_flag, "t", "127.0.0.1", "IPアドレスの指定")
}

func execute(cmdstr string) (string, error) {
	cmdstr = strings.TrimSpace(cmdstr)
	splitedcmd := strings.Split(cmdstr, " ")
	if cmdstr == "" {
		return "", errors.New("must input command")
	}

	output, err := exec.Command(splitedcmd[0], splitedcmd[1:]...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (nc *NetCat) run() {
	if listen_flag {
		nc.listen()
	} else {
		nc.send()
	}
}

func (nc *NetCat) send() {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port_flag))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(nc.buffer))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	fmt.Println(string(buf[:n]))
	if err != nil {
		panic(err)
	}
}

func (nc *NetCat) listen() {
	psock, err := net.Listen("tcp", fmt.Sprintf(":%d", port_flag))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := psock.Accept()
		if err != nil {
			panic(err)
		}
		go nc.handle(conn)
	}

}

func (nc *NetCat) handle(conn net.Conn) {
	defer conn.Close()
	if execute_flag != "" {
		output, err := execute(execute_flag)
		if err != nil {
			panic(err)
		}
		conn.Write([]byte(output))
		io.Copy(conn, conn)
	}
}

func main() {
	flag.Parse()
	fmt.Println(listen_flag)
	fmt.Println(execute_flag)
	var buffer string
	if listen_flag {
		buffer = ""
	} else {
		buffer = bufio.NewScanner(os.Stdin).Text()
	}
	nc := NetCat{
		buffer: buffer,
	}
	nc.run()
}
