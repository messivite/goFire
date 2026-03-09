package handlers

import (
	"net/http"

	"github.com/mustafaaksoy/goFire/config"
)

func Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GoFire</title>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@600;700&display=swap" rel="stylesheet">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #0d1117;
      font-family: 'JetBrains Mono', monospace;
      color: #58a6ff;
      overflow: hidden;
    }
    .container {
      text-align: center;
      animation: fadeIn 0.8s ease-out;
    }
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(20px); }
      to   { opacity: 1; transform: translateY(0); }
    }
    .box {
      display: inline-block;
      padding: 2.5rem 3.5rem;
      border: 2px solid #58a6ff;
      border-radius: 6px;
      position: relative;
    }
    .box::before {
      content: '';
      position: absolute;
      inset: -1px;
      border-radius: 6px;
      background: linear-gradient(135deg, rgba(88,166,255,0.15), transparent);
      z-index: -1;
    }
    .fire { font-size: 3rem; margin-bottom: 0.5rem; }
    .brand {
      font-size: 2.8rem;
      font-weight: 700;
      color: #fff;
      letter-spacing: 0.15em;
      text-shadow: 0 0 30px rgba(88,166,255,0.4);
    }
    .version {
      margin-top: 0.75rem;
      font-size: 0.85rem;
      color: #8b949e;
    }
    .links {
      margin-top: 1.8rem;
      display: flex;
      gap: 1.5rem;
      justify-content: center;
      font-size: 0.8rem;
    }
    .links a {
      color: #58a6ff;
      text-decoration: none;
      padding: 0.4rem 0.8rem;
      border: 1px solid #30363d;
      border-radius: 4px;
      transition: all 0.2s;
    }
    .links a:hover {
      border-color: #58a6ff;
      background: rgba(88,166,255,0.1);
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="box">
      <div class="brand">GoFire</div>
      <div class="version">v` + config.Version + `</div>
      <div class="links">
        <a href="/api/health">/api/health</a>
      </div>
    </div>
  </div>
</body>
</html>`

	w.Write([]byte(html))
}
