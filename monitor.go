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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
    "sync"
)

var refreshLock = sync.Mutex{}

const PROC_DIR = "/proc"

type s_proc struct {
	sysInfo
	Pids map[int]procPidInfo
}

var ErrorNotExists = errors.New("pid not exists")

var Proc = s_proc{
	Pids: make(map[int]procPidInfo, 100),
}
var Interval time.Duration = time.Second * 1

func Pid(pid int) (pi procPidInfo, err error) {
	pi, ok := Proc.Pids[pid]
	if !ok {
		err = ErrorNotExists
	}
	return
}

func Ncpu() int {
	return Proc.Ncpu
}

func Cpu() float64 {
	return Proc.Cpu
}

func Mem() uint64 {
	return Proc.Mem
}

// initial the proc stat
func init() {
	go Refresh()
}

func readFile(filename string) (data []byte, err error) {
	return exec.Command("/bin/cat", filename).Output()
}


// refresh proc states
func Refresh() {
    refreshLock.Lock()
    defer refreshLock.Unlock()

	Proc.Update()

	pids, err := Pids()
	if err != nil {
		log.Println(err)
		return
	}
	uid := time.Now().UnixNano()
	for _, pid := range pids {
		pi := procPidInfo{Pid: pid, random: uid}
		pi.Update()
		Proc.Pids[pid] = pi
	}
	// clean not existed pids
	for pid, pinfo := range Proc.Pids {
		if pinfo.random != uid {
			delete(Proc.Pids, pid)
		}
	}
}

func Hostname() (name string, err error) {
	data, err := readFile("/proc/sys/kernel/hostname")
	if err != nil {
		return
	}
	return strings.TrimSpace(string(data)), nil
}

// Refresh Proc information in a gorountine
func GoRefresh() {
	go func() {
		for {
			Refresh()
			time.Sleep(Interval)
		}
	}()
}

func ls(dir string) (names []string, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	defer f.Close()
	return f.Readdirnames(-1)
}

// Get all pids
func Pids() ([]int, error) {
	names, err := ls("/proc")
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

type procStat struct {
	User, Nice, Sys, Idle, Iowait, Irq, Softirq uint64
	Ncpu                                        int
}

// Update from /proc/stst
func (s *procStat) Update() (st procStat, err error) {
	data, err := readFile(filepath.Join(PROC_DIR, "stat"))
	if err != nil {
		log.Println(err)
		return
	}
	var ig string
	fmt.Sscanln(string(data), &ig, &s.User, &s.Nice, &s.Sys, &s.Idle, &s.Iowait, &s.Irq, &s.Softirq)
	st = *s
	// Cpu Count
	re, _ := regexp.Compile("\ncpu[0-9]")
	s.Ncpu = len(re.FindAllString(string(data), 100))
	return
}

type sysInfo struct {
	Cpu  float64
	Mem  uint64
	Ncpu int
	St   procStat
}

// Update Sysinf (include Cpu, Mem, Ncpu)
func (si *sysInfo) Update() (err error) {
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
	sum := func(s *procStat) uint64 {
		return s.User + s.Nice + s.Sys + s.Iowait + s.Irq + s.Softirq
	}
	s1, s2 := sum(&org), sum(&cur)
	si.Cpu = float64(s2-s1) / float64(s2-s1+cur.Idle-org.Idle)
	si.Ncpu = si.St.Ncpu
	return
}

type procPidStat struct {
	Pid                          int
	State                        string
	Ppid, Pgrp                   int
	Session                      int
	Utime, Stime, Cutime, Cstime uint64
}
type procPidInfo struct {
	Pid  int
	Exe  string
	Root string
	Cwd  string

	Cpu  float32
	Mem  uint64
	Fd   []string
	Stat procPidStat

	random int64
}

// Update procPidInfo
func (pi *procPidInfo) Update() (err error) {
	basedir := filepath.Join(PROC_DIR, strconv.Itoa(pi.Pid))
	pi.Exe, _ = os.Readlink(filepath.Join(basedir, "exe"))
	pi.Cwd, _ = os.Readlink(filepath.Join(basedir, "cwd"))
	pi.Root, _ = os.Readlink(filepath.Join(basedir, "root"))

	// FIXME: finish it
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
func (ps *procPidStat) Update() (err error) {
	content, err := readFile(filepath.Join(PROC_DIR, strconv.Itoa(ps.Pid), "stat"))
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
