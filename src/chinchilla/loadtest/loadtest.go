package loadtest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type LoadTester struct {
	ApiPort string
	Socket  websocket.Conn
	Stop    bool
	Running []*exec.Cmd
}

type TestParams struct {
	TypeOne     int    `json:"type_1"`
	TypeTwo     int    `json:"type_2"`
	TypeThree   int    `json:"type_3"`
	AlgType     string `json:"alg_type"`
	WorkerCount int    `json:"worker_count"`
}

type workResp struct {
	Type  uint8        `json:"type"`
	Id    uint32       `json:"id"`
	Data  workRespData `json:"data"` // Can be json kind of thing or string
	WId   uint32       `json:"work_id"`
	PTime float64      `json:"return_time"`
	QVal  float64      `json:"q_val"`
}
type workRespData struct {
	Desc string `json:"description"`
}

type response struct {
	response *http.Response
}

func NewLoadTester(apiPort string, c websocket.Conn) *LoadTester {
	return &LoadTester{apiPort, c, false, nil}
}

func (t *LoadTester) LoadTest(w http.ResponseWriter, r *http.Request, p *TestParams) {

	numRequests := 4
	ch := make(chan *response, numRequests*1000)
	url := "http://localhost:" + t.ApiPort + "/api/"
	reqPerType := [3]int{p.TypeOne, p.TypeTwo, p.TypeThree}

	t.Running = startUp(p.AlgType, p.WorkerCount, t.ApiPort)

	time.Sleep(1000 * time.Millisecond)

	fmt.Println(p.AlgType)
	fmt.Println(p.WorkerCount)

	go func() {
		for !t.Stop {
			time.Sleep(1000 * time.Millisecond)
			for j := 0; j < 3; j++ {
				for i := 0; i < reqPerType[j]; i++ {
					go func(v int) {
						fmt.Println("sending")
						resp, err := http.Get(url + strconv.Itoa(v+1) + "/test")
						if err == nil {
							ch <- &response{resp}
						}
					}(j)
				}
			}
		}
		fmt.Print("halting sender")
	}()

	go func() {
		for !t.Stop {
			fmt.Println("trying to read")
			select {
			case r := <-ch:
				temp_struct := new(workResp)
				fmt.Println(r)
				body, err := ioutil.ReadAll(r.response.Body)
				if err != nil {
					fmt.Println("error")
				}
				err = json.Unmarshal(body, &temp_struct)
				if err != nil {
					fmt.Println("error")
				}
				fmt.Println("responding: ", temp_struct)
				t.Socket.WriteJSON(temp_struct)
			}
		}
	}()

}

func (t *LoadTester) StopTest() {
	killServer(t.Running)
}

func (t *LoadTester) KillWorker(w http.ResponseWriter, r *http.Request, wid int) {

	t.Running = removeWorker(t.Running, wid)

}

func (t *LoadTester) AddWorker() {
	t.Running = addWorker(t.Running)
}

func startUp(algo string, numWorker int, apiPort string) []*exec.Cmd {
	running := make([]*exec.Cmd, 1)
	running[0] = exec.Command("./scheduler", apiPort, "9020", algo) // start scheduler
	ip := getIp()
	for i := 1; i < numWorker+1; i++ {
		running = append(running, exec.Command("./worker", ip+":9020", strconv.Itoa(i)))
	}
	running[0].Start()
	time.Sleep(1000 * time.Millisecond)
	for i := 1; i < len(running); i++ {
		fmt.Println("running here")
		running[i].Start()
	}

	return running
}

func addWorker(running []*exec.Cmd) []*exec.Cmd {
	fmt.Println("adding worker")
	numWorkers := len(running)
	ip := getIp()

	for i := 1; i < numWorkers; i++ {
		if running[i] == nil {
			running[i] = exec.Command("./worker", ip+":9020", strconv.Itoa(i))
			running[i].Start()
			return running
		}
	}
	running = append(running, exec.Command("./worker", ip+":9020", strconv.Itoa(numWorkers)))
	running[numWorkers].Start()
	return running

}

func getIp() string {
	ips, _ := net.InterfaceAddrs()
	ip := ips[1].String()
	if loc := strings.Index(ip, "/"); loc != -1 {
		return ip[:loc]
	}
	return ip
}

func removeWorker(running []*exec.Cmd, i int) []*exec.Cmd {
	fmt.Println("Removing worker")
	running[i].Process.Kill()
	running[i] = nil
	return running
}

func killServer(running []*exec.Cmd) {
	for i := 0; i < len(running); i++ {
		if running[i] != nil {
			running[i].Process.Kill()
		}
	}
}
