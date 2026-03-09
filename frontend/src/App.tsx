import { useEffect, useState } from "react"
import { fetchEducation } from "./api/education"
import { EducationSchema } from "proto/education/v1/education_pb"
import type { MessageShape } from "@bufbuild/protobuf"

type Education = MessageShape<typeof EducationSchema>

export default function App() {

  const [education, setEducation] = useState<Education[]>([])

  useEffect(() => {
    fetchEducation().then(setEducation)
  }, [])

  return (
    <div style={{ padding: 40 }}>
      <h1>Education</h1>

      {education.map(e => (
        <div key={e.id}>
          <h3>{e.institution}</h3>
          <p>{e.degree}</p>
          <p>{e.startDate} - {e.endDate}</p>
        </div>
      ))}
    </div>
  )
}