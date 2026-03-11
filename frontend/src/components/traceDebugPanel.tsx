import { useEffect, useState, useCallback } from "react"
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

function hasError(span: JaegerSpan): boolean {
  return span.tags.some(t => t.key === "error" && t.value === true) ||
    (span.warnings != null && span.warnings.length > 0)
}

export default function TraceDebugPanel() {
  const [visible, setVisible] = useState(false)
  const [trace, setTrace] = useState<JaegerTrace | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastFetch, setLastFetch] = useState<Date | null>(null)
  const traceId = getTraceId()

  const fetchTrace = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const url = `${ServiceEndpoints.jaegerQuery}/api/traces/${traceId}`
      const res = await fetch(url)
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const json: JaegerResponse = await res.json()
      setTrace(json.data?.[0] ?? null)
      setLastFetch(new Date())
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }, [traceId])

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "d" && e.altKey) setVisible(v => !v)
    }
    window.addEventListener("keydown", handler)
    return () => window.removeEventListener("keydown", handler)
  }, [])

  useEffect(() => {
    if (visible && !trace) fetchTrace()
  }, [visible, trace, fetchTrace])

  if (!visible) {
    return (
      <button
        onClick={() => setVisible(true)}
        style={{
          position: "fixed", bottom: 12, right: 12, zIndex: 9999,
          background: "rgba(0,0,0,0.6)", color: "#7dd3fc",
          border: "1px solid #0ea5e9", borderRadius: 6,
          padding: "4px 10px", fontSize: 11, cursor: "pointer",
          fontFamily: "monospace",
        }}
        title="Alt+D to toggle"
      >
        trace
      </button>
    )
  }

  const jaegerUiUrl = `${ServiceEndpoints.jaegerQuery}/trace/${traceId}`

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

  return (
    <div style={{
      position: "fixed", bottom: 12, right: 12, zIndex: 9999,
      background: "rgba(15,23,42,0.96)", color: "#e2e8f0",
      border: "1px solid #334155", borderRadius: 8,
      fontFamily: "monospace", fontSize: 11,
      width: 520, maxHeight: "80vh",
      display: "flex", flexDirection: "column",
      boxShadow: "0 8px 32px rgba(0,0,0,0.5)",
    }}>
      {/* Header */}
      <div style={{
        display: "flex", alignItems: "center", justifyContent: "space-between",
        padding: "8px 12px", borderBottom: "1px solid #334155",
        background: "rgba(30,41,59,0.8)", borderRadius: "8px 8px 0 0",
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
          <button onClick={() => setVisible(false)} style={{
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
                    display: "inline-block", maxWidth: 300, verticalAlign: "bottom",
                    whiteSpace: "nowrap",
                  }}>
                    {httpUrl.replace(/^https?:\/\/[^/]+/, "")}
                  </span>
                )}
                {err && span.warnings && (
                  <span style={{ marginLeft: 8, color: "#fbbf24" }}>{span.warnings[0]}</span>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
