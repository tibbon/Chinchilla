package schedule

import (
	"chinchilla/types"
	"fmt"
)

func RoundRobin(workers *types.MapQ) uint32 {
	for k, v := range workers.M {
		workers.L.Lock()
		if !workers.M[k].Sent {
			v.Sent = true
			workers.M[k] = v
			workers.L.Unlock()
			return k
		}
		workers.L.Unlock()
	}
	for k, v := range workers.M {
		workers.L.Lock()
		v.Sent = false
		workers.M[k] = v
		workers.L.Unlock()
	}
	for k, v := range workers.M {
		v.Sent = true
		workers.M[k] = v
		return k
	}
	return 0
}

func ShortestQ(workers *types.MapQ, typ uint8) uint32 {
	fmt.Println("in shortest Q")
	first := true
	var node uint32
	var min float64
	for k, v := range workers.M {
		if first {
			node = k
			min = v.QVal
			first = false
		} else {
			if v.QVal <= min {
				min = v.QVal
				node = k
			}
		}
		fmt.Printf("node %d has a QVal of %f\n", k, v.QVal)
	}
	for k, v := range workers.M[node].AvgTimes {
		fmt.Printf("type %d has qval of %f\n", k, v)
	}
	workers.L.Lock()
	tmp := workers.M[node]
	tmp.QVal += workers.M[node].AvgTimes[typ] + .1
	workers.M[node] = tmp
	workers.L.Unlock()
	fmt.Printf("chose node %d\n", node)
	return node
}
