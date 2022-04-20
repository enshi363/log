package log

import (
	"log"
	"os"
	"testing"
	"time"
)

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()
	f, _ := os.Open(os.DevNull)
	New(f, "", log.Lshortfile)
	for i := 0; i < b.N; i++ {
		Debug("123")
	}
}

func TestFluentLog(t *testing.T) {
	NewRemoteSyslog("tcp", "172.31.2.43:55140", "app-log")
	//New(os.Stdout, "", log.Lmicroseconds)
	//Warn("aa")
	for {
		Debug("aa")
		time.Sleep(5 * time.Second)
	}
}
