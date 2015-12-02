package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/abrander/alerto/api"
	_ "github.com/abrander/alerto/config"
	"github.com/abrander/alerto/monitor"
	_ "github.com/abrander/alerto/plugins/dns"
	_ "github.com/abrander/alerto/plugins/http"
	_ "github.com/abrander/alerto/plugins/icmpping"
	_ "github.com/abrander/alerto/plugins/localtransport"
	_ "github.com/abrander/alerto/plugins/noop"
	_ "github.com/abrander/alerto/plugins/pidof"
	_ "github.com/abrander/alerto/plugins/ssh"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go api.Run(wg)

	wg.Add(1)
	monitor.Loop(wg)

	wg.Wait()
}
