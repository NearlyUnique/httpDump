package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	bright  = "\033[1m"
	clear   = "\033[0m"
	black   = 30
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	magenta = 35
	cyan    = 36
	white   = 37
)

func colour(c int, text string) string {
	return fmt.Sprintf("\033[%dm%s%s", c, text, clear)
}
func colourB(c int, text string) string {
	return fmt.Sprintf("\033[%dm%s%s%s", c, bright, text, clear)
}

func mapToStringSlice(headers map[string][]string) string {
	s := ""
	for k, v := range headers {
		s += fmt.Sprintf("%s = %s\n", k, strings.Join(v, ","))
	}
	return s
}

func main() {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "9000"
	}
	addr := ":" + port
	fmt.Printf("listening on %s\n", colour(blue, addr))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("-------------------\nrequest:\n%s %s\n\nheaders:\n%v\nbody:\n%s\n",
			colour(yellow, r.Method),
			colour(green, r.URL.String()),
			colourB(magenta, mapToStringSlice(r.Header)),
			colour(cyan, string(buf)),
		)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"time\":\"%d\"}", time.Now().Unix())
	})
	http.ListenAndServe(addr, nil)
}
