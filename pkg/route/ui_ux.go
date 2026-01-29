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

func htmlShell(title, body, script string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>%s</title>
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <style>
    :root{
      color-scheme: dark;
      --bg:#050505; --card:#0f0f10; --border:#222;
      --text:#f5f5f5; --muted:#9ca3af; --accent:#e5e5e5;
      --good:#22c55e; --warn:#f59e0b; --dim:#6b7280;
      --focus:#2a2a2a;
    }
    *{box-sizing:border-box;margin:0;padding:0}
    body{
      min-height:100vh;
      font-family:system-ui,-apple-system,BlinkMacSystemFont,"SF Pro Text",sans-serif;
      background:radial-gradient(circle at top,#111 0,#050505 55%%);
      color:var(--text);
      display:flex; align-items:center; justify-content:center;
      padding:2rem 1.5rem;
    }
    .wrap{width:100%%;max-width:1100px}
    .card{
      border-radius:1.25rem;
      border:1px solid var(--border);
      background:radial-gradient(circle at top left,#151515 0,var(--card) 50%%,#050505 100%%);
      padding:1.5rem 1.5rem 1.25rem;
    }
    .top{display:flex;align-items:flex-start;justify-content:space-between;gap:1rem;margin-bottom:1rem}
    .eyebrow{font-size:.7rem;letter-spacing:.22em;text-transform:uppercase;color:var(--muted);margin-bottom:.5rem}
    h1{font-size:1.35rem;line-height:1.2;margin-bottom:.4rem}
    p{font-size:.9rem;line-height:1.6;color:var(--muted)}
    .pill-row{display:flex;flex-wrap:wrap;gap:.45rem;margin-top:.85rem}
    .pill{
      font-size:.7rem;text-transform:uppercase;letter-spacing:.16em;
      padding:.25rem .6rem;border-radius:999px;border:1px solid var(--border);
      color:var(--muted);display:inline-flex;gap:.4rem;align-items:center
    }
    .pill strong{color:var(--accent);font-weight:600}
    .actions{display:flex;gap:.5rem;align-items:center}
    button{
      cursor:pointer;border:1px solid var(--border);background:#0b0b0c;color:var(--accent);
      border-radius:.75rem;padding:.5rem .75rem;font-size:.85rem;
    }
    button:hover{background:#111}
    .hint{font-size:.75rem;color:var(--dim);margin-top:.35rem;text-align:right}
    input{
      width:100%%;margin-top:.4rem;padding:.55rem .7rem;border-radius:.75rem;
      border:1px solid var(--border);background:#0b0b0c;color:var(--accent);outline:none;
    }
    .table{
      width:100%%;margin-top:1rem;border:1px solid var(--border);
      border-radius:1rem;overflow:hidden;
    }
    table{width:100%%;border-collapse:collapse}
    thead th{
      text-align:left;font-size:.72rem;letter-spacing:.18em;text-transform:uppercase;color:var(--muted);
      background:#0b0b0c;border-bottom:1px solid var(--border);padding:.75rem .9rem;
    }
    tbody td{
      padding:.75rem .9rem;border-bottom:1px solid rgba(34,34,34,.65);
      vertical-align:top;font-size:.9rem;color:var(--accent);
    }
    tbody tr:hover{background:rgba(255,255,255,.02)}
    .muted{color:var(--muted);font-size:.82rem}
    .right{white-space:nowrap}
    code{font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace;font-size:.8rem;color:var(--muted)}
    .footer{
      margin-top:1rem;padding-top:.85rem;border-top:1px solid var(--border);
      display:flex;justify-content:space-between;gap:.75rem;font-size:.75rem;color:var(--muted);
    }
    a{color:var(--accent);text-decoration:none}
    a:hover{text-decoration:underline}

    /* Nav */
    .nav{display:flex;flex-wrap:wrap;gap:.5rem;margin-top:.75rem}
    .nav a{
      display:inline-flex;gap:.5rem;align-items:center;
      border:1px solid var(--border);border-radius:999px;
      padding:.35rem .65rem;font-size:.78rem;color:var(--accent);
      background:#0b0b0c;
      transition: background .12s ease, border-color .12s ease;
    }
    .nav a:hover{background:#111}
    .nav a.active{
      border-color: #3a3a3a;
      background: #141414;
    }
    .ico{
      width:16px;height:16px;display:inline-block;
    }
    .ico svg{width:16px;height:16px;display:block;fill:none;stroke:var(--muted);stroke-width:1.8;stroke-linecap:round;stroke-linejoin:round}
    .nav a.active .ico svg{stroke:var(--accent)}
    pre{margin-top:1rem;background:#050505;border:1px solid var(--border);border-radius:1rem;padding:1rem;overflow:auto}

	    /* Stack badges */
    .stack{
      display:flex;align-items:center;gap:.55rem;
      margin-top:.65rem;flex-wrap:wrap;
    }
    .stack-badge{
      display:inline-flex;align-items:center;gap:.45rem;
      border:1px solid var(--border);
      background:#0b0b0c;
      border-radius:999px;
      padding:.28rem .6rem;
      font-size:.78rem;
      color:var(--muted);
    }
    .stack-badge svg{width:16px;height:16px;display:block}
    .stack-badge strong{color:var(--accent);font-weight:600}

  </style>
</head>
<body>
  <div class="wrap">
    <div class="card">%s</div>
  </div>

  <script>
    const $ = (sel) => document.querySelector(sel)

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

    %s
  </script>
</body>
</html>`, html.EscapeString(title), body, script)
}

func deviconGo() string {
	// simple Go mark (monochrome-friendly)
	return `<svg viewBox="0 0 128 128" aria-hidden="true" role="img">
  <path fill="currentColor" d="M64 14c-28.7 0-52 17.8-52 39.7v20.6C12 96.2 35.3 114 64 114s52-17.8 52-39.7V53.7C116 31.8 92.7 14 64 14Zm0 12c21.7 0 40 12.8 40 27.7v20.6C104 89.2 85.7 102 64 102s-40-12.8-40-27.7V53.7C24 38.8 42.3 26 64 26Z"/>
  <path fill="currentColor" d="M52 56c-7.7 0-14 5.4-14 12s6.3 12 14 12 14-5.4 14-12-6.3-12-14-12Zm0 8c2.8 0 5 1.8 5 4s-2.2 4-5 4-5-1.8-5-4 2.2-4 5-4Z"/>
  <path fill="currentColor" d="M82 56c-7.7 0-14 5.4-14 12s6.3 12 14 12 14-5.4 14-12-6.3-12-14-12Zm0 8c2.8 0 5 1.8 5 4s-2.2 4-5 4-5-1.8-5-4 2.2-4 5-4Z"/>
  <path fill="currentColor" d="M73 88H55a6 6 0 0 1 0-12h18a6 6 0 0 1 0 12Z"/>
</svg>`
}

func deviconSupabase() string {
	// Supabase "S" mark (monochrome-friendly)
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

func uiTop(title, desc, nav string) string {
	return fmt.Sprintf(`
<div class="top">
  <div>
    <div class="eyebrow">Forgeon · UI</div>
    <h1>%s</h1>
    <p>%s</p>
    %s
    %s
  </div>
</div>`,
		html.EscapeString(title),
		html.EscapeString(desc),
		uiStackBadges(), // ✅ added here
		nav,
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
	b.WriteString(`<div class="nav">`)
	for _, it := range items {
		cl := ""
		if it.href == activePath {
			cl = ` class="active"`
		}
		b.WriteString(fmt.Sprintf(`<a%s href="%s"><span class="ico">%s</span>%s</a>`,
			cl,
			html.EscapeString(it.href),
			it.icon,
			html.EscapeString(it.label),
		))
	}
	b.WriteString(`</div>`)
	return b.String()
}

// ------------------------------------------------------------
// Inline SVG icons (no emoji, no external deps)
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
	) + `
<div class="table">
  <table>
    <thead><tr><th>Page</th><th>Purpose</th><th>API used</th></tr></thead>
    <tbody>
      <tr><td><a href="/api/v1/ui/users">/api/v1/ui/users</a></td><td class="muted">List users</td><td class="muted"><code>/api/v1/users</code></td></tr>
      <tr><td><a href="/api/v1/ui/projects">/api/v1/ui/projects</a></td><td class="muted">List projects</td><td class="muted"><code>/api/v1/projects</code></td></tr>
      <tr><td><a href="/api/v1/ui/tasks">/api/v1/ui/tasks</a></td><td class="muted">List tasks</td><td class="muted"><code>/api/v1/tasks</code></td></tr>
      <tr><td><a href="/api/v1/ui/books">/api/v1/ui/books</a></td><td class="muted">List books</td><td class="muted"><code>/api/v1/books</code></td></tr>
    </tbody>
  </table>
</div>

<div class="footer">
  <span>Tip: FE can inspect JSON quickly by opening the API links.</span>
  <span class="pill">ui-playground</span>
</div>
`
	return c.Type("html", "utf-8").SendString(htmlShell("Forgeon · UI", body, ""))
}

// -----------------------------
// Tasks UI
// -----------------------------

func UITasksPage(c *fiber.Ctx) error {
	now := time.Now().Format(time.RFC3339)

	body := fmt.Sprintf(`
<div class="top">
  <div>
    <div class="eyebrow">Forgeon · Tasks UI</div>
    <h1>Tasks dashboard</h1>
    <p>Calls <code>/api/v1/tasks</code> and renders real data.</p>
    %s
  </div>

  <div style="min-width:320px">
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
  </div>
</div>

<div class="table">
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

<div class="footer">
  <span>JSON: <a href="/api/v1/tasks" target="_blank">/api/v1/tasks</a></span>
  <span class="pill">tasks-ui</span>
</div>

<pre><code id="raw"></code></pre>
`, uiNav("/api/v1/ui/tasks"), html.EscapeString(now))

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

// -----------------------------
// Users UI
// -----------------------------

func UIUsersPage(c *fiber.Ctx) error {
	body := fmt.Sprintf(`
<div class="top">
  <div>
    <div class="eyebrow">Forgeon · Users UI</div>
    <h1>Users list</h1>
    <p>Calls <code>/api/v1/users</code> and renders real data.</p>
    %s
  </div>

  <div style="min-width:320px">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/users?page=1&page_size=10" />
  </div>
</div>

<div class="table">
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

<div class="footer">
  <span>JSON: <a href="/api/v1/users" target="_blank">/api/v1/users</a></span>
  <span class="pill">users-ui</span>
</div>

<pre><code id="raw"></code></pre>
`, uiNav("/api/v1/ui/users"))

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

// -----------------------------
// Projects UI
// -----------------------------

func UIProjectsPage(c *fiber.Ctx) error {
	body := fmt.Sprintf(`
<div class="top">
  <div>
    <div class="eyebrow">Forgeon · Projects UI</div>
    <h1>Projects list</h1>
    <p>Calls <code>/api/v1/projects</code> and renders real data.</p>
    %s
  </div>

  <div style="min-width:320px">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/projects?page=1&page_size=10" />
  </div>
</div>

<div class="table">
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

<div class="footer">
  <span>JSON: <a href="/api/v1/projects" target="_blank">/api/v1/projects</a></span>
  <span class="pill">projects-ui</span>
</div>

<pre><code id="raw"></code></pre>
`, uiNav("/api/v1/ui/projects"))

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

// -----------------------------
// Books UI
// -----------------------------

func UIBooksPage(c *fiber.Ctx) error {
	body := fmt.Sprintf(`
<div class="top">
  <div>
    <div class="eyebrow">Forgeon · Books UI</div>
    <h1>Books list</h1>
    <p>Calls <code>/api/v1/books</code> and renders real data.</p>
    %s
  </div>

  <div style="min-width:320px">
    <div class="actions">
      <button id="reload">Reload</button>
      <span class="muted">Status: <strong id="status">…</strong></span>
    </div>
    <div class="hint">API URL</div>
    <input id="api" value="/api/v1/books?page=1&page_size=10" />
  </div>
</div>

<div class="table">
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

<div class="footer">
  <span>JSON: <a href="/api/v1/books" target="_blank">/api/v1/books</a></span>
  <span class="pill">books-ui</span>
</div>

<pre><code id="raw"></code></pre>
`, uiNav("/api/v1/ui/books"))

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
