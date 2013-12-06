Monitor of Linux Process
=========================
[![Build Status](https://drone.io/github.com/shxsun/monitor/status.png)](https://drone.io/github.com/shxsun/monitor/latest)

Used in Linux system.

Parse /proc info and get Cpu, Mem, Hostname, Pids (and each state)

Not maintained now. Please refer to the new project <https://github.com/cloudfoundry/gosigar>
## Download

    go get github.com/shxsun/monitor

## Example (see more in example dir)
    
    package main

    import "fmt"
    import "time"
    import "github.com/shxsun/monitor"

    func main(){
        // monitor.Interval = time.Second * 2   // default 1s, (refresh gap)
        monitor.GoRefresh()
        for {
            fmt.Printf("Cpu usage: %.1f%%\n", monitor.Cpu() * 100)
            time.Sleep(5e8) // 0.5s
        }
    }
