#!/usr/bin/env python3
"""Static file server for the browser-inspector dist/ with caching disabled.

The whole app (including the multi-MB wasm) is inlined into a single
index.html, so a stale browser cache silently serves OLD code — which makes
"I rebuilt but nothing changed" bugs very confusing. Sending `Cache-Control:
no-store` on every response forces the browser to refetch the freshly built
page each load. Pair with the per-build nonce the Makefile stamps into the
HTML (<meta name="lg-build">) to confirm which build is live.

Usage: python3 serve.py [port]   (defaults to 8000; serves the current dir)
"""
import http.server
import sys


class NoCacheHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header("Cache-Control", "no-store, max-age=0")
        self.send_header("Pragma", "no-cache")
        super().end_headers()


if __name__ == "__main__":
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8000
    http.server.test(HandlerClass=NoCacheHandler, port=port)
