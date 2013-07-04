package main

import (
	"fmt"
	"github.com/shxsun/monitor"
	"time"
)

func main() {
	monitor.GoRefresh()
	for {
		time.Sleep(1e9)
		fmt.Printf("Cpu: %.2f%%\n", monitor.Cpu()*100)
	}
}
