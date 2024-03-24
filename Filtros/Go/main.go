package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	"github.com/valyala/fasthttp"
)

type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Conf   Config            `json:"Conf"`
	Start  time.Time         `json:"Start"`
	Filtro map[uint32][]byte `json:"Filtro"`
	Db     *ledis.DB         `json:"Db"`
}
type silentLogger struct{}

func (sl *silentLogger) Printf(format string, args ...interface{}) {}

func main() {

	var port string
	var dbname string
	if runtime.GOOS == "windows" {
		port = ":81"
		dbname = "C:/Diego/AllinApp/Go/Bases/BaseCompleto/LedisDB/01"
	} else {
		port = ":81"
		dbname = "/var/LedisDB/01"
	}

	handler := &MyHandler{
		Start:  time.Now(),
		Conf:   Config{Tiempo: 1 * time.Second},
		Filtro: make(map[uint32][]byte, 0),
		Db:     LedisConfig(dbname),
	}

	handler.SaveDB(100000, 1000000)

	con := context.Background()
	con, cancel := context.WithCancel(con)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					handler.Conf.init()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-con.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	Server := &fasthttp.Server{Handler: handler.HandleFastHTTP, Logger: &silentLogger{}}
	Server.ListenAndServe(port)

	if err := run(con, handler, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func (h *MyHandler) SearchMemFiltro(Id uint32) []byte {
	if Pro, foundPro := h.Filtro[Id]; foundPro {
		if Pro[0] == 0 {
			return Pro
		} else {
			return []byte{}
		}
	} else {
		val, _ := h.Db.Get(Uint32_4Bytes(Id))
		return val
	}
}

func (h *MyHandler) SearchDbFiltro(n []byte) []byte {
	/*
		val, _ := h.Db.Get(basic.Uint32ToBytes4(id))
		if len(val) > 0 {
			ctx.SetBody(val)
		}
	*/
	return []byte{0, 1}
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/busqueda1":

			var Id uint32 = Param(ctx.QueryArgs().Peek("p"))
			ctx.SetBody(h.SearchMemFiltro(Id))

		case "/busqueda2":

			var Id uint32 = Param(ctx.QueryArgs().Peek("p"))
			if Pro, foundPro := h.Filtro[Id]; foundPro {
				ctx.SetBody(Pro)
			} else {
				val, _ := h.Db.Get(Uint32_4Bytes(Id))
				ctx.SetBody(val)
			}

		case "/stats":
			fmt.Fprintf(ctx, "Stats")
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 10 * time.Second
	fmt.Printf(".")
}
func (c *Config) init() {
	var tick = flag.Duration("tick", 1*time.Second, "Ticking interval")
	c.Tiempo = *tick
}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {
	c.Conf.init()
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-con.Done():
			return nil
		case <-time.Tick(c.Conf.Tiempo):
			c.StartDaemon()
		}
	}
}

func (h *MyHandler) SaveDB(cmem, cdisk int) {
	for i := 0; i < cmem; i++ {
		h.Filtro[uint32(i)] = GetBytes(1024)
	}
	for i := cmem; i < cmem+cdisk; i++ {
		h.Db.Set(Uint32_4Bytes(uint32(i)), GetBytes(1024))
	}
}
func GetBytes(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[uint8(i)] = uint8(i % 256)
	}
	return bytes
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func LedisConfig(path string) *ledis.DB {
	cfg := lediscfg.NewConfigDefault()
	cfg.DataDir = path
	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)
	return db
}
func Uint32_4Bytes(num uint32) []byte {
	b := make([]byte, 4)
	b[0] = uint8(num / 16777216)
	b[1] = uint8(num / 65536)
	b[2] = uint8(num / 256)
	b[3] = uint8(num % 256)
	return b
}
func Param(data []byte) uint32 {
	var x uint32
	for _, c := range data {
		x = x*10 + uint32(c-'0')
	}
	return x
}
