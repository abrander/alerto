package main

import (
	"math/rand"
	"sync"
	"time"

	_ "github.com/abrander/alerto/agent/dns"
	_ "github.com/abrander/alerto/agent/http"
	_ "github.com/abrander/alerto/agent/icmpping"
	_ "github.com/abrander/alerto/agent/noop"
	_ "github.com/abrander/alerto/agent/ssh"
	"github.com/abrander/alerto/monitor"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	monitor.Loop(wg)

	wg.Wait()
}
