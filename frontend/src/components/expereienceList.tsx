import { useEffect, useState } from "react"
import { fetchExperience } from "../services/experience"
import type { Experience } from "@cv/proto/experience/v1/experience_pb"

export default function ExperienceList({limit, offset}: {limit: number; offset: number}) {
  const [experience, setExperience] = useState<Experience[]>([])

  useEffect(() => {
    fetchExperience({limit, offset}).then(setExperience)
  }, [])


  return (
    <section>
      {experience.map((e) => (         
            <div className="experience-item" key={e.id}>
                <div className="company-name">{e.company}</div>
                <div className="job-title">{e.title}</div>
                <div className="job-dates">{e.startDate} – {e.endDate}</div>
                <div className="location">{e.location}</div>
                <div className="job-bullets">
                    {e.highlights.map((h) => (     
                    <div key={h} className="bullet">{h}</div>
                    ))}
                </div>
                <div className="tech-stack">
                    {e.skills.map((k, i) => (
                    <span key={k} className="tech-stack-skill">{k}{i < e.skills.length - 1 ? ", " : ""}</span>
                    ))}
                </div>
            </div>
      ))}
    </section>
  )
}