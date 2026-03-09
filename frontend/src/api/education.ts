import { create, type MessageShape } from "@bufbuild/protobuf"
import { EducationSchema } from "@cv/proto/education/v1/education_pb"
import { fromJson } from "@bufbuild/protobuf"

export async function fetchEducation(): Promise<MessageShape<typeof EducationSchema>[]> {
    const res = await fetch("http://localhost:8080/education")

    const json = await res.json()

    return json.map((e: any) => fromJson(EducationSchema, e))
}