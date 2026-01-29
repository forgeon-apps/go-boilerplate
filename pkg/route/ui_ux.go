package route

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterUI(route fiber.Router) {
	route.Get("/ui", UIHomePage)

	route.Get("/ui/tasks", UITasksPage)
	route.Get("/ui/users", UIUsersPage)
	route.Get("/ui/projects", UIProjectsPage)
	route.Get("/ui/books", UIBooksPage)
}

// ------------------------------------------------------------
// HTML shell + shared helpers
// ------------------------------------------------------------

// IMPORTANT:
// We DO NOT use fmt.Sprintf for the shell template because CSS/JS contains lots of '%' like '0%' '100%'
// which breaks fmt formatting and produces "%!s(MISSING)" artifacts.
const htmlShellTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>{{TITLE}}</title>
  <meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover" />
  <meta name="color-scheme" content="dark" />
  <style>
    :root{
      color-scheme: dark;

      /* Core palette */
      --bg:#000;
      --card:#0a0a0a;
      --panel:#070707;

      --border:#1f1f1f;
      --border2:#2c2c2c;

      --text:#f5f5f5;
      --muted:#9ca3af;
      --accent:#ffffff;

      --dim:#6b7280;

      --shadow: 0 18px 55px rgba(0,0,0,.72);
      --radius: 1.1rem;
      --pad: clamp(.85rem, 2.2vw, 1.25rem);

      /* Grid */
      --grid-size: 28px;
      --grid-line: rgba(255,255,255,.10);
      --grid-glow: rgba(255,255,255,.05);
    }

    *{box-sizing:border-box;margin:0;padding:0}
    html{height:100%}
    body{min-height:100%}

    /* Page */
    body{
      font-family:system-ui,-apple-system,BlinkMacSystemFont,"SF Pro Text",sans-serif;
      color:var(--text);
      padding:
        calc(var(--pad) + env(safe-area-inset-top))
        calc(var(--pad) + env(safe-area-inset-right))
        calc(var(--pad) + env(safe-area-inset-bottom))
        calc(var(--pad) + env(safe-area-inset-left));
      overflow-x:hidden;
      background: var(--bg);
      position: relative;
    }

    /* Grid overlay */
    body::before{
      content:"";
      position: fixed;
      inset: 0;
      pointer-events: none;
      z-index: 0;
      background:
        radial-gradient(800px circle at 18% -10%, var(--grid-glow), transparent 60%),
        radial-gradient(700px circle at 86% -6%,  var(--grid-glow), transparent 65%),
        repeating-linear-gradient(
          90deg,
          var(--grid-line) 0px,
          var(--grid-line) 1px,
          transparent 1px,
          transparent var(--grid-size)
        ),
        repeating-linear-gradient(
          0deg,
          var(--grid-line) 0px,
          var(--grid-line) 1px,
          transparent 1px,
          transparent var(--grid-size)
        );
      opacity: .95;
    }

    /* Layout wrapper */
    .wrap{
      width:100%;
      max-width:1100px;
      margin:0 auto;
      position: relative;
      z-index: 1;
    }

    .card{
      border-radius:var(--radius);
      border:1px solid var(--border);
      background:
        radial-gradient(circle at top left, rgba(255,255,255,.06) 0, rgba(255,255,255,0) 55%),
        radial-gradient(circle at bottom right, rgba(255,255,255,.04) 0, rgba(255,255,255,0) 60%),
        var(--card);
      box-shadow: var(--shadow);
      padding: clamp(.9rem, 2vw, 1.1rem);
      position:relative;
      overflow:hidden; /* keep everything clipped */
    }

    /* Floating Home */
    .home-fab{
      position: sticky;
      top: calc(.25rem + env(safe-area-inset-top));
      z-index: 10;
      display:inline-flex;
      align-items:center;
      gap:.5rem;
      padding:.45rem .65rem;
      border-radius:999px;
      border:1px solid var(--border);
      background: rgba(0,0,0,.78);
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
      box-shadow: 0 10px 24px rgba(0,0,0,.55);
      color: var(--accent);
      text-decoration:none;
      margin-bottom:.75rem;
      width: fit-content;
    }
    .home-fab:hover{background: rgba(0,0,0,.9); border-color: var(--border2)}
    .home-fab .ico svg{stroke: var(--accent)}

    /* Header block */
    .top{
      display:grid;
      grid-template-columns: 1fr 360px;
      gap:1rem;
      align-items:start;
      margin-bottom:.9rem;
    }

    .eyebrow{
      font-size:.7rem;
      letter-spacing:.22em;
      text-transform:uppercase;
      color:var(--muted);
      margin-bottom:.45rem
    }

    h1{
      font-size:1.35rem;
      line-height:1.15;
      margin-bottom:.4rem;
      letter-spacing:-.01em;
    }
    p{font-size:.92rem;line-height:1.6;color:var(--muted)}

    /* Rainbow headline */
    .rainbow{
      background: linear-gradient(90deg,
        #ff3b3b, #ffb13b, #fff13b, #3bff7a, #3bbcff, #7a3bff, #ff3bbf, #ff3b3b);
      background-size: 300% 100%;
      -webkit-background-clip: text;
      background-clip: text;
      color: transparent;
      animation: rainbowMove 7s linear infinite;
    }
    @keyframes rainbowMove{
      0%{background-position:0% 50%}
      100%{background-position:100% 50%}
    }
    @media (prefers-reduced-motion: reduce){
      .rainbow{animation:none}
    }

    /* Stack badges */
    .stack{display:flex;align-items:center;gap:.55rem;margin-top:.65rem;flex-wrap:wrap}
    .stack-badge{
      display:inline-flex;align-items:center;gap:.45rem;
      border:1px solid var(--border);
      background: rgba(0,0,0,.60);
      border-radius:999px;
      padding:.28rem .6rem;
      font-size:.78rem;
      color:var(--muted);
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
    }
    .stack-badge svg{width:16px;height:16px;display:block}
    .stack-badge strong{color:var(--accent);font-weight:600}

    /* Pills */
    .pill-row{display:flex;flex-wrap:wrap;gap:.45rem;margin-top:.75rem}
    .pill{
      font-size:.7rem;text-transform:uppercase;letter-spacing:.16em;
      padding:.22rem .55rem;border-radius:999px;border:1px solid var(--border);
      color:var(--muted);display:inline-flex;gap:.4rem;align-items:center;
      background: rgba(0,0,0,.30);
    }
    .pill strong{color:var(--accent);font-weight:600}

    /* Panel */
    .panel{
      border:1px solid var(--border);
      background: rgba(0,0,0,.55);
      border-radius:1rem;
      padding:.9rem;
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
      overflow:hidden;
    }

    .actions{
      display:flex;
      gap:.5rem;
      align-items:center;
      justify-content:space-between;
      flex-wrap:wrap;
      min-width: 0;
    }

    button{
      cursor:pointer;
      border:1px solid var(--border);
      background:#000;
      color:var(--accent);
      border-radius:.8rem;
      padding:.55rem .85rem;
      font-size:.88rem;
      transition: transform .08s ease, background .12s ease, border-color .12s ease;
      white-space:nowrap;
    }
    button:hover{background:#0b0b0b;border-color:var(--border2)}
    button:active{transform: translateY(1px)}

    .hint{font-size:.75rem;color:var(--dim);margin-top:.55rem}

    input{
      width:100%;
      margin-top:.45rem;
      padding:.62rem .75rem;
      border-radius:.8rem;
      border:1px solid var(--border);
      background:#000;
      color:var(--accent);
      outline:none;
      max-width: 100%;
      font-size: 16px; /* ✅ prevent iOS zoom */
    }
    input:focus{border-color: var(--border2)}

    .muted{color:var(--muted);font-size:.82rem}
    .right{white-space:nowrap}

    code{
      font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace;
      font-size:.82rem;color:var(--muted)
    }

    /* ------------------------------------------------------------
       ✅ NAV: desktop/tablet = scroll row + fade edges
       ✅ mobile = grid menu (no horizontal scroll)
    ------------------------------------------------------------ */
    .nav-wrap{
      position: relative;
      width: 100%;
      max-width: 100%;
      margin-top:.75rem;
    }

    /* Default: horizontal scroller */
    .nav{
      display:flex;
      gap:.5rem;
      width:100%;
      max-width:100%;
      overflow-x:auto;
      overflow-y:hidden;
      padding: .15rem .1rem .35rem .1rem;
      -webkit-overflow-scrolling: touch;
      scrollbar-width: none;
      scroll-snap-type: x proximity;
      overscroll-behavior-x: contain;
      white-space: nowrap;
    }
    .nav::-webkit-scrollbar{display:none}

    .nav-wrap::before,
    .nav-wrap::after{
      content:"";
      position:absolute;
      top:0;
      bottom:0;
      width: 26px;
      pointer-events:none;
      z-index: 2;
    }
    .nav-wrap::before{
      left:0;
      background: linear-gradient(90deg, rgba(10,10,10,1) 0%, rgba(10,10,10,0) 100%);
    }
    .nav-wrap::after{
      right:0;
      background: linear-gradient(270deg, rgba(10,10,10,1) 0%, rgba(10,10,10,0) 100%);
    }

    .nav a{
      flex:0 0 auto;
      display:inline-flex;
      gap:.5rem;
      align-items:center;
      border:1px solid var(--border);
      border-radius:999px;
      padding:.42rem .75rem;
      font-size:.8rem;
      color:var(--accent);
      background: rgba(0,0,0,.65);
      transition: background .12s ease, border-color .12s ease;
      white-space:nowrap;
      scroll-snap-align: start;
      text-decoration:none;
      max-width: 100%;
    }
    .nav a:hover{background: rgba(0,0,0,.85); border-color: var(--border2)}
    .nav a.active{border-color:#3a3a3a;background: rgba(0,0,0,.92)}

    .ico{width:16px;height:16px;display:inline-block}
    .ico svg{
      width:16px;height:16px;display:block;
      fill:none;stroke:var(--muted);stroke-width:1.8;
      stroke-linecap:round;stroke-linejoin:round
    }
    .nav a.active .ico svg{stroke: var(--accent)}

    /* ✅ Mobile: grid menu (fixes the “weird” look) */
    @media (max-width: 520px){
      .nav-wrap::before,
      .nav-wrap::after{ display:none; }

      .nav{
        overflow: visible;
        white-space: normal;
        scroll-snap-type: none;
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: .5rem;
        padding: .25rem 0 0 0;
      }

      .nav a{
        width: 100%;
        justify-content: center;
        padding: .58rem .7rem;
        border-radius: .9rem;
        background: rgba(0,0,0,.78);
      }
    }
    @media (min-width: 420px) and (max-width: 520px){
      .nav{ grid-template-columns: repeat(3, minmax(0, 1fr)); }
    }

    /* Table container */
    .table{
      width:100%;
      margin-top:1rem;
      border:1px solid var(--border);
      border-radius:1rem;
      overflow:hidden;
      background: rgba(0,0,0,.35);
    }
    .table-scroll{
      width:100%;
      overflow:auto;
      -webkit-overflow-scrolling: touch;
    }
    table{width:100%;border-collapse:collapse;min-width:720px}
    thead th{
      text-align:left;
      font-size:.72rem;
      letter-spacing:.18em;
      text-transform:uppercase;
      color:var(--muted);
      background: rgba(0,0,0,.80);
      border-bottom:1px solid var(--border);
      padding:.8rem .9rem;
      position: sticky;
      top: 0;
      z-index: 1;
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
    }
    tbody td{
      padding:.78rem .9rem;
      border-bottom:1px solid rgba(255,255,255,.08);
      vertical-align:top;
      font-size:.92rem;
      color:var(--accent);
    }
    tbody tr:hover{background:rgba(255,255,255,.03)}

    pre{
      margin-top:1rem;
      background:#000;
      border:1px solid var(--border);
      border-radius:1rem;
      padding:1rem;
      overflow:auto;
      max-height:320px;
    }

    .footer{
      margin-top:1rem;
      padding-top:.85rem;
      border-top:1px solid var(--border);
      display:flex;
      justify-content:space-between;
      gap:.75rem;
      font-size:.75rem;
      color:var(--muted);
      flex-wrap:wrap;
    }

    a{color:var(--accent);text-decoration:none}
    a:hover{text-decoration:underline}

    /* Breakpoints */
    @media (max-width: 860px){
      .top{grid-template-columns: 1fr}
      .panel{padding:.85rem}
      table{min-width: 640px}
    }
    @media (max-width: 520px){
      .panel{padding:.78rem;border-radius:.95rem}
      .actions{align-items: stretch}
      button{width:100%}
    }
    @media (max-width: 420px){
      body{
        padding:
          calc(.9rem + env(safe-area-inset-top))
          calc(.75rem + env(safe-area-inset-right))
          calc(.9rem + env(safe-area-inset-bottom))
          calc(.75rem + env(safe-area-inset-left));
      }
      .card{padding:.95rem}
      h1{font-size:1.18rem}
      p{font-size:.9rem}
      .actions{gap:.55rem}
      .actions > *{flex:1 1 auto}
    }
  </style>
</head>
<body>
  <div class="wrap">
    <a class="home-fab" href="/api/v1/ui" data-scroll-top="1" aria-label="Back to home">
      <span class="ico">{{ICON_HOME}}</span>
      Home
    </a>
    <div class="card">{{BODY}}</div>
  </div>

  <script>
    const $ = (sel) => document.querySelector(sel)

    // Always scroll to top when navigating via internal links.
    try { history.scrollRestoration = 'manual' } catch(e) {}

    function forceTopOnce(){
      try{ sessionStorage.setItem('__forgeon_force_top__','1') }catch(e){}
    }
    function consumeForceTop(){
      try{
        const v = sessionStorage.getItem('__forgeon_force_top__')
        if(v === '1'){
          sessionStorage.removeItem('__forgeon_force_top__')
          window.scrollTo({ top: 0, left: 0, behavior: 'instant' })
        }
      }catch(e){}
    }

    window.addEventListener('pageshow', () => { consumeForceTop() })

    document.addEventListener('click', (ev) => {
      const a = ev.target && ev.target.closest ? ev.target.closest('a') : null
      if(!a) return
      if(a.hasAttribute('target')) return
      const href = a.getAttribute('href') || ''
      if(!href) return
      if(href.startsWith('#')) return
      if(href.startsWith('http://') || href.startsWith('https://')) return

      forceTopOnce()
      window.scrollTo({ top: 0, left: 0, behavior: 'instant' })
    })

    function fmtDate(s){
      if(!s) return ''
      const d = new Date(s)
      if(Number.isNaN(d.getTime())) return s
      return d.toLocaleString()
    }

    function htmlEscape(s){
      return String(s)
        .replaceAll('&','&amp;')
        .replaceAll('<','&lt;')
        .replaceAll('>','&gt;')
        .replaceAll('"','&quot;')
        .replaceAll("'",'&#39;')
    }

    async function fetchJSON(url){
      const res = await fetch(url, { headers: { 'Accept': 'application/json' } })
      if(!res.ok) throw new Error('HTTP ' + res.status)
      return await res.json()
    }

    {{SCRIPT}}
  </script>
</body>
</html>`

func htmlShell(title, body, script string) string {
	s := htmlShellTemplate
	s = strings.ReplaceAll(s, "{{TITLE}}", html.EscapeString(title))
	s = strings.ReplaceAll(s, "{{BODY}}", body)
	s = strings.ReplaceAll(s, "{{SCRIPT}}", script)
	s = strings.ReplaceAll(s, "{{ICON_HOME}}", iconHome())
	return s
}

// ------------------------------------------------------------
// Icons / badges
// ------------------------------------------------------------

func deviconGo() string {
	return `<svg viewBox="0 0 128 128" aria-hidden="true" role="img">
  <path fill="currentColor" d="M64 14c-28.7 0-52 17.8-52 39.7v20.6C12 96.2 35.3 114 64 114s52-17.8 52-39.7V53.7C116 31.8 92.7 14 64 14Zm0 12c21.7 0 40 12.8 40 27.7v20.6C104 89.2 85.7 102 64 102s-40-12.8-40-27.7V53.7C24 38.8 42.3 26 64 26Z"/>
  <path fill="currentColor" d="M52 56c-7.7 0-14 5.4-14 12s6.3 12 14 12 14-5.4 14-12-6.3-12-14-12Zm0 8c2.8 0 5 1.8 5 4s-2.2 4-5 4-5-1.8-5-4 2.2-4 5-4Z"/>
  <path fill="currentColor" d="M82 56c-7.7 0-14 5.4-14 12s6.3 12 14 12 14-5.4 14-12-6.3-12-14-12Zm0 8c2.8 0 5 1.8 5 4s-2.2 4-5 4-5-1.8-5-4 2.2-4 5-4Z"/>
  <path fill="currentColor" d="M73 88H55a6 6 0 0 1 0-12h18a6 6 0 0 1 0 12Z"/>
</svg>`
}

func deviconSupabase() string {
	return `<svg viewBox="0 0 128 128" aria-hidden="true" role="img">
  <path fill="currentColor" d="M77.6 7.7c-2.5-3-7.2-3-9.7 0L21.4 61.4c-1.6 1.9-2 4.6-1 6.9 1 2.3 3.3 3.8 5.8 3.8h28.5l-4.7 48.6c-.3 3 1.6 5.8 4.6 6.6 3 .8 6.2-.5 7.7-3.2l44.5-75.2c1.2-2 1.2-4.5.1-6.5-1.1-2.1-3.3-3.3-5.7-3.3H73.5L77.6 7.7Zm-5.3 43.4h24.2L62.8 108.6 66.7 64H32.2L72.3 18.3l-4.9 32.8Z"/>
</svg>`
}

func uiStackBadges() string {
	return fmt.Sprintf(`
<div class="stack">
  <span class="stack-badge" title="Golang">
    <span style="color:var(--accent)">%s</span>
    <strong>Go</strong>
    <span>Fiber</span>
  </span>

  <span class="stack-badge" title="Supabase Postgres">
    <span style="color:var(--accent)">%s</span>
    <strong>Supabase</strong>
    <span>Postgres</span>
  </span>
</div>`, deviconGo(), deviconSupabase())
}

func uiTop(title, desc, nav, right string) string {
	// right can be empty -> it collapses fine on mobile
	return fmt.Sprintf(`
<div class="top">
  <div>
    <div class="eyebrow">Forgeon · UI</div>
    <h1 class="rainbow">%s</h1>
    <p>%s</p>
    %s
    %s
  </div>
  %s
</div>`,
		html.EscapeString(title),
		html.EscapeString(desc),
		uiStackBadges(),
		nav,
		right,
	)
}

func uiNav(activePath string) string {
	type item struct {
		href  string
		label string
		icon  string
	}
	items := []item{
		{href: "/api/v1/ui", label: "Home", icon: iconHome()},
		{href: "/api/v1/ui/users", label: "Users", icon: iconUser()},
		{href: "/api/v1/ui/projects", label: "Projects", icon: iconBox()},
		{href: "/api/v1/ui/tasks", label: "Tasks", icon: iconCheckCircle()},
		{href: "/api/v1/ui/books", label: "Books", icon: iconBook()},
	}

	var b strings.Builder
	b.WriteString(`<div class="nav-wrap"><div class="nav">`)
	for _, it := range items {
		cl := ""
		if it.href == activePath {
			cl = ` class="active"`
		}
		b.WriteString(fmt.Sprintf(`<a%s href="%s" data-scroll-top="1"><span class="ico">%s</span>%s</a>`,
			cl,
			html.EscapeString(it.href),
			it.icon,
			html.EscapeString(it.label),
		))
	}
	b.WriteString(`</div></div>`)
	return b.String()
}

// ------------------------------------------------------------
// Inline SVG icons
// ------------------------------------------------------------

func iconHome() string {
	return `<svg viewBox="0 0 24 24"><path d="M3 10.5 12 3l9 7.5"/><path d="M5 10v10h14V10"/></svg>`
}
func iconUser() string {
	return `<svg viewBox="0 0 24 24"><path d="M20 21a8 8 0 0 0-16 0"/><path d="M12 13a4 4 0 1 0-4-4 4 4 0 0 0 4 4Z"/></svg>`
}
func iconBox() string {
	return `<svg viewBox="0 0 24 24"><path d="M21 8.5 12 3 3 8.5 12 14 21 8.5Z"/><path d="M3 8.5V20l9 4 9-4V8.5"/><path d="M12 14v10"/></svg>`
}
func iconCheckCircle() string {
	return `<svg viewBox="0 0 24 24"><path d="M22 12a10 10 0 1 1-10-10 10 10 0 0 1 10 10Z"/><path d="m7.5 12.2 3 3 6-6"/></svg>`
}
func iconBook() string {
	return `<svg viewBox="0 0 24 24"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 3H20v18H6.5A2.5 2.5 0 0 1 4 18.5V5.5A2.5 2.5 0 0 1 6.5 3Z"/></svg>`
}

// ------------------------------------------------------------
// Pages
// ------------------------------------------------------------

func UIHomePage(c *fiber.Ctx) error {
	body := uiTop(
		"UI demo pages for FE - (Golang + Supabase)",
		"HTML + fetch pages backed by your real JSON endpoints.",
		uiNav("/api/v1/ui"),
		"",
	) + `
<div class="table">
  <div class="table-scroll">
	<table>
		<thead><tr><th>Page</th><th>Purpose</th><th>API used</th></tr></thead>
		<tbody>
		<tr><td><a href="/api/v1/ui/users" data-scroll-top="1">/api/v1/ui/users</a></td><td class="muted">List users</td><td class="muted"><code>/api/v1/users</code></td></tr>
		<tr><td><a href="/api/v1/ui/projects" data-scroll-top="1">/api/v1/ui/projects</a></td><td class="muted">List projects</td><td class="muted"><code>/api/v1/projects</code></td></tr>
		<tr><td><a href="/api/v1/ui/tasks" data-scroll-top="1">/api/v1/ui/tasks</a></td><td class="muted">List tasks</td><td class="muted"><code>/api/v1/tasks</code></td></tr>
		<tr><td><a href="/api/v1/ui/books" data-scroll-top="1">/api/v1/ui/books</a></td><td class="muted">List books</td><td class="muted"><code>/api/v1/books</code></td></tr>
		</tbody>
	</table>
  </div>
</div>

<div class="footer">
  <span>Tip: FE can inspect JSON quickly by opening the API links.</span>
  <span class="pill">ui-playground</span>
</div>
`
	return c.Type("html", "utf-8").SendString(htmlShell("Forgeon · UI", body, ""))
}

func UITasksPage(c *fiber.Ctx) error {
	now := time.Now().Format(time.RFC3339)

	right := fmt.Sprintf(`
  <div class="panel">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/tasks?page=1&page_size=10" />
    <div class="pill-row">
      <div class="pill">count: <strong id="count">0</strong></div>
      <div class="pill">page: <strong id="page">1</strong></div>
      <div class="pill">page_size: <strong id="page_size">10</strong></div>
    </div>
    <div class="muted" style="margin-top:.35rem">Rendered: <code>%s</code></div>
  </div>`, html.EscapeString(now))

	body := uiTop(
		"Tasks dashboard",
		"Calls /api/v1/tasks and renders real data.",
		uiNav("/api/v1/ui/tasks"),
		right,
	) + `
<div class="table">
  <div class="table-scroll">
  <table>
    <thead>
      <tr>
        <th>Title</th>
        <th>Project ID</th>
        <th>Status</th>
        <th>Due</th>
        <th>Created</th>
      </tr>
    </thead>
    <tbody id="tbody">
      <tr><td colspan="5" class="muted">Loading…</td></tr>
    </tbody>
  </table>
</div>
</div>

<div class="footer">
  <span>JSON: <a href="/api/v1/tasks" target="_blank">/api/v1/tasks</a></span>
  <span class="pill">tasks-ui</span>
</div>

<pre><code id="raw"></code></pre>
`

	script := `
function render(payload){
  $('#count').textContent = payload.count ?? 0
  $('#page').textContent = payload.page ?? 1
  $('#page_size').textContent = payload.page_size ?? 10

  const rows = (payload.tasks || []).map((t) => {
    return ''
      + '<tr>'
      + ' <td>' + htmlEscape(t.title || '') + '<div class="muted"><code>' + htmlEscape(t.id || '') + '</code></div></td>'
      + ' <td class="muted"><code>' + htmlEscape(t.project_id || '') + '</code></td>'
      + ' <td class="muted"><code>' + htmlEscape(t.status || '') + '</code></td>'
      + ' <td class="right muted">' + htmlEscape(fmtDate(t.due_at)) + '</td>'
      + ' <td class="right muted">' + htmlEscape(fmtDate(t.created_at)) + '</td>'
      + '</tr>'
  }).join('')

  $('#tbody').innerHTML = rows || '<tr><td colspan="5" class="muted">No tasks.</td></tr>'
  $('#raw').textContent = JSON.stringify(payload, null, 2)
}

async function load(){
  const url = $('#api').value.trim()
  $('#status').textContent = 'Loading…'
  try{
    const data = await fetchJSON(url)
    render(data)
    $('#status').textContent = 'OK'
  }catch(e){
    $('#status').textContent = 'Error'
    $('#tbody').innerHTML = '<tr><td colspan="5" class="muted">' + htmlEscape(e.message || String(e)) + '</td></tr>'
    $('#raw').textContent = ''
  }
}
$('#reload').addEventListener('click', load)
load()
`
	return c.Status(200).Type("html", "utf-8").SendString(htmlShell("Forgeon · Tasks UI", body, script))
}

func UIUsersPage(c *fiber.Ctx) error {
	right := `
  <div class="panel">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/users?page=1&page_size=10" />
  </div>`

	body := uiTop(
		"Users list",
		"Calls /api/v1/users and renders real data.",
		uiNav("/api/v1/ui/users"),
		right,
	) + `
<div class="table">
  <div class="table-scroll">
  <table>
    <thead>
      <tr>
        <th>User</th>
        <th>Email</th>
        <th>Flags</th>
        <th>Created</th>
      </tr>
    </thead>
    <tbody id="tbody">
      <tr><td colspan="4" class="muted">Loading…</td></tr>
    </tbody>
  </table>
</div>
</div>

<div class="footer">
  <span>JSON: <a href="/api/v1/users" target="_blank">/api/v1/users</a></span>
  <span class="pill">users-ui</span>
</div>

<pre><code id="raw"></code></pre>
`

	script := `
function render(payload){
  const arr = payload.users || payload.items || []
  const rows = arr.map((u) => {
    const flags = []
    if (u.is_admin) flags.push('admin')
    if (u.is_active) flags.push('active')
    if (u.is_deleted) flags.push('deleted')
    return ''
      + '<tr>'
      + ' <td>' + htmlEscape(u.username || '') + '<div class="muted"><code>' + htmlEscape(u.id || '') + '</code></div></td>'
      + ' <td class="muted">' + htmlEscape(u.email || '') + '</td>'
      + ' <td class="muted"><code>' + htmlEscape(flags.join(',')) + '</code></td>'
      + ' <td class="right muted">' + htmlEscape(fmtDate(u.created_at)) + '</td>'
      + '</tr>'
  }).join('')
  $('#tbody').innerHTML = rows || '<tr><td colspan="4" class="muted">No users.</td></tr>'
  $('#raw').textContent = JSON.stringify(payload, null, 2)
}

async function load(){
  const url = $('#api').value.trim()
  $('#status').textContent = 'Loading…'
  try{
    const data = await fetchJSON(url)
    render(data)
    $('#status').textContent = 'OK'
  }catch(e){
    $('#status').textContent = 'Error'
    $('#tbody').innerHTML = '<tr><td colspan="4" class="muted">' + htmlEscape(e.message || String(e)) + '</td></tr>'
    $('#raw').textContent = ''
  }
}
$('#reload').addEventListener('click', load)
load()
`
	return c.Type("html", "utf-8").SendString(htmlShell("Forgeon · Users UI", body, script))
}

func UIProjectsPage(c *fiber.Ctx) error {
	right := `
  <div class="panel">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/projects?page=1&page_size=10" />
  </div>`

	body := uiTop(
		"Projects list",
		"Calls /api/v1/projects and renders real data.",
		uiNav("/api/v1/ui/projects"),
		right,
	) + `
<div class="table">
  <div class="table-scroll">
  <table>
    <thead>
      <tr>
        <th>Project</th>
        <th>Owner</th>
        <th>Description</th>
        <th>Created</th>
      </tr>
    </thead>
    <tbody id="tbody">
      <tr><td colspan="4" class="muted">Loading…</td></tr>
    </tbody>
  </table>
  </div>
</div>

<div class="footer">
  <span>JSON: <a href="/api/v1/projects" target="_blank">/api/v1/projects</a></span>
  <span class="pill">projects-ui</span>
</div>

<pre><code id="raw"></code></pre>
`

	script := `
function render(payload){
  const arr = payload.projects || payload.items || []
  const rows = arr.map((p) => {
    return ''
      + '<tr>'
      + ' <td>' + htmlEscape(p.name || '') + '<div class="muted"><code>' + htmlEscape(p.id || '') + '</code></div></td>'
      + ' <td class="muted"><code>' + htmlEscape(p.owner_user_id || '') + '</code></td>'
      + ' <td class="muted">' + htmlEscape(p.description || '') + '</td>'
      + ' <td class="right muted">' + htmlEscape(fmtDate(p.created_at)) + '</td>'
      + '</tr>'
  }).join('')
  $('#tbody').innerHTML = rows || '<tr><td colspan="4" class="muted">No projects.</td></tr>'
  $('#raw').textContent = JSON.stringify(payload, null, 2)
}

async function load(){
  const url = $('#api').value.trim()
  $('#status').textContent = 'Loading…'
  try{
    const data = await fetchJSON(url)
    render(data)
    $('#status').textContent = 'OK'
  }catch(e){
    $('#status').textContent = 'Error'
    $('#tbody').innerHTML = '<tr><td colspan="4" class="muted">' + htmlEscape(e.message || String(e)) + '</td></tr>'
    $('#raw').textContent = ''
  }
}
$('#reload').addEventListener('click', load)
load()
`
	return c.Type("html", "utf-8").SendString(htmlShell("Forgeon · Projects UI", body, script))
}

func UIBooksPage(c *fiber.Ctx) error {
	right := `
  <div class="panel">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/books?page=1&page_size=10" />
  </div>`

	body := uiTop(
		"Books list",
		"Calls /api/v1/books and renders real data.",
		uiNav("/api/v1/ui/books"),
		right,
	) + `
<div class="table">
  <div class="table-scroll">
  <table>
    <thead>
      <tr>
        <th>Book</th>
        <th>User ID</th>
        <th>Author</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody id="tbody">
      <tr><td colspan="4" class="muted">Loading…</td></tr>
    </tbody>
  </table>
</div>
</div>

<div class="footer">
  <span>JSON: <a href="/api/v1/books" target="_blank">/api/v1/books</a></span>
  <span class="pill">books-ui</span>
</div>

<pre><code id="raw"></code></pre>
`

	script := `
function render(payload){
  const arr = payload.books || payload.items || []
  const rows = arr.map((b) => {
    return ''
      + '<tr>'
      + ' <td>' + htmlEscape(b.title || '') + '<div class="muted"><code>' + htmlEscape(b.id || '') + '</code></div></td>'
      + ' <td class="muted"><code>' + htmlEscape(b.user_id || '') + '</code></td>'
      + ' <td class="muted">' + htmlEscape(b.author || '') + '</td>'
      + ' <td class="muted"><code>' + htmlEscape(String(b.status ?? "")) + '</code></td>'
      + '</tr>'
  }).join('')
  $('#tbody').innerHTML = rows || '<tr><td colspan="4" class="muted">No books.</td></tr>'
  $('#raw').textContent = JSON.stringify(payload, null, 2)
}

async function load(){
  const url = $('#api').value.trim()
  $('#status').textContent = 'Loading…'
  try{
    const data = await fetchJSON(url)
    render(data)
    $('#status').textContent = 'OK'
  }catch(e){
    $('#status').textContent = 'Error'
    $('#tbody').innerHTML = '<tr><td colspan="4" class="muted">' + htmlEscape(e.message || String(e)) + '</td></tr>'
    $('#raw').textContent = ''
  }
}
$('#reload').addEventListener('click', load)
load()
`
	return c.Type("html", "utf-8").SendString(htmlShell("Forgeon · Books UI", body, script))
}
