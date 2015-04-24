package send

import (
	"chinchilla/mssg"
	"chinchilla/schedule"
	"chinchilla/types"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Thread to send responses back to hosts
func Client(RespQueue chan mssg.WorkResp, jobs *types.MapJ) {
	for {
		resp := <-RespQueue
		json_resp, _ := json.Marshal(resp)
		jobs.L.Lock()
		w := jobs.M[resp.WId].W
		// allow cross domain AJAX requests
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		_, err := w.Write(json_resp)
		close(jobs.M[resp.WId].Sem)
		jobs.L.Unlock()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		}

	}
}

func Scheduler(w http.ResponseWriter, ReqQueue chan mssg.WorkReq, typ int, arg1 string, id uint32, jobs *types.MapJ) {
	jobs.L.Lock()
	jobs.M[id] = types.Job{W: w, Sem: make(chan struct{})}
	jobs.L.Unlock()
	ReqQueue <- mssg.WorkReq{Type: uint8(typ), Arg1: arg1, WId: id}

}

func ReScheduler(r mssg.WorkReq, ReqQueue chan mssg.WorkReq) {
	ReqQueue <- r
}

func Node(ReqQueue chan mssg.WorkReq, workers *types.MapQ, algo string) {

	for {
		var node uint32
		req := <-ReqQueue
		if len(workers.M) == 0 {
			fmt.Println("No workers you dangus")
			os.Exit(1)
		}
		req.STime = time.Now()

		if algo == "rr" {
			node = schedule.RoundRobin(workers, req.Type)
		} else if algo == "sq" {
			node = schedule.ShortestQ(workers, req.Type)
		} else {
			node = schedule.ShortestQ(workers, req.Type)
		}

		workers.L.Lock()
		tmp := workers.M[node]
		tmp.Reqs = append(workers.M[node].Reqs, req)
		tmp.QLen = workers.M[node].QLen + 1
		workers.M[node] = tmp
		err := workers.M[node].Enc.Encode(req)

		workers.L.Unlock()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		}
	}
}
