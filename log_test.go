package log

import (
	"testing"
)

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()
	//f, _ := os.Open(os.DevNull)
	//New(f, "", log.Lshortfile)
	NewRemoteSyslog("tcp", "172.31.2.43:55140", "app-log", 0)
	for i := 0; i < b.N; i++ {
		Warn("123")
	}
}

func TestFluentLog(t *testing.T) {
	NewRemoteSyslog("udp", "172.31.2.43:55140", "app-log", 0)
	//New(os.Stdout, "", log.Lmicroseconds)
	Warn("aa")
	/* for {
		Debug("aa")
		time.Sleep(5 * time.Second)
	} */
}
func TestLocal(t *testing.T) {
	//NewRemoteSyslog("udp", "172.31.2.43:55140", "app-log", 0)
	//New(os.Stdout, "", log.Lmicroseconds)
	//ReleasePool()
	//Warn("aa")
	/* for {
		Debug("aa")
		time.Sleep(5 * time.Second)
	} */
}
