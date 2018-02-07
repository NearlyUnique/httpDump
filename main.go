package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	blue      = color.New(color.FgBlue).SprintFunc()
	yellow    = color.New(color.FgYellow).SprintFunc()
	green     = color.New(color.FgGreen).SprintFunc()
	cyan      = color.New(color.FgCyan).SprintFunc()
	magentaBg = color.New(color.BgMagenta).SprintFunc()
	magentaFg = color.New(color.FgMagenta).SprintFunc()
)

func mapToStringSlice(headers map[string][]string) string {
	s := ""
	for k, v := range headers {
		s += fmt.Sprintf("%s = %s\n", k, strings.Join(v, ","))
	}
	return s
}

func main() {
	pretty := flag.Bool("pretty", true, "When true pretty print")
	useColor := flag.Bool("color", true, "When true use coloured output")
	curl := flag.Bool("curl", false, "When true report in curl format")
	flag.Parse()

	port := os.Getenv("PORT")

	color.NoColor = !*useColor

	if len(port) == 0 {
		port = "9000"
	}
	addr := ":" + port
	printf("listening on %s\n", blue(addr))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		body := string(buf)
		if *pretty {
			printf("-------------------\nrequest:\n%s %s\n\nheaders:\n%v\nbody:\n%s\n",
				yellow(r.Method),
				green(r.URL.String()),
				magentaBg(mapToStringSlice(r.Header)),
				cyan(body),
			)
		}
		if *curl {
			printCurl(r, body)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"time\":\"%d\"}", time.Now().Unix())
	})
	http.ListenAndServe(addr, nil)
}

func printCurl(r *http.Request, body string) {
	printf(blue("curl"))
	printf("   -X %s \\\n", yellow(r.Method))
	for k, v := range r.Header {
		for _, vi := range v {
			printf("  -H \"%s: %s\" \\\n", magentaBg(k), magentaFg(vi))
		}
	}
	if len(body) > 0 {
		printf("   -d '%s' \\\n", body)
	}
	printf("   \"https://%s%s\" \n", r.Host, green(r.URL.String()))
}

func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(color.Output, format, a...)
}
