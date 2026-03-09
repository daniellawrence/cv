import { useEffect, useState } from "react"
import { fetchQrcode, type Qrcode } from "../services/qrcode"

export default function QRcode({ encode_url }: { encode_url: string }) {
  const [qrcode, setQrcode] = useState<Qrcode | null>(null)

  useEffect(() => {
    fetchQrcode({ encode_url }).then(setQrcode)
  }, [encode_url])

  if (!qrcode) {
    return <div>Loading...</div>
  }

  return (
    <>
      <div className="qr-code" id="qrcode">
        <a href={encode_url}>
          <img src={`data:image/png;base64,${qrcode.imageBase64}`} width={150} height={150} />          
          <small>{encode_url}</small>
        </a>
      </div>
    </>
  )
}