import { fromJson, type MessageShape } from "@bufbuild/protobuf"
import { IdentitySchema } from "@cv/proto/identity/v1/identity_pb"
import { ServiceEndpoints } from "./endpoints"

export type Identity = MessageShape<typeof IdentitySchema>

export async function fetchIdentity({ id }: { id: string }): Promise<Identity[]> {
    const url = `${ServiceEndpoints.identity}/${id}`
    const res = await fetch(url)

    if (!res.ok) {
        throw new Error(`Failed to fetch identity`)
    }

    const json = await res.json()

    return fromJson(IdentitySchema, json)
}