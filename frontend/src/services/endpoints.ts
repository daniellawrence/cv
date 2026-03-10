const env = import.meta.env

const runtime = (window as any).RUNTIME_CONFIG ?? {}

export const ServiceEndpoints = {
    experience:
        runtime.EXPERIENCE_SERVICE_URL ??
        env.EXPERIENCE_SERVICE_URL ??
        "http://localhost:8080/experience",

    education:
        runtime.EDUCATION_SERVICE_URL ??
        env.EDUCATION_SERVICE_URL ??
        "http://localhost:8081/education",

    interest:
        runtime.INTEREST_SERVICE_URL ??
        env.INTEREST_SERVICE_URL ??
        "http://localhost:8082/interest",

    identity:
        runtime.IDENTITY_SERVICE_URL ??
        env.IDENTITY_SERVICE_URL ??
        "http://localhost:8083/identity",

    qrcode:
        runtime.QRCODE_SERVICE_URL ??
        env.QRCODE_SERVICE_URL ??
        "http://localhost:8084/qrcode",
}