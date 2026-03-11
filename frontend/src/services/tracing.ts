function randomHex(bytes: number): string {
  return Array.from(crypto.getRandomValues(new Uint8Array(bytes)))
    .map(b => b.toString(16).padStart(2, "0"))
    .join("")
}

const pageTraceId = randomHex(16) // reset on every page load

export function tracedFetch(url: string, init?: RequestInit): Promise<Response> {
  const traceId = pageTraceId
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
