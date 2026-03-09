import { fromJson, type MessageShape } from "@bufbuild/protobuf"
import { ExperienceSchema } from "@cv/proto/experience/v1/experience_pb"

export type Experience = MessageShape<typeof ExperienceSchema>

export async function fetchExperience({ limit, offset }: { limit: number; offset: number }): Promise<Experience[]> {
    const res = await fetch("http://localhost:8080/experience")

    if (!res.ok) {
        throw new Error(`Failed to fetch experience`)
    }

    const json = await res.json()
    // TODO: change the pagination to match the limit and then use pages.
    const experience = json.experience.slice(offset, offset + limit)
    return experience.map(e => fromJson(ExperienceSchema, e))
}