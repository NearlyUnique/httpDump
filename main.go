package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
	useColour := flag.Bool("no-colour", false, "When true use coloured output")
	curl := flag.Bool("curl", false, "When true report in curl format")
	portNum := flag.Int("port", 0, "port to listen on, overrides PORT env var")
	response := flag.String("response", `{"time":%d}`, "response, default replaces %d with time")
	flag.Parse()

	port := resolvePort(*portNum, os.Getenv("PORT"))

	color.NoColor = *useColour

	addr := ":" + port
	printf("listening on %s\n", blue(addr))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			printf("ERROR: reading body: %v", err)
		} else {
			defer r.Body.Close()
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
		}
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(*response, "%d") {
			fmt.Fprintf(w, *response, time.Now().Unix())
		} else {
			fmt.Fprint(w, *response)
		}
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

func resolvePort(p int, e string) string {
	if p == 0 && e != "" {
		n, err := strconv.Atoi(e)
		if err != nil {
			printf("basd value in env var PORT '%s'", e)
			os.Exit(1)
		}
		p = n
	}
	if p == 0 {
		p = 9000
	}
	return strconv.Itoa(p)
}
