import { useEffect, useState, useCallback, useRef } from "react"
import { getTraceId } from "../services/tracing"
import { ServiceEndpoints } from "../services/endpoints"

interface JaegerSpan {
  spanID: string
  operationName: string
  startTime: number
  duration: number
  processID: string
  tags: { key: string; type: string; value: unknown }[]
  warnings: string[] | null
  references?: { refType: string; traceID: string; spanID: string }[]
}

// ── Service graph helpers ──────────────────────────────────────────────────

const NODE_W = 130
const NODE_H = 38
const H_GAP = 52
const V_GAP = 14
const PAD = 14

function buildServiceGraph(trace: JaegerTrace) {
  const spanMap = new Map(trace.spans.map(s => [s.spanID, s]))
  const edgeSet = new Set<string>()
  const edges: { from: string; to: string }[] = []
  const nodes = new Set<string>(Object.values(trace.processes).map(p => p.serviceName))

  for (const span of trace.spans) {
    if (!span.references) continue
    const toSvc = trace.processes[span.processID]?.serviceName
    for (const ref of span.references) {
      if (ref.refType !== "CHILD_OF") continue
      const parent = spanMap.get(ref.spanID)
      if (!parent) continue
      const fromSvc = trace.processes[parent.processID]?.serviceName
      if (fromSvc && toSvc && fromSvc !== toSvc) {
        const key = `${fromSvc}|${toSvc}`
        if (!edgeSet.has(key)) {
          edgeSet.add(key)
          edges.push({ from: fromSvc, to: toSvc })
        }
      }
    }
  }
  return { nodes: [...nodes], edges }
}

function layoutGraph(nodes: string[], edges: { from: string; to: string }[]) {
  const outEdges = new Map<string, string[]>()
  const inDegree = new Map<string, number>()
  for (const n of nodes) { outEdges.set(n, []); inDegree.set(n, 0) }
  for (const e of edges) {
    outEdges.get(e.from)?.push(e.to)
    inDegree.set(e.to, (inDegree.get(e.to) ?? 0) + 1)
  }

  const layer = new Map<string, number>()
  const queue = nodes.filter(n => inDegree.get(n) === 0)
  queue.forEach(n => layer.set(n, 0))
  for (let qi = 0; qi < queue.length; qi++) {
    const n = queue[qi]
    for (const child of outEdges.get(n) ?? []) {
      const nl = (layer.get(n) ?? 0) + 1
      if (!layer.has(child) || layer.get(child)! < nl) {
        layer.set(child, nl)
        queue.push(child)
      }
    }
  }
  nodes.forEach(n => { if (!layer.has(n)) layer.set(n, 0) })

  const byLayer: string[][] = []
  for (const [n, l] of layer) {
    while (byLayer.length <= l) byLayer.push([])
    byLayer[l].push(n)
  }

  const pos = new Map<string, { x: number; y: number }>()
  byLayer.forEach((ln, l) =>
    ln.forEach((n, i) =>
      pos.set(n, { x: PAD + l * (NODE_W + H_GAP), y: PAD + i * (NODE_H + V_GAP) })
    )
  )

  const svgW = PAD * 2 + byLayer.length * NODE_W + Math.max(0, byLayer.length - 1) * H_GAP
  const maxRows = Math.max(...byLayer.map(l => l.length))
  const svgH = PAD * 2 + maxRows * NODE_H + Math.max(0, maxRows - 1) * V_GAP

  return { pos, svgW, svgH }
}

function ServiceGraph({ trace }: { trace: JaegerTrace }) {
  const { nodes, edges } = buildServiceGraph(trace)
  const { pos, svgW, svgH } = layoutGraph(nodes, edges)

  return (
    <div style={{
      position: "fixed", left: 12, top: "50%", transform: "translateY(-50%)",
      zIndex: 9998,
      background: "rgba(15,23,42,0.97)",
      border: "1px solid #334155",
      borderRadius: 8,
      padding: "10px 12px",
      fontFamily: "monospace",
      boxShadow: "0 8px 32px rgba(0,0,0,0.6)",
    }}>
      <div style={{ color: "#7dd3fc", fontWeight: "bold", fontSize: 11, marginBottom: 8 }}>
        service graph
      </div>
      <svg width={svgW} height={svgH} style={{ display: "block", overflow: "visible" }}>
        <defs>
          <marker id="sg-arrow" markerWidth="7" markerHeight="5" refX="6" refY="2.5" orient="auto">
            <polygon points="0 0, 7 2.5, 0 5" fill="#0ea5e9" />
          </marker>
        </defs>
        {edges.map((e, i) => {
          const f = pos.get(e.from)
          const t = pos.get(e.to)
          if (!f || !t) return null
          const x1 = f.x + NODE_W, y1 = f.y + NODE_H / 2
          const x2 = t.x - 1,     y2 = t.y + NODE_H / 2
          const mx = (x1 + x2) / 2
          return (
            <path key={i}
              d={`M ${x1} ${y1} C ${mx} ${y1}, ${mx} ${y2}, ${x2} ${y2}`}
              fill="none" stroke="#0ea5e9" strokeWidth={1.5}
              markerEnd="url(#sg-arrow)" opacity={0.75}
            />
          )
        })}
        {nodes.map(n => {
          const p = pos.get(n)
          if (!p) return null
          const label = n.length > 16 ? n.slice(0, 15) + "…" : n
          return (
            <g key={n}>
              <rect x={p.x} y={p.y} width={NODE_W} height={NODE_H}
                rx={6} fill="rgba(30,41,59,0.95)" stroke="#475569" strokeWidth={1} />
              <text x={p.x + NODE_W / 2} y={p.y + NODE_H / 2 + 4}
                textAnchor="middle" fill="#7dd3fc" fontSize={11} fontFamily="monospace">
                {label}
              </text>
            </g>
          )
        })}
      </svg>
    </div>
  )
}

