/*
Monitor /proc

Get Cpu, Mem, and process's Cpu, Mem.
*/

package monitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const BASE_DIR = "/proc"

var ProcPidList = make(map[int]ProcInfo, 1000) // TODO  pid->procinfo
var ProcExeList = make(map[string][]int, 1000) // TODO  exe -> pid
func init() {                                  // TODO
}

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

// /proc/[pid]/stat
// https://www.kernel.org/doc/man-pages/online/pages/man5/proc.5.html
func ProcPidStat(pid int) ([]byte, error) {
	filename := `/proc/` + strconv.FormatInt(int64(pid), 10) + `/stat`
	return ioutil.ReadFile(filename)
}

type Stat struct {
	Utime, Stime uint64
}

type ProcInfo struct {
	Pid  int
	Exe  string
	Root string
	Cwd  string
	Fd   []string
	Stat ProcStat
}

// Update ProcInfo
func (pi *ProcInfo) Update() (err error) {
	basedir := filepath.Join(BASE_DIR, strconv.Itoa(pi.Pid))
	pi.Exe, _ = os.Readlink(filepath.Join(basedir, "exe"))
	pi.Cwd, _ = os.Readlink(filepath.Join(basedir, "cwd"))
	pi.Root, _ = os.Readlink(filepath.Join(basedir, "root"))

	//pi.Stat.Update()
	// read fd
	fs, err := ioutil.ReadDir(filepath.Join(basedir, "fd"))
	if err != nil {
		return
	}
	pi.Fd = make([]string, 0, 100)
	for _, f := range fs {
		p, err := os.Readlink(filepath.Join(basedir, "fd", f.Name()))
		if err != nil {
			continue
		}
		pi.Fd = append(pi.Fd, p)
	}
	return
}

type ProcStat struct {
	Pid                          int
	State                        string
	Ppid, Pgrp                   int
	Session                      int
	Utime, Stime, Cutime, Cstime uint64
}

// Read from /proc/[num]/stat
func (ps *ProcStat) Update() (err error) {
	content, err := ioutil.ReadFile(filepath.Join(BASE_DIR, strconv.Itoa(ps.Pid), "stat"))
	if err != nil {
		return
	}
	var ig string // ignore variable

	fmt.Sscan(string(content),
		&ig, &ig, &ps.State, &ps.Ppid, &ps.Pgrp, &ps.Session,
		&ig, &ig, &ig, &ig, &ig, &ig, &ig,
		&ps.Utime, &ps.Stime, &ps.Cutime, &ps.Cstime)
	//	fmt.Println("PID:", state, ppid, utime, stime, cutime, cstime)
	return
}
