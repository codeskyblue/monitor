/*
Monitor /proc

Get Cpu, Mem, and process's Cpu, Mem.

Variable:Proc contains all information about system current status.

Proc struct {
	Cpu	 float64
	Mem  uint64
	Ncpu int
	....
*/

package monitor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const PROC_DIR = "/proc"

type s_proc struct {
	SysInfo
	Pids map[int]ProcPidInfo
}

var Proc = s_proc{}

//var ProcPidList = make(map[int]ProcInfo, 1000) // TODO  pid->procinfo
//var ProcExeList = make(map[string][]int, 1000) // TODO  exe -> pid
//var SystemInfo = SysInfo{}

// Refresh Proc information
func init() {
	Proc.Pids = make(map[int]ProcPidInfo, 1000)
	go func() {
		for {
			Proc.Update()
			time.Sleep(1 * time.Second)

			pids, err := Pids()
			if err != nil {
				log.Println(err)
				continue
			}
			uid := time.Now().UnixNano()
			for _, pid := range pids {
				pi := ProcPidInfo{Pid: pid, random: uid}
				pi.Update()
				Proc.Pids[pid] = pi
			}
			for pid, pinfo := range Proc.Pids {
				if pinfo.random != uid {
					delete(Proc.Pids, pid)
				}
			}
		}
	}()
}

// Get all pids
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

type ProcStat struct {
	User, Nice, Sys, Idle, Iowait, Irq, Softirq uint64
}

// Update from /proc/stst
func (s *ProcStat) Update() (st ProcStat, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	var ig string
	fmt.Fscanln(f, &ig, &s.User, &s.Nice, &s.Sys, &s.Idle, &s.Iowait, &s.Irq, &s.Softirq)
	st = *s
	return
}

type SysInfo struct {
	Cpu  float64
	Mem  uint64
	Ncpu int
	St   ProcStat
}

// Update Sysinf (include Cpu, Mem, Ncpu)
func (si *SysInfo) Update() (err error) {
	org, err := si.St.Update()
	if err != nil {
		return
	}
	//fmt.Println(org)
	time.Sleep(1 * time.Second) // sleep for 1 second
	cur, err := si.St.Update()
	if err != nil {
		return
	}
	sum := func(s *ProcStat) uint64 {
		return s.User + s.Nice + s.Sys + s.Iowait + s.Irq + s.Softirq
	}
	s1, s2 := sum(&org), sum(&cur)
	si.Cpu = float64(s2-s1) / float64(s2-s1+cur.Idle-org.Idle)
	return
}

type ProcPidStat struct {
	Pid                          int
	State                        string
	Ppid, Pgrp                   int
	Session                      int
	Utime, Stime, Cutime, Cstime uint64
}
type ProcPidInfo struct {
	Pid  int
	Exe  string
	Root string
	Cwd  string

	Cpu  float32
	Mem  uint64
	Fd   []string
	Stat ProcPidStat

	random int64
}

// Update ProcPidInfo
func (pi *ProcPidInfo) Update() (err error) {
	basedir := filepath.Join(PROC_DIR, strconv.Itoa(pi.Pid))
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

// Read from /proc/[num]/stat
func (ps *ProcPidStat) Update() (err error) {
	content, err := ioutil.ReadFile(filepath.Join(PROC_DIR, strconv.Itoa(ps.Pid), "stat"))
	if err != nil {
		return
	}
	var ig string // ignore variable

	fmt.Sscan(string(content),
		&ig, &ig, &ps.State, &ps.Ppid, &ps.Pgrp, &ps.Session,
		&ig, &ig, &ig, &ig, &ig, &ig, &ig,
		&ps.Utime, &ps.Stime, &ps.Cutime, &ps.Cstime)
	return
}
