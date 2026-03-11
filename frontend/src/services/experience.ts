import { fromJson, type MessageShape } from "@bufbuild/protobuf"
import { ExperienceSchema } from "@cv/proto/experience/v1/experience_pb"
import { ServiceEndpoints } from "./endpoints";

export type Experience = MessageShape<typeof ExperienceSchema>

export async function fetchExperience({ limit, offset }: { limit: number; offset: number }): Promise<Experience[]> {
    const url = `${ServiceEndpoints.experience}/${offset}/${limit}`
    const res = await fetch(url)

    if (!res.ok) {
        throw new Error(`Failed to fetch experience`)
    }

    const json = await res.json()
    return json.experience.map(e => fromJson(ExperienceSchema, e))
}