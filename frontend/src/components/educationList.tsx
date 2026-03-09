import { useEffect, useState } from "react"
import { fetchEducation, type Education } from "../services/education"

export default function Education() {
  const [education, setEducation] = useState<Education[]>([])

  useEffect(() => {
    fetchEducation().then(setEducation)
  }, [])

  return (
    <section>
      <h2>Education</h2>

      {education.map((e) => (
        <div key={e.id}>
          <strong>{e.institution}</strong>
          <div>{e.degree}</div>
          <div>{e.startDate} – {e.endDate}</div>
        </div>
      ))}
    </section>
  )
}