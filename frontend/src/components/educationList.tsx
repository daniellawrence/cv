import { useEffect, useState } from "react"
import { fetchEducation, type Education } from "../services/education"

export default function EducationList() {
  const [education, setEducation] = useState<Education[]>([])

  useEffect(() => {
    fetchEducation().then(setEducation)
  }, [])

  return (
    <section>
      <div className="section-title">Education</div>
      {education.length ? (
        education.map((e) => (
            <div className="education-item" key={e.id}>
                <div className="degree">{e.degree}</div>
                <div className="school">{e.institution}</div>
                <div className="year">{e.startDate} – {e.endDate}</div>
            </div>
      ))
      ) : (
            <div className="education-item">
              <div className="degree">xxxxx, xxxxxx xxxxxx</div>
              <div className="school">xxxxxxxx xxxxxxxxx</div>
              <div className="year">xxxx-xx – xxxx-xx</div>
            </div>
      )}

    </section>
  )
}