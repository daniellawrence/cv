const SESSION_TRACE_KEY = "cv_trace_id"

function randomHex(bytes: number): string {
  return Array.from(crypto.getRandomValues(new Uint8Array(bytes)))
    .map(b => b.toString(16).padStart(2, "0"))
    .join("")
}

function getSessionTraceId(): string {
  let id = sessionStorage.getItem(SESSION_TRACE_KEY)
  if (!id) {
    id = randomHex(16) // 32 hex chars = 128-bit trace id
    sessionStorage.setItem(SESSION_TRACE_KEY, id)
  }
  return id
}

export function tracedFetch(url: string, init?: RequestInit): Promise<Response> {
  const traceId = getSessionTraceId()
  const spanId = randomHex(8) // 16 hex chars = 64-bit span id
  const traceparent = `00-${traceId}-${spanId}-01`
  console.log(traceId, spanId, traceparent)

  return fetch(url, {
    ...init,
    headers: {
      ...init?.headers,
      traceparent,
    },
  })
}
