package main

import (
	"bufio"
	"errors"
	"flag"
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
	flag.BoolVar(&cmd_flag, "c", true, "対話型シェルの初期化")
	flag.StringVar(&execute_flag, "e", "", "指定のコマンドの実行")
	flag.BoolVar(&listen_flag, "l", true, "通信待受モード")
	flag.IntVar(&port_flag, "p", 5555, "ポート番号の指定")
	flag.StringVar(&target_flag, "t", "192.168.1.203", "IPアドレスの指定")
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

}

func main() {
	var buffer string
	if listen_flag {
		buffer = ""
	} else {
		buffer = bufio.NewScanner(os.Stdin).Text()
	}
	nc := NetCat{
		buffer: buffer,
	}
	execute("ls -al")
}
