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
        fmt.Fprintf(w, "<h1>Hi there, I love %s! I'm <font color=red>%s</font>. I have saw that you are %s.</h1>", r.URL.Path[1:], ip, r.RemoteAddr)
}

func main() {
	cmd := exec.Command("sh", "-c", `ifconfig eth0 |grep inet|grep -v inet6|awk '{print $2}'`)
	if out, err := cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		ip = strings.Trim(string(out), "\n")
		fmt.Println(ip)
	}
	http.HandleFunc("/", handler)
        http.ListenAndServe("0.0.0.0:80", nil)
}

