package server

import "net/http"

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

const indexHTML = `<!doctype html>
<html lang="pt-BR">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>llm-memory</title>
<style>
:root{color-scheme:dark;--bg:#0b0f14;--panel:#111821;--muted:#8ba0b3;--txt:#e6edf3;--line:#243244;--accent:#7ee787;--bad:#ff7b72}*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--txt);font:14px/1.45 system-ui,Segoe UI,Roboto,Arial}header{padding:18px 22px;border-bottom:1px solid var(--line);display:flex;gap:16px;align-items:center;justify-content:space-between}h1{font-size:18px;margin:0}.wrap{display:grid;grid-template-columns:360px 1fr;gap:16px;padding:16px}.card{background:var(--panel);border:1px solid var(--line);border-radius:12px;padding:14px}label{display:block;margin:8px 0 4px;color:var(--muted)}input,textarea,select{width:100%;background:#0d141c;color:var(--txt);border:1px solid var(--line);border-radius:8px;padding:8px}textarea{min-height:110px}button{background:#1f6feb;color:white;border:0;border-radius:8px;padding:9px 12px;cursor:pointer}button.danger{background:#9b2428}.row{display:flex;gap:8px}.row>*{flex:1}.muted{color:var(--muted)}.mem{border:1px solid var(--line);border-radius:10px;padding:12px;margin:10px 0;background:#0d141c}.mem pre{white-space:pre-wrap;margin:8px 0}.tag{display:inline-block;color:#111;background:var(--accent);padding:2px 6px;border-radius:999px;margin-right:4px;font-size:12px}.err{color:var(--bad)}code{color:var(--accent)}@media(max-width:900px){.wrap{grid-template-columns:1fr}}
</style>
</head>
<body>
<header><h1>☣️ llm-memory</h1><div class="muted" id="cfg">carregando config...</div></header>
<div class="wrap">
  <section class="card">
    <h2>Nova / editar memória</h2>
    <input id="id" placeholder="id opcional" />
    <div class="row"><div><label>type</label><select id="type"><option>preference</option><option>fact</option><option>decision</option><option>task</option><option>note</option><option>relationship</option></select></div><div><label>scope</label><select id="scope"><option>global</option><option>project</option><option>session</option><option>private</option></select></div></div>
    <label>subject</label><input id="subject" value="botmaster" />
    <label>content</label><textarea id="content"></textarea>
    <div class="row"><div><label>source.kind</label><input id="sourceKind" value="gui" /></div><div><label>source.ref</label><input id="sourceRef" value="local" /></div></div>
    <div class="row"><div><label>confidence</label><input id="confidence" type="number" min="0" max="1" step="0.01" value="0.90" /></div><div><label>tags, vírgula</label><input id="tags" placeholder="style, preference" /></div></div>
    <p><button onclick="saveMemory()">Salvar</button> <button onclick="clearForm()">Limpar</button></p>
    <p id="formMsg" class="muted"></p>
  </section>
  <main class="card">
    <h2>Buscar</h2>
    <div class="row"><input id="q" placeholder="texto FTS: respostas diretas" onkeydown="if(event.key==='Enter') search()"/><input id="filterSubject" placeholder="subject" /></div>
    <p><button onclick="search()">Buscar</button> <button onclick="loadEvents()">Eventos</button></p>
    <div id="out"></div>
  </main>
</div>
<script>
async function api(path, opt={}){const r=await fetch(path,{headers:{'content-type':'application/json'},...opt});const j=await r.json().catch(()=>null);if(!r.ok)throw new Error((j&&j.error)||r.statusText);return j}
function val(id){return document.getElementById(id).value.trim()}
function memoryFromForm(){return {id:val('id')||undefined,type:val('type'),subject:val('subject'),content:val('content'),source:{kind:val('sourceKind'),ref:val('sourceRef')},scope:val('scope'),confidence:parseFloat(val('confidence')||'0.9'),tags:val('tags')?val('tags').split(',').map(s=>s.trim()).filter(Boolean):[],embedding_refs:{}}}
async function saveMemory(){try{const m=await api('/api/memories',{method:'POST',body:JSON.stringify(memoryFromForm())});document.getElementById('formMsg').textContent='salvo: '+m.id;document.getElementById('id').value=m.id;search()}catch(e){document.getElementById('formMsg').innerHTML='<span class="err">'+e.message+'</span>'}}
function clearForm(){for(const id of ['id','content','tags'])document.getElementById(id).value='';document.getElementById('confidence').value='0.90'}
async function search(){const body={text:val('q'),subject:val('filterSubject'),limit:50};const items=await api('/api/search',{method:'POST',body:JSON.stringify(body)});renderMemories(items)}
function renderMemories(items){const out=document.getElementById('out');out.innerHTML=items.map(m=>'<div class="mem"><b>'+escapeHTML(m.type)+'</b> <code>'+escapeHTML(m.id)+'</code> <span class="muted">'+escapeHTML(m.scope)+' · '+escapeHTML(m.subject)+' · conf '+m.confidence+'</span><pre>'+escapeHTML(m.content)+'</pre><div>'+(m.tags||[]).map(t=>'<span class="tag">'+escapeHTML(t)+'</span>').join('')+'</div><p><button onclick=\'editMemory('+JSON.stringify(m).replaceAll("'","&#39;")+')\'>Editar</button> <button class="danger" onclick="forget(\''+m.id+'\')">Forget</button></p></div>').join('')||'<p class="muted">sem resultados</p>'}
function editMemory(m){document.getElementById('id').value=m.id;document.getElementById('type').value=m.type;document.getElementById('scope').value=m.scope;document.getElementById('subject').value=m.subject;document.getElementById('content').value=m.content;document.getElementById('sourceKind').value=m.source.kind;document.getElementById('sourceRef').value=m.source.ref;document.getElementById('confidence').value=m.confidence;document.getElementById('tags').value=(m.tags||[]).join(', ')}
async function forget(id){if(!confirm('Forget '+id+'?'))return;await api('/api/memories/'+id,{method:'DELETE'});search()}
async function loadEvents(){const items=await api('/api/events?limit=100');document.getElementById('out').innerHTML='<h3>Eventos</h3>'+items.map(e=>'<div class="mem"><b>'+escapeHTML(e.kind)+'</b> <code>'+escapeHTML(e.id)+'</code><pre>'+escapeHTML(e.payload)+'</pre><span class="muted">'+escapeHTML(e.source.kind)+':'+escapeHTML(e.source.ref)+' · '+e.created_at+'</span></div>').join('')}
function escapeHTML(s){return String(s||'').replace(/[&<>"]/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c]))}
api('/api/config').then(c=>document.getElementById('cfg').textContent='db '+c.database.path+' · llm '+c.llm.provider+'/'+(c.llm.model||'-')+' · embedding '+c.embedding.provider+'/'+(c.embedding.model||'-'));search();
</script>
</body>
</html>`
