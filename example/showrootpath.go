package main

import (
	"fmt"
	"github.com/shxsun/monitor"
)

func main() {
	fmt.Println("start")
	monitor.Refresh()
	monitor.GoRefresh()
	fmt.Println("CPU: ", monitor.Ncpu())
	pids, _ := monitor.Pids()
	fmt.Println("Len proc:", len(pids))
	if len(pids) < 100 {
		return
	}
	for _, pid := range pids[40:100] {
		pi, err := monitor.Pid(pid)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(pid, "CWD:", pi.Exe, len(pi.Fd))
	}
}
