import { fromJson, type MessageShape } from "@bufbuild/protobuf"
import { InterestSchema } from "@cv/proto/interest/v1/interest_pb"
import { ServiceEndpoints } from "./endpoints";

export type Interest = MessageShape<typeof InterestSchema>

export async function fetchInterest({ id }: { id: string }): Promise<Interest> {
    const url = `${ServiceEndpoints.interest}/${id}`
    const res = await fetch(url)

    if (!res.ok) {
        throw new Error(`Failed to fetch interest`)
    }

    const json = await res.json()
    // TODO: change the pagination to match the limit and then use pages.
    return fromJson(InterestSchema, json)
}