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
      if (found && found.spans.length >= 30) hasTrace.current = true
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
        onClick={() => setPanelOpen(true)}
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
