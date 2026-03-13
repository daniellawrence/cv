const env = import.meta.env

const runtime = (window as { RUNTIME_CONFIG?: Record<string, string> }).RUNTIME_CONFIG ?? {}

export const ServiceEndpoints = {
    experience:
        runtime.EXPERIENCE_SERVICE_URL ??
        env.EXPERIENCE_SERVICE_URL ??
        "http://experience.localhost/experience",

    education:
        runtime.EDUCATION_SERVICE_URL ??
        env.EDUCATION_SERVICE_URL ??
        "http://education.localhost/education",

    interest:
        runtime.INTEREST_SERVICE_URL ??
        env.INTEREST_SERVICE_URL ??
        "http://interest.localhost/interest",

    identity:
        runtime.IDENTITY_SERVICE_URL ??
        env.IDENTITY_SERVICE_URL ??
        "http://identity.localhost/identity",

    qrcode:
        runtime.QRCODE_SERVICE_URL ??
        env.QRCODE_SERVICE_URL ??
        "http://qrcode.localhost/qrcode",

    otlp:
        runtime.OTEL_EXPORTER_OTLP_ENDPOINT ??
        env.OTEL_EXPORTER_OTLP_ENDPOINT ??
        "http://tracing.localhost/v1/traces",

    jaegerQuery:
        runtime.JAEGER_QUERY_URL ??
        env.JAEGER_QUERY_URL ??
        "http://tracing.localhost",
}