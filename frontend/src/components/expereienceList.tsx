import { useEffect, useState } from "react"
import { fetchExperience, type Experience } from "../services/experience"

export default function ExperienceList({limit, offset}: {limit: number; offset: number}) {
  const [experience, setExperience] = useState<Experience[]>([])

  useEffect(() => {
    fetchExperience({limit, offset}).then(setExperience)
  }, [])


  return (
    <section>
      <div className="section-title">Experience</div>

      {experience.map((e) => (         
            <div class="experience-item" key={e.id}>
                <div class="company-name">{e.company}</div>
                <div class="job-title">{e.title}</div>
                <div class="job-dates">{e.startDate} – {e.endDate}</div>
                <div class="location">{e.location}</div>
                <div class="job-bullets">
                    {e.highlights.map((h) => (     
                    <div class="bullet">{h}</div>
                    ))}
                </div>
                <div class="tech-stack">
                    {e.skills.map((k) => (     
                    <span class="tech-stack-skill">{k}, </span>
                    ))}
                </div>
            </div>
      ))}
    </section>
  )
}