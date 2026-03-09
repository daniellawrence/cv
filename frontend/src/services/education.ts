import { fromJson, type MessageShape } from "@bufbuild/protobuf"
import { EducationSchema } from "@cv/proto/education/v1/education_pb"
import { ServiceEndpoints } from "./endpoints"

export type Education = MessageShape<typeof EducationSchema>

export async function fetchEducation(): Promise<Education[]> {
    const res = await fetch(ServiceEndpoints.education)

    if (!res.ok) {
        throw new Error(`Failed to fetch education`)
    }

    const json = await res.json()

    return json.education.map(e => fromJson(EducationSchema, e))
}