interface JaegerProcess {
  serviceName: string
}

interface JaegerTrace {
  traceID: string
  spans: JaegerSpan[]
  processes: Record<string, JaegerProcess>
}

interface JaegerResponse {
  data: JaegerTrace[]
}

function formatDuration(us: number): string {
  if (us < 1000) return `${us}µs`
  if (us < 1_000_000) return `${(us / 1000).toFixed(1)}ms`
  return `${(us / 1_000_000).toFixed(2)}s`
}

function isNoise(warning: string): boolean {
  return warning.includes("clock skew") || warning.includes("invalid parent span")
}

function hasError(span: JaegerSpan): boolean {
  return span.tags.some(t => t.key === "error" && t.value === true) ||
    (span.warnings != null && span.warnings.some(w => !isNoise(w)))
}

export default function TraceDebugPanel() {
  const [panelOpen, setPanelOpen] = useState(false)
  const [trace, setTrace] = useState<JaegerTrace | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastFetch, setLastFetch] = useState<Date | null>(null)
  const hasTrace = useRef(false)
  const traceId = getTraceId()

  const fetchTrace = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const url = `${ServiceEndpoints.jaegerQuery}/api/traces/${traceId}`
      const res = await fetch(url)
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const json: JaegerResponse = await res.json()
      const found = json.data?.[0] ?? null
      setTrace(found)
      if (found && found.spans.length >= 32) hasTrace.current = true
      setLastFetch(new Date())
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }, [traceId])

  // Poll every second until traces are found
  useEffect(() => {
    fetchTrace()
    const interval = setInterval(() => {
      if (hasTrace.current) { clearInterval(interval); return }
      fetchTrace()
    }, 1000)
    return () => clearInterval(interval)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "d" && e.altKey) setPanelOpen(v => !v)
    }
    window.addEventListener("keydown", handler)
    return () => window.removeEventListener("keydown", handler)
  }, [])

  const spanRows = trace
    ? [...trace.spans].sort((a, b) => a.startTime - b.startTime)
    : []

  const totalDuration = trace
    ? Math.max(...trace.spans.map(s => s.startTime + s.duration)) -
      Math.min(...trace.spans.map(s => s.startTime))
    : 0

  const services = trace
    ? [...new Set(Object.values(trace.processes).map(p => p.serviceName))]
    : []

  const hasTraces = trace !== null && spanRows.length > 0
  const hasErrors = spanRows.some(hasError)
  const jaegerUiUrl = `${ServiceEndpoints.jaegerQuery}/trace/${traceId}`

  return (
    <>
      {/* Prominent trigger button — rendered in document flow below the CV card */}
      <button
        className={`trace-trigger${hasTraces ? " ready" : ""}`}
        onClick={() => setPanelOpen(v => !v)}
        title="Alt+D to toggle"
      >
        <span className={`trace-dot${hasTraces ? " live" : ""}`} />
        {loading && !trace && "fetching traces…"}
        {!loading && !hasTraces && !error && "waiting for traces"}
        {error && <span style={{ color: "#f87171" }}>trace unavailable</span>}
        {hasTraces && (
          <>
            <span style={{ color: "#38bdf8", fontWeight: "bold" }}>
              {spanRows.length} spans
            </span>
            <span style={{ color: "#64748b" }}>·</span>
            <span>{services.join(", ")}</span>
            <span style={{ color: "#64748b" }}>·</span>
            <span style={{ color: hasErrors ? "#f87171" : "#4ade80" }}>
              {formatDuration(totalDuration)}
            </span>
            {hasErrors && <span style={{ color: "#f87171" }}>· errors</span>}
          </>
        )}
      </button>

      {/* Service graph — left side */}
      {panelOpen && trace && <ServiceGraph trace={trace} />}

      {/* Detail panel overlay */}
      {panelOpen && (
        <div style={{
          position: "fixed", bottom: 12, right: 12, zIndex: 9999,
          background: "rgba(15,23,42,0.97)", color: "#e2e8f0",
          border: "1px solid #334155", borderRadius: 8,
          fontFamily: "monospace", fontSize: 11,
          width: 540, maxHeight: "80vh",
          display: "flex", flexDirection: "column",
          boxShadow: "0 8px 32px rgba(0,0,0,0.6)",
        }}>
          {/* Header */}
          <div style={{
            display: "flex", alignItems: "center", justifyContent: "space-between",
            padding: "8px 12px", borderBottom: "1px solid #334155",
            background: "rgba(30,41,59,0.9)", borderRadius: "8px 8px 0 0",
          }}>
            <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
              <span style={{ color: "#7dd3fc", fontWeight: "bold" }}>trace debug</span>
              {trace && (
                <span style={{ color: "#64748b" }}>
                  {spanRows.length} spans · {services.join(", ")} · {formatDuration(totalDuration)}
                </span>
              )}
            </div>
            <div style={{ display: "flex", gap: 6 }}>
              <button onClick={fetchTrace} disabled={loading} style={{
                background: "none", border: "1px solid #475569", color: "#94a3b8",
                borderRadius: 4, padding: "2px 7px", cursor: "pointer", fontSize: 10,
              }}>
                {loading ? "…" : "refresh"}
              </button>
              <a href={jaegerUiUrl} target="_blank" rel="noreferrer" style={{
                background: "none", border: "1px solid #475569", color: "#94a3b8",
                borderRadius: 4, padding: "2px 7px", fontSize: 10, textDecoration: "none",
              }}>
                jaeger ↗
              </a>
              <button onClick={() => setPanelOpen(false)} style={{
                background: "none", border: "none", color: "#94a3b8",
                cursor: "pointer", fontSize: 14, lineHeight: 1, padding: "0 2px",
              }}>×</button>
            </div>
          </div>

          {/* Trace ID */}
          <div style={{ padding: "6px 12px", borderBottom: "1px solid #1e293b", color: "#64748b" }}>
            <span style={{ color: "#475569" }}>id: </span>
            <span style={{ color: "#a5b4fc" }}>{traceId}</span>
            {lastFetch && (
              <span style={{ color: "#334155", float: "right" }}>
                {lastFetch.toLocaleTimeString()}
              </span>
            )}
          </div>

          {/* Error */}
          {error && (
            <div style={{ padding: "8px 12px", color: "#f87171", borderBottom: "1px solid #1e293b" }}>
              {error}
            </div>
          )}

          {/* Spans */}
          <div style={{ overflowY: "auto", flex: 1 }}>
            {spanRows.length === 0 && !loading && !error && (
              <div style={{ padding: "16px 12px", color: "#475569", textAlign: "center" }}>
                no spans yet — they may still be flushing
              </div>
            )}
            {spanRows.map(span => {
              const service = trace!.processes[span.processID]?.serviceName ?? span.processID
              const err = hasError(span)
              const httpStatus = span.tags.find(t => t.key === "http.status_code")?.value
              const httpUrl = span.tags.find(t => t.key === "http.url")?.value as string | undefined
              return (
                <div key={span.spanID} style={{
                  padding: "5px 12px",
                  borderBottom: "1px solid #1e293b",
                  borderLeft: `3px solid ${err ? "#ef4444" : "#334155"}`,
                }}>
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "baseline" }}>
                    <span style={{ color: err ? "#f87171" : "#7dd3fc" }}>{span.operationName}</span>
                    <span style={{ color: "#64748b" }}>{formatDuration(span.duration)}</span>
                  </div>
                  <div style={{ color: "#475569", marginTop: 2 }}>
                    <span style={{ color: "#334155" }}>{service}</span>
                    {httpStatus !== undefined && (
                      <span style={{ marginLeft: 8, color: Number(httpStatus) >= 400 ? "#f87171" : "#4ade80" }}>
                        HTTP {String(httpStatus)}
                      </span>
                    )}
                    {httpUrl && (
                      <span style={{
                        marginLeft: 8, color: "#475569",
                        overflow: "hidden", textOverflow: "ellipsis",
                        display: "inline-block", maxWidth: 320, verticalAlign: "bottom",
                        whiteSpace: "nowrap",
                      }}>
                        {httpUrl.replace(/^https?:\/\/[^/]+/, "")}
                      </span>
                    )}
                    {err && span.warnings && span.warnings.filter(w => !isNoise(w)).map((w, i) => (
                      <span key={i} style={{ marginLeft: 8, color: "#fbbf24" }}>{w}</span>
                    ))}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      )}
    </>
  )
}
