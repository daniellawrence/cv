import { useEffect, useState } from "react"
import { fetchInterest, type Interest } from "../services/interest"

export function InterestList({ id }: { id: string }) {
  const [interest, setInterest] = useState<Interest | null>(null)

  useEffect(() => {
    fetchInterest({ id }).then(setInterest)
  }, [id])

  if (!interest) {
    return <div>Loading...</div>
  }

  return (
      <div className="skills-section">

      <div className="skill-category">{interest.type}</div>
        <div className="skill-text">
        {interest.names.map((n) => (     
          <span key={n} className="tech-stack-skill">{n}, </span>
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