package bytes

import ( //"bytes"
	//"strconv"
	//"strings"
	"testing"
)

type MyHandler struct {
	Filtro map[uint64][]byte `json:"Filtro"`
}

/*
func Benchmark_Prueba1(b *testing.B) {

	h := &MyHandler{
		Filtro: make(map[uint64][]byte, 0),
	}
	h.SaveDB(100000)

	b.ResetTimer()
	for m := 0; m < b.N; m++ {
		if Pro, foundPro := h.Filtro[75637]; foundPro {
			Silence(Pro)
		}
	}
}
func Benchmark_Prueba2(b *testing.B) {

	h := &MyHandler{
		Filtro: make(map[uint64][]byte, 0),
	}
	h.SaveDB(100000)

	b.ResetTimer()
	for m := 0; m < b.N; m++ {
		if Pro, foundPro := h.Filtro[75637]; foundPro {
			if Pro[0]%2 == 0 {
				Silence(Pro)
			}
		}
	}
}
*/

func Benchmark_Prueba1(b *testing.B) {
	a := []byte("ABCDEF")
	for m := 0; m < b.N; m++ {
		Param1(a)
	}
}
func Benchmark_Prueba2(b *testing.B) {
	a := []byte("A7s4")
	for m := 0; m < b.N; m++ {
		Param2(a)
	}
}

func Benchmark_Prueba3(b *testing.B) {

	a := []byte("/")
	b.ResetTimer()
	for m := 0; m < b.N; m++ {
		x := PathCtx(a)
		if x == 47 {

		}
	}
}
func Benchmark_Prueba4(b *testing.B) {

	a := []byte("/")
	b.ResetTimer()
	for m := 0; m < b.N; m++ {
		if len(a) == 2 {
			if a[1] == 97 {

			} else if a[1] == 98 {

			} else {

			}
		}
	}
}

func PathCtx(b []byte) uint32 {
	var x uint32
	for _, c := range b {
		x = x*256 + uint32(c)
	}
	return x
}

type Ips struct {
	IPv4_Lista  []uint32  `json:"IPv4_Lista"`
	IPv4_Lista2 [][4]byte `json:"IPv4_Lista2"`
	CountIp     int       `json:"CountIp"`
	Active      bool      `json:"Active"`
}

func (i *Ips) BuscarIpBalanceador(ip []byte) bool {

	if len(ip) == 4 {
		ips := uint32(ip[0])*16777216 + uint32(ip[1])*65536 + uint32(ip[2])*256 + uint32(ip[3])
		for _, n := range i.IPv4_Lista {
			if n == ips {
				return true
			}
		}
	} else {
		return false
	}
	return false
}
func (i *Ips) BuscarIpBalanceado2(ip []byte) bool {

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
func (i *Ips) BuscarIpBalanceador3(ip []byte) bool {

	if len(ip) == 4 {
		for _, n := range i.IPv4_Lista2 {
			if n[3] == ip[3] && n[2] == ip[2] && n[1] == ip[1] && n[0] == ip[0] {
				return true
			}
		}
	} else {
		return false
	}
	return false
}

func Param(data []byte) uint32 {
	var x uint32
	for _, c := range data {
		x = x*10 + uint32(c-'0')
	}
	return x
}
func Param1(data []byte) uint32 {
	var x uint32
	for _, c := range data {
		x = x*26 + uint32(c-'A')
	}
	return x
}
func Param2(data []byte) uint32 {
	var x uint32
	for _, c := range data {
		if c > 64 && c < 91 {
			x = x*26 + uint32(c-'A')
		} else if c > 96 && c < 123 {
			x = x*26 + uint32(c-'a')
		} else {
			x = x*26 + uint32(c-'0')
		}
	}
	return x
}

func Silence(b []byte) {

}
func (h *MyHandler) SaveDB(cmem int) {
	b := GetBytes(1024)
	for i := 0; i < cmem; i++ {
		h.Filtro[uint64(i)] = b
	}
}
func GetBytes(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[uint8(i)] = uint8(i % 256)
	}
	return bytes
}
