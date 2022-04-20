package log

import (
	"log"
	"os"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()
	f, _ := os.Open(os.DevNull)
	New(f, "", log.Lshortfile)
	for i := 0; i < b.N; i++ {
		Debug("123")
	}
}
