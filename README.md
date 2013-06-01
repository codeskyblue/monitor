Monitor of Linux Process
=========================
2013-5-28 sun shengxiang

Used in Linux system.

Parse /proc info and get Cpu, Mem, Hostname, Pids (and each state)

## Download

    go get github.com/shxsun/monitor

## Example
    
    package main

    import "fmt"
    import "time"
    import "github.com/shxsun/monitor"

    func main(){
        // monitor.Interval = time.Second * 2   // default 1s, (refresh gap)
        for {
            fmt.Printf("Cpu usage: %.1f%%\n", monitor.Cpu() * 100)
            time.Sleep(5e8) // 0.5s
        }
    }
