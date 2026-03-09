import { useEffect, useState } from "react"
import EducationList from "./components/educationList"
import ExperienceList from "./components/expereienceList"

export default function App() {


  return (
    <>
        <div className="container">
            <div className="header">
                <div className="header-left">
                    <h1>Daniel Lawrence</h1>
                </div>
            </div>

            <div className="left-column">
                <div className="section-title">Experience (past 10 years)</div>
                <ExperienceList limit={5} offset={0} />
            </div>

            <div className="right-column">
                <div className="section-title">Experience (cont.)</div>
                <ExperienceList limit={5} offset={5} />
                <EducationList />
            </div>
        </div>
    </>
  )
}