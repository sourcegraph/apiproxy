package main_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Example() {
	// Start a target server.
	targetRequestCount := 0
	targetResponseBody := []byte("qux")
	targetMux := http.NewServeMux()
	targetMux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		targetRequestCount++
		w.Header().Add("Cache-Control", "max-age=60")
		w.Write(targetResponseBody)
	})
	target := httptest.NewServer(targetMux)
	defer target.Close()

	// Start apiproxy.
	cmd := exec.Command(program, "-http=:8090", "-never-revalidate", "-only-revalidate-older-than=24h", target.URL)
	cmd.Start()
	defer cmd.Process.Kill()
	time.Sleep(150 * time.Millisecond)

	// Hit the target server via apiproxy.
	httpGet("http://localhost:8090/foo")
	httpGet("http://localhost:8090/foo")

	// Output:
	// Got response: qux (cached: false)
	// Got response: qux (cached: true)
}

func httpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting /foo: %s\n", err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %s\n", err)
	}

	fmt.Printf("Got response: %s (cached: %v)\n", data, resp.Header.Get("X-From-Cache") != "")
}

var program string

func init() {
	var err error

	// The executable name will be the directory name.
	if program, err = os.Getwd(); err != nil {
		log.Fatal(err)
	}
	program = filepath.Base(program)

	if _, err = exec.LookPath(program); err != nil {
		if err.(*exec.Error).Err == exec.ErrNotFound {
			if err = exec.Command("go", "install").Run(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}
