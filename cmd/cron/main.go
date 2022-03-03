package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/godump/acdb"
	"github.com/godump/cron"
	"github.com/godump/doa"
	"github.com/godump/gracefulexit"
)

type Conf struct {
	VMPath string
	VM     string
	Benchs []string
}

var (
	cDb   = acdb.Doc("./db")
	cConf = func() Conf {
		conf := Conf{}
		data := doa.Try(os.ReadFile("./conf.json")).([]byte)
		doa.Nil(json.Unmarshal(data, &conf))
		return conf
	}()
)

type Item struct {
	Time     time.Time
	Duration int64
	CommitID string
}

type Main struct {
	CommitID string
}

func (m *Main) UpdateCommitID() {
	cmd := exec.Command("sh", "-c", "git rev-parse HEAD")
	cmd.Dir = cConf.VMPath
	m.CommitID = strings.TrimSpace(string(doa.Try(cmd.Output()).([]byte)))
}

func (m *Main) UpdateVM() {
	var cmd *exec.Cmd
	cmd = exec.Command("sh", "-c", "git pull origin rvv")
	cmd.Dir = cConf.VMPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	doa.Nil(cmd.Run())
	cmd = exec.Command("sh", "-c", "cargo build --release --example asm64 --features=asm")
	cmd.Dir = cConf.VMPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	doa.Nil(cmd.Run())
}

func (m *Main) CaseStorage(path string) []Item {
	r := []Item{}
	if cDb.GetDecode(filepath.Base(path), &r) != nil {
		return []Item{}
	} else {
		return r
	}
}

func (m *Main) CaseElapses(path string) int64 {
	log.Println(path)
	tic := time.Now()
	cmd := exec.Command(cConf.VM, path)
	out := doa.Try(cmd.Output()).([]byte)
	log.Println(strings.TrimSpace(string(out)))
	toc := time.Since(tic)
	return toc.Milliseconds()
}

func (m *Main) CaseOnce(path string) {
	d := m.CaseElapses(path)
	s := m.CaseStorage(path)
	s = append(s, Item{Time: time.Now(), Duration: d, CommitID: m.CommitID})
	cDb.SetEncode(filepath.Base(path), s)
}

func (m *Main) Once() {
	for _, e := range cConf.Benchs {
		m.CaseOnce(e)
	}
}

func (m *Main) Cron() {
	if cDb.Has("pm_on") && doa.Try(cDb.GetUint64("pm_on")).(uint64) != 0 {
		log.Println("main: duplicate definition")
		return
	}
	chanPing := cron.Dayz()
	chanExit := gracefulexit.Chan()
	done := 0
	for {
		select {
		case <-chanPing:
			m.UpdateVM()
			m.UpdateCommitID()
			m.Once()
		case <-chanExit:
			done = 1
		}
		if done != 0 {
			break
		}
	}
	log.Println("main: done")
}

func NewMain() *Main {
	m := &Main{}
	m.UpdateCommitID()
	return m
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "once":
		NewMain().Once()
	case "cron":
		NewMain().Cron()
	}
}
