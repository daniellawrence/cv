import { useEffect, useState } from "react"
import { fetchQrcode } from "../services/qrcode"
import type { GenerateQRCodeResponse } from "@cv/proto/qrcode/v1/qrcode_pb"

export default function QRcode({ encode_url }: { encode_url: string }) {
  const [qrcode, setQrcode] = useState<GenerateQRCodeResponse | null>(null)
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    fetchQrcode({ encode_url }).then(setQrcode)
  }, [encode_url])

  return (
    <>
      <div className="qr-code" id="qrcode">
        <a href={encode_url}>
          <img
            src={qrcode ? `data:image/png;base64,${qrcode.imageBase64}` : undefined}
            width={150}
            height={150}
            onLoad={() => setVisible(true)}
            style={{ filter: visible ? "blur(0px)" : "blur(8px)", transition: "filter 0.4s ease" }}
          />
        </a>
      </div>
    </>
  )
}