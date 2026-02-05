package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

/* ========= COLORS ========= */
var (
	Red    = color("\033[1;31m%s\033[0m")
	Green  = color("\033[1;32m%s\033[0m")
	Yellow = color("\033[1;33m%s\033[0m")
	Blue   = color("\033[1;34m%s\033[0m")
	Cyan   = color("\033[1;36m%s\033[0m")
	White  = color("\033[1;97m%s\033[0m")
)

func color(c string) func(...interface{}) string {
	return func(a ...interface{}) string {
		return fmt.Sprintf(c, fmt.Sprint(a...))
	}
}

/* ========= BANNER ========= */
func banner() {
	fmt.Println(Cyan(`
██████╗ ██╗   ██╗██████╗  █████╗ ███████╗███████╗
██╔══██╗╚██╗ ██╔╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██████╔╝ ╚████╔╝ ██████╔╝███████║███████╗███████╗
██╔═══╝   ╚██╔╝  ██╔═══╝ ██╔══██║╚════██║╚════██║
██║        ██║   ██║     ██║  ██║███████║███████║
╚═╝        ╚═╝   ╚═╝     ╚═╝  ╚═╝╚══════╝╚══════╝

        403 BYPASSER | Bug Bounty Edition
        Author: algamil7x
`))
}

/* ========= HELPERS ========= */
func loadPayloads(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s", path)
	}
	lines := strings.Split(string(data), "\n")
	seen := make(map[string]bool)
	var out []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if !seen[l] {
			seen[l] = true
			out = append(out, l)
		}
	}
	return out
}

func baseline(url string) (int, int64) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode, resp.ContentLength
}

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return Green(code)
	case code >= 300 && code < 400:
		return Yellow(code)
	case code >= 400 && code < 500:
		return Red(code)
	default:
		return Blue(code)
	}
}

/* ========= TESTS ========= */

func printResult(code int, target string, length int64) {
	fmt.Printf(
		"[%s] %s | len=%d\n",
		statusColor(code),
		Green(target), // 👈 اللينك أخضر ومميّز
		length,
	)
}

func testMethods(url string, baseStatus int, baseLen int64) {
	fmt.Println(Cyan("\n[+] HTTP Method Bypass"))
	for _, m := range loadPayloads("payloads/methods") {
		req, _ := http.NewRequest(m, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != baseStatus || resp.ContentLength != baseLen {
			printResult(resp.StatusCode, m+" "+url, resp.ContentLength)
		}
	}
}

func testURLPayloads(url string, baseStatus int, baseLen int64) {
	fmt.Println(Cyan("\n[+] URL Payload Bypass"))
	for _, p := range loadPayloads("payloads/url") {
		full := url + p
		resp, err := http.Get(full)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != baseStatus || resp.ContentLength != baseLen {
			printResult(resp.StatusCode, full, resp.ContentLength)
		}
	}
}

func testHeaderNames(url string, baseStatus int, baseLen int64) {
	fmt.Println(Cyan("\n[+] Header Name Bypass"))
	for _, h := range loadPayloads("payloads/header-name") {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set(h, "127.0.0.1")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != baseStatus || resp.ContentLength != baseLen {
			printResult(resp.StatusCode, h+" → "+url, resp.ContentLength)
		}
	}
}

func testHeaderPayloads(url string, baseStatus int, baseLen int64) {
	fmt.Println(Cyan("\n[+] Header Payload Bypass"))
	for _, p := range loadPayloads("payloads/header-payload") {
		parts := strings.SplitN(p, ":", 2)
		if len(parts) != 2 {
			continue
		}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != baseStatus || resp.ContentLength != baseLen {
			printResult(resp.StatusCode, p+" → "+url, resp.ContentLength)
		}
	}
}

/* ========= MAIN ========= */
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: 403-bypasser <url>")
		os.Exit(1)
	}
	url := os.Args[1]

	banner()

	baseStatus, baseLen := baseline(url)
	fmt.Printf("[*] Baseline: %d | len=%d\n", baseStatus, baseLen)

	testMethods(url, baseStatus, baseLen)
	testURLPayloads(url, baseStatus, baseLen)
	testHeaderNames(url, baseStatus, baseLen)
	testHeaderPayloads(url, baseStatus, baseLen)
}
