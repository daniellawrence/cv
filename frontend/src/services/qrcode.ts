import { fromJson, type MessageShape } from "@bufbuild/protobuf"
import { GenerateQRCodeResponseSchema } from "@cv/proto/qrcode/v1/qrcode_pb"
import { ServiceEndpoints } from "./endpoints";

export type Qrcode = MessageShape<typeof GenerateQRCodeResponseSchema>

export async function fetchQrcode({ encode_url }: { encode_url: string }): Promise<Qrcode> {
    const url = `${ServiceEndpoints.qrcode}?url=${encode_url}`
    const res = await fetch(url)

    if (!res.ok) {
        throw new Error(`Failed to fetch qrcode`)
    }

    const json = await res.json()
    // TODO: change the pagination to match the limit and then use pages.
    return fromJson(GenerateQRCodeResponseSchema, json)
}