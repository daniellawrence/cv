import { useEffect, useState } from "react"
import EducationList from "./components/educationList"
import ExperienceList from "./components/expereienceList"
import InterestList from "./components/interestList"
import Header from "./components/header"

export default function App() {


  return (
    <>
        <div className="container">
            <Header id="dsl" />
            
            <div className="left-column">
                <div className="section-title">Experience (past 10 years)</div>
                <ExperienceList limit={4} offset={0} />
            </div>

            <div className="right-column">
                <div className="section-title">Experience (cont.)</div>
                <ExperienceList limit={3} offset={4} />
                <EducationList />
                <InterestList id="tech" />
                <InterestList id="skills" />
                <InterestList id="hobbies" />
            </div>
        </div>
    </>
  )
}