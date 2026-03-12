import { useEffect, useState } from "react"
import { fetchInterest } from "../services/interest"
import type { Interest } from "@cv/proto/interest/v1/interest_pb"

export function InterestList({ id }: { id: string }) {
  const [interest, setInterest] = useState<Interest | null>(null)

  useEffect(() => {
    fetchInterest({ id }).then(setInterest)
  }, [id])

  if (!interest) {
    return <div>Loading...</div>
  }

  const { names } = interest

  return (
      <div className="skills-section">

      <div className="skill-category">{interest.type}</div>
        <div className="skill-text">
        {names.map((n, i) => (
          <span key={n} className="tech-stack-skill">{n}{i < names.length - 1 ? ", " : ""}</span>
        ))}
        </div>
      </div>
  )
}



export default function Interest() {
  return (
    <section>
      <div className="section-title">Talk to me about</div>
      <InterestList id="tech" />
      <InterestList id="skills" />
      <InterestList id="hobbies" />
    </section>
  )
}