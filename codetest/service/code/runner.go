package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	//创建一个命令
	cmd := exec.Command("go", "run", "runnerCode/runCode.go")
	//配置相关数据
	var stderr, out bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinpip, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(stdinpip, "20 10\n")
	//运行
	if err := cmd.Run(); err != nil {
		fmt.Println(err, stderr.String())
	}
	fmt.Println(out.String())

}
