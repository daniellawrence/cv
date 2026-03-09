import { useEffect, useState } from "react"
import { fetchIdentity, type Identity } from "../services/identity"
import QRcode from "./qrcode"

export default function Header({ id }: { id: string }) {
  const [identity, setIdentity] = useState<Identity | null>(null)

  useEffect(() => {
    fetchIdentity({ id }).then(setIdentity)
  }, [id])

  if (!identity) {
    return <div>Loading...</div>
  }
  

  return (
    <>
        <div className="header">
            <div className="header-left">
                <h1>{identity.name}</h1>
                <p><a href="mailto: {identity.email}">{identity.email}</a></p>
            </div>
            <QRcode encode_url={identity.linkedin} />
        </div>
    </>
  )
}