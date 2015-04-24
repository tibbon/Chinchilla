package main

import (
	// "fmt"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {

	running := StartUp("rr", 5)
	running = AddWorker(running)
	running = RemoveWorker(running, 3)
	running = AddWorker(running)
	// KillServer(running)
	// fmt.Println("Killing them alll")
	// for i := 0; i < len(running); i++ {
	// 	running[i].Process.Kill()
	// }

}

func StartUp(algo string, numWorker int) []*exec.Cmd {
	running := make([]*exec.Cmd, 1)
	running[0] = exec.Command("./scheduler", "8080", "8081", algo) // start scheduler
	ip := GetIp()
	for i := 1; i < numWorker+1; i++ {
		running = append(running, exec.Command("./worker", ip+":8081", strconv.Itoa(i)))
	}
	running[0].Start()
	time.Sleep(100 * time.Millisecond)
	for i := 1; i < len(running); i++ {
		running[i].Start()
	}

	return running
}

func AddWorker(running []*exec.Cmd) []*exec.Cmd {
	fmt.Println("adding worker")
	numWorkers := len(running)
	ip := GetIp()

	for i := 1; i < numWorkers; i++ {
		if running[i] == nil {
			running[i] = exec.Command("./worker", ip+":8081", strconv.Itoa(i))
			running[i].Start()
			return running
		}
	}
	running = append(running, exec.Command("./worker", ip+":8081", strconv.Itoa(numWorkers)))
	running[numWorkers].Start()
	return running

}

func GetIp() string {
	ips, _ := net.InterfaceAddrs()
	ip := ips[1].String()
	if loc := strings.Index(ip, "/"); loc != -1 {
		return ip[:loc]
	}
	return ip
}

func RemoveWorker(running []*exec.Cmd, i int) []*exec.Cmd {
	fmt.Println("Removing worker")
	running[i].Process.Kill()
	running[i] = nil
	return running
}

func KillServer(running []*exec.Cmd) {
	for i := 0; i < len(running); i++ {
		running[i].Process.Kill()
	}
}
