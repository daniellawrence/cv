const env = import.meta.env

export const ServiceEndpoints = {
    experience: env.EXPERIENCE_SERVICE_URL ?? "http://localhost:8080/experience",
    education: env.EDUCATION_SERVICE_URL ?? "http://localhost:8081/education",
    interest: env.INTEREST_SERVICE_URL ?? "http://localhost:8082/interest",
}