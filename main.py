#!/usr/bin/env python3
import requests
import sys
from pathlib import Path

# ========= COLORS =========
GREEN  = "\033[1;32m"
RED    = "\033[1;31m"
YELLOW = "\033[1;33m"
BLUE   = "\033[1;34m"
CYAN   = "\033[1;36m"
RESET  = "\033[0m"

def c(color, text):
    return f"{color}{text}{RESET}"

# ========= BANNER =========
def banner():
    print(c(CYAN, r"""
██████╗ ██╗   ██╗██████╗  █████╗ ███████╗███████╗
██╔══██╗╚██╗ ██╔╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██████╔╝ ╚████╔╝ ██████╔╝███████║███████╗███████╗
██╔═══╝   ╚██╔╝  ██╔═══╝ ██╔══██║╚════██║╚════██║
██║        ██║   ██║     ██║  ██║███████║███████║
╚═╝        ╚═╝   ╚═╝     ╚═╝  ╚═╝╚══════╝╚══════╝

        403 BYPASSER | Python Edition
        Author: algamil7x
"""))

# ========= HELPERS =========
def load_payloads(name):
    path = Path("payloads") / name
    if not path.exists():
        return []
    return [l.strip() for l in path.read_text().splitlines()
            if l.strip() and not l.startswith("#")]

def status_color(code):
    if 200 <= code < 300:
        return GREEN
    if 300 <= code < 400:
        return YELLOW
    if 400 <= code < 500:
        return RED
    return BLUE

def print_result(code, target, length):
    print(f"[{c(status_color(code), code)}] {c(GREEN, target)} | len={length}")

# ========= HTTP =========
session = requests.Session()
session.headers.update({
    "User-Agent": "403-bypasser/1.0"
})
TIMEOUT = 10

def request(method, url, headers=None):
    try:
        r = session.request(
            method=method,
            url=url,
            headers=headers,
            timeout=TIMEOUT,
            allow_redirects=False,
            verify=True
        )
        return r.status_code, len(r.content)
    except Exception:
        return None, None

# ========= MAIN LOGIC =========
def main():
    if len(sys.argv) != 2:
        print("Usage: python3 main.py <url>")
        sys.exit(1)

    url = sys.argv[1].rstrip("/")
    banner()

    # ---- Baseline ----
    base_status, base_len = request("GET", url)
    if base_status is None:
        print(c(RED, "[!] Baseline failed (network / DNS issue)"))
        return

    print(f"[*] Baseline: {base_status} | len={base_len}")

    # ---- URL Payloads ----
    print(c(CYAN, "\n[+] URL Payload Bypass"))
    for p in load_payloads("url"):
        full = url + p
        code, length = request("GET", full)
        if code is None:
            continue
        if code != base_status or length != base_len:
            print_result(code, full, length)

    # ---- HTTP Methods ----
    print(c(CYAN, "\n[+] HTTP Method Bypass"))
    for m in load_payloads("methods"):
        code, length = request(m, url)
        if code is None:
            continue
        if code != base_status or length != base_len:
            print_result(code, f"{m} {url}", length)

    # ---- Header Names ----
    print(c(CYAN, "\n[+] Header Name Bypass"))
    for h in load_payloads("header-name"):
        code, length = request("GET", url, headers={h: "127.0.0.1"})
        if code is None:
            continue
        if code != base_status or length != base_len:
            print_result(code, f"{h} → {url}", length)

    # ---- Header Payloads ----
    print(c(CYAN, "\n[+] Header Payload Bypass"))
    for hp in load_payloads("header-payload"):
        if ":" not in hp:
            continue
        k, v = hp.split(":", 1)
        code, length = request("GET", url, headers={k.strip(): v.strip()})
        if code is None:
            continue
        if code != base_status or length != base_len:
            print_result(code, f"{hp} → {url}", length)

if __name__ == "__main__":
    main()

