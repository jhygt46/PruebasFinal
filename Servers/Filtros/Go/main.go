package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
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
	Ips    *Ips              `json:"Ips"`
	Conf   Config            `json:"Conf"`
	Start  time.Time         `json:"Start"`
	Filtro map[uint64][]byte `json:"Filtro"`
	Db     *ledis.DB         `json:"Db"`
}
type Ips struct {
	IPv4_Lista []uint32 `json:"IPv4_Lista"`
	CountIp    int      `json:"CountIp"`
	Active     bool     `json:"Active"`
}

type silentLogger struct{}

func (sl *silentLogger) Printf(format string, args ...interface{}) {}

func main() {

	var port string
	var dbname string
	if runtime.GOOS == "windows" {
		port = ":82"
		dbname = "C:/Diego/AllinApp/Go/Bases/BaseCompleto/LedisDB/01"
	} else {
		port = ":82"
		dbname = "/var/LedisDB/Filtros"
	}

	handler := &MyHandler{
		Start:  time.Now(),
		Conf:   Config{Tiempo: 1 * time.Second},
		Filtro: make(map[uint64][]byte, 0),
		Db:     LedisConfig(dbname),
		Ips:    &Ips{Active: false, IPv4_Lista: []uint32{167323543, 167323549, 167323542, 167323546}, CountIp: 4},
	}

	//handler.SaveDB(10, 0)

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

	go func() {
		server1 := &fasthttp.Server{Handler: handler.HandleFastHTTP, Logger: &silentLogger{}}
		server1.ListenAndServe(port)
	}()
	if err := run(con, handler, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	if !h.Ips.Active || h.Ips.BuscarIpBalanceador(ctx.RemoteIP()) {

		if len(ctx.Path()) == 2 {
			if ctx.Path()[1] == 97 {
				ctx.SetBody(h.SearchMemFiltro(ParamAlpha(ctx.QueryArgs().Peek("p")), 0))
			} else if ctx.Path()[1] == 98 {

				p := ParamAlpha(ctx.QueryArgs().Peek("p"))
				if p == 52 {
					// MEMORY
					ExecMonitoring("memory")
					ctx.SetBody([]byte{65})
				} else if p == 53 {
					ExecMonitoring("disk")
					ctx.SetBody([]byte{65})
				}

			}
		}
	} else {
		// ALERT IP NOT FOUND
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 10 * time.Second
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

// DAEMON //

// BASIC //
func LedisConfig(path string) *ledis.DB {
	cfg := lediscfg.NewConfigDefault()
	cfg.DataDir = path
	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)
	return db
}
func Uint64_NBytes(num uint64, c int) []byte {
	b := make([]byte, 6)
	b[0] = uint8(num / 1099511627776)
	b[1] = uint8(num / 4294967296)
	b[2] = uint8(num / 16777216)
	b[3] = uint8(num / 65536)
	b[4] = uint8(num / 256)
	b[5] = uint8(num % 256)
	return b[:c]
}
func ParamAlpha(data []byte) uint64 {
	var x uint64
	for _, c := range data {
		if c > 64 && c < 91 {
			x = x*64 + uint64(c-65)
		} else if c > 96 && c < 123 {
			x = x*64 + uint64(c-71)
		} else {
			x = x*64 + uint64(c+4)
		}
	}
	return x
}
func Join(c1, c2 []byte) []byte {
	r := make([]byte, len(c1)+len(c2))
	copy(r, c1)
	copy(r[len(c1):], c2)
	return r
}
func (i *Ips) BuscarIpBalanceador(ip net.IP) bool {

	if len(ip) == 4 {
		ips := uint32(ip[0])*16777216 + uint32(ip[1])*65536 + uint32(ip[2])*256 + uint32(ip[3])
		for j := 0; j < i.CountIp; j++ {
			if i.IPv4_Lista[j] == ips {
				return true
			}
		}
	} else {
		return false
	}
	return false
}
func ExecMonitoring(x string) {

	fmt.Println("Montioring", x)

	cmd := exec.Command("monitoring", x)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error al ejecutar el programa:", err)
		return
	}

	result := strings.TrimSpace(string(output))
	fmt.Println(output)

	// Mostrar la información devuelta por el programa
	fmt.Println("Información devuelta por el programa en Go:", result)

}

// BASIC //

// FILTRO //
func (h *MyHandler) SearchMemFiltro(Id, x uint64) []byte {

	var c int = 6
	var num uint64 = Id % 4294967296

	if x > 0 {
		c = 4
		num = Id % 281474976710656
		if x > 1 {
			return []byte{}
		}
	}
	x++

	if Per, foundPro := h.Filtro[num]; foundPro {
		fmt.Println("Per:", Per)
		if Per[0]%2 == 0 {
			return Per
		} else {
			nnum, n := GetNuevoNum(Per, num)
			return Join(Per[n:], h.SearchMemFiltro(nnum, x))
		}
	} else {
		val, _ := h.Db.Get(Uint64_NBytes(num, c))
		if val[0]%2 == 0 {
			return val
		} else {
			nnum, n := GetNuevoNum(val, num)
			return Join(val[n:], h.SearchMemFiltro(nnum, x))
		}
	}
}
func GetNuevoNum(b []byte, num uint64) (uint64, int) {

	fmt.Println("GetNuevoNum")
	var len int = len(b)
	var i int = 0
	var nnum uint64
	for {
		if len < i+1 {
			break
		}
		fmt.Printf("%v -", b[i])
		i++
	}
	return nnum, i
}

// FILTRO //

// TEST ENCODE DECODE //
type Filtro struct {
	Tipo    byte      `json:"Tipo"`
	Filtros []Filtros `json:"Filtros"`
	Evals   []string  `json:"Evals"`
}
type Filtros struct {
	Tipo    byte            `json:"Tipo"`
	Numero  int             `json:"Numero"`
	Nombre  string          `json:"Nombre"`
	Valores []FiltroValores `json:"Valores"`
	Order   int             `json:"Order"`
}
type FiltroValores struct {
	Nombre  string `json:"Nombre"`
	Numero  int    `json:"Numero"`
	Padre   int    `json:"Padre"`
	Element int    `json:"Element"`
	Order   int    `json:"Order"`
}

func (h *MyHandler) SaveDB(cmem, cdisk int) {

	Filtros := []Filtros{
		Filtros{Tipo: 0, Order: 1, Numero: 0, Nombre: "Marca", Valores: []FiltroValores{FiltroValores{Numero: 1, Nombre: "Audi", Order: 2, Padre: 1, Element: 1}, FiltroValores{Numero: 2, Nombre: "Volvo", Order: 1, Padre: 1, Element: 1}, FiltroValores{Numero: 3, Nombre: "Chevrolet", Order: 0, Padre: 1, Element: 1}}},
		Filtros{Tipo: 0, Order: 0, Numero: 1, Nombre: "Modelo", Valores: []FiltroValores{FiltroValores{Numero: 1, Nombre: "A3", Order: 0, Padre: 1, Element: 1}}},
	}
	Evals := []string{"Eval1", "Eval2", "Eval3", "Eval4", "Eval5", "Eval6"}
	Filtro := Filtro{Tipo: 0, Filtros: Filtros, Evals: Evals}

	b := EncodeFiltro(Filtro)
	fmt.Println(b)

	for i := 0; i < cmem; i++ {
		h.Filtro[uint64(i)] = GetBytes(1024)
	}
	for i := cmem; i < cmem+cdisk; i++ {
		h.Db.Set(Uint64_NBytes(uint64(i), 0), GetBytes(1024))
	}
}
func EncodeFiltro(x Filtro) []byte {

	var buf []byte
	var bufnom []byte
	var bufopc []byte
	var num uint8 = 0
	var posnom int = 0

	buf = append(buf, x.Tipo)

	for i, filtro := range x.Filtros {

		fmt.Println(filtro.Nombre, i, num)

		buf = append(buf, filtro.Tipo)
		buf = append(buf, NuMmax2Bytes(filtro.Numero)...)
		buf = append(buf, NuMmax2Bytes(filtro.Order)...)

		if x.Tipo == 0 {
			buf = append(buf, uint8(len(filtro.Nombre)))
			buf = append(buf, []byte(filtro.Nombre)...)
		} else {
			posnom = GetFiltrosNombres(&bufnom, filtro.Nombre)
			buf = append(buf, NuMmax2Bytes(posnom)...)
		}

		for _, valores := range filtro.Valores {

			bufopc = append(bufopc, NuMmax2Bytes(valores.Numero)...)
			bufopc = append(bufopc, NuMmax2Bytes(valores.Order)...)
			bufopc = append(bufopc, NuMmax2Bytes(valores.Padre)...)

			if x.Tipo == 0 {
				bufopc = append(bufopc, uint8(len(valores.Nombre)))
				bufopc = append(bufopc, []byte(valores.Nombre)...)
			} else {
				posnom = GetFiltrosNombres(&bufnom, valores.Nombre)
				bufopc = append(bufopc, NuMmax2Bytes(posnom)...)
			}
		}

	}
	buf = append(buf, bufnom...)

	return buf
}
func NuMmax2Bytes(n int) []byte {
	b := make([]byte, 2)
	if n < 200 {
		b[0] = uint8(n)
		return b[:1]
	} else {
		b[0] = uint8(n/256 + 200)
		b[1] = uint8(n%256 - 200)
		return b
	}
}
func GetFiltrosNombres(b *[]byte, nom string) int {

	var i int = 0
	var len int = len(*b)
	var pos int = 0

	if len == 0 {
		return 0
	}

	for {

		nameLen := int((*b)[i])
		currentName := string((*b)[i+1 : i+1+nameLen])

		if nom == currentName {
			return pos
		} else {
			pos++
		}

		i += nameLen + 1
		if len == i {
			//*b = append(*b, uint8(len(nom)))
			*b = append(*b, []byte(nom)...)
			return pos
		}

	}
}
func GetBytes2(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i <= size; i++ {
		bytes[uint8(i)] = uint8(i % 256)
	}
	return bytes
}
func GetBytes(size int) []byte {
	return []byte{1, 2, 3, 4, 5}
}

// TEST ENCODE DECODE //
