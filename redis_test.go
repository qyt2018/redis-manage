package main

import(
	"testing"
	"net"
	"fmt"
)

var (
	test_client net.Conn
)

func init(){
	redisServer := &server{"", "", "", ""}
	test_client, _ = conn(redisServer)
	defer test_client.Close()
}

func Testexec(t *testing.T){
	cmd := ""
	out, erro := exec(cmd, test_client)
	if erro != nil {
		t.Errorf("exec command: %s", erro.Error())
		return
	}
	want := "hash"
	if(out[0] != want){
		t.Errorf("cmd:%s, want is %s, now is %s",cmd, want, out[0])
	}
}

func TestgetType(t *testing.T){
	key := ""
	out, _ := getType(key, test_client)
	if out == "hash" {
		t.Errorf("key:%s, want is hash, now is %s",key, out)
	}
}

func BenchmarkgetType(b *testing.B){
	for i := 0; i < b.N; i++ {
		getType("", test_client)
	}
}

func ExamplegetType(){
	out, _ := getType("", test_client)
	fmt.Println(out)
}