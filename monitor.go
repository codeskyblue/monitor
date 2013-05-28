package monitor

import (
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
)

func Pids() ([]int, error) {
    f, err := os.Open(`/proc`)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    names, err := f.Readdirnames(-1)
    if err != nil {
        return nil, err
    }
    pids := make([]int, 0, len(names))
    for _, name := range names {
        if pid, err := strconv.ParseInt(name, 10, 0); err == nil {
            pids = append(pids, int(pid))
        }
    }
    return pids, nil
}

func ProcPidStat(pid int) ([]byte, error) {
    // /proc/[pid]/stat
    // https://www.kernel.org/doc/man-pages/online/pages/man5/proc.5.html
    filename := `/proc/` + strconv.FormatInt(int64(pid), 10) + `/stat`
    return ioutil.ReadFile(filename)
}

/*
func main() {
    pids, err := Pids()
    if err != nil {
        fmt.Println("pids:", err)
        return
    }
    if len(pids) > 0 {
        pid := pids[0]
        stat, err := ProcPidStat(pid)
        if err != nil {
            fmt.Println("pid:", pid, err)
            return
        }
        fmt.Println(`/proc/[pid]/stat:`, string(stat))
    }
}
*/
