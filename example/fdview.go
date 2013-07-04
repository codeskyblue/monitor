package main

import (
	"flag"
	"fmt"
	"github.com/shxsun/monitor"
	"strings"
    "os"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s <program path>\n", os.Args[0])
		return
	}

	monitor.Refresh()
	monitor.GoRefresh()
	fmt.Println("NCPU: ", monitor.Ncpu())
	fmt.Println("CPU: ", monitor.Cpu())
	pids, _ := monitor.Pids()
	fmt.Println("Len proc:", len(pids))

	for _, pid := range pids {
		pi, err := monitor.Pid(pid)
		if err != nil {
			continue
		}

		if strings.Contains(pi.Exe, flag.Arg(0)) {
		fmt.Println("----------------------------")
			fmt.Println("Pid: ", pid)
			fmt.Println("Exe: ", pi.Exe)
            fmt.Println("Fd size: ", len(pi.Fd))
		}
	}
}
