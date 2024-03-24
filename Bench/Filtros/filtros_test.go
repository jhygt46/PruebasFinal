package bytes

import ( //"bytes"
	//"strconv"
	//"strings"
	"testing"
)

type MyHandler struct {
	Filtro map[uint64][]byte `json:"Filtro"`
}

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
