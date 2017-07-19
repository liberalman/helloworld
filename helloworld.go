package main

import (
        "fmt"
        "os/exec"
        "net/http"
        "os"
        "strings"
)

var ip string

func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<h1>Hello World! I'm %s. I saw that you are %s.</h1>", ip, r.RemoteAddr)
}

func main() {
	ip = os.Getenv("MY_HOST")
	cmd := exec.Command("sh", "-c", `ifconfig eth0 |grep inet|grep -v inet6|awk '{print $2}'`)
	if out, err := cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		ip = "<font color=blue>" + ip + "</font> <font color=red>" + strings.Trim(string(out), "\n") + "</font>"
		fmt.Println(ip)
	}
	http.HandleFunc("/", handler)
        http.ListenAndServe("0.0.0.0:80", nil)
}

