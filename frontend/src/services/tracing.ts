import { WebTracerProvider, StackContextManager, BatchSpanProcessor } from '@opentelemetry/sdk-trace-web'
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch'
import { registerInstrumentations } from '@opentelemetry/instrumentation'
import { context, trace, Span } from '@opentelemetry/api'
import { resourceFromAttributes } from '@opentelemetry/resources'
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions'
import { ServiceEndpoints } from './endpoints'

const exporter = new OTLPTraceExporter({
  url: ServiceEndpoints.otlp,
})

const provider = new WebTracerProvider({
  resource: resourceFromAttributes({ [ATTR_SERVICE_NAME]: 'cv-frontend' }),
  spanProcessors: [new BatchSpanProcessor(exporter, { scheduledDelayMillis: 1000 })],
})

provider.register({
  contextManager: new StackContextManager(),
})

registerInstrumentations({
  instrumentations: [
    new FetchInstrumentation({
      propagateTraceHeaderCorsUrls: [/localhost/],
    }),
  ],
})

const tracer = trace.getTracer('cv-frontend')

const pageLoadSpan: Span = tracer.startSpan('page-load')
export const pageLoadContext = trace.setSpan(context.active(), pageLoadSpan)

export function tracedFetch(url: string, init?: RequestInit): Promise<Response> {
  return context.with(pageLoadContext, () => fetch(url, init))
}

export function getTraceId(): string {
  return pageLoadSpan.spanContext().traceId
}

// End the page-load span once the page has finished loading so it is exported
// before child spans arrive at Jaeger, preventing "invalid parent span IDs" warnings.
window.addEventListener('load', () => pageLoadSpan.end(), { once: true })