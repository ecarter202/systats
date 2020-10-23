package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gocraft/web"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/pkg/errors"
)

type Context struct{}

func main() {
	ctx := new(Context)

	r := web.New(Context{})
	r.Get("/system/stats", (*ctx).Stats)

	log.Fatal(http.ListenAndServe(":5555", r))
}

// Stats is the controller for getting / returning stats.
func (ctx *Context) Stats(rw web.ResponseWriter, req *web.Request) {
	stats, err := getStats()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("error getting stats : " + err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(stats)
	rw.Write(b)
}

// getStats will obtain the system usage stats.
func getStats() (stats *Stats, err error) {
	stats = new(Stats)

	memory, err := memory.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting memory stats")
	}
	stats.Memory = &MemoryUsage{
		Total:  memory.Total,
		Used:   memory.Used,
		Cached: memory.Cached,
		Free:   memory.Free,
	}

	before, err := cpu.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting cpu stats (before)")
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting cpu stats (after)")
	}
	total := float64(after.Total - before.Total)
	stats.CPU = &CPUUsage{
		User:   float64(after.User-before.User) / total * 100,
		System: float64(after.System-before.System) / total * 100,
		Idle:   float64(after.Idle-before.Idle) / total * 100,
	}

	return stats, nil
}
