import { useEffect, useState } from "react"
import EducationList from "./components/educationList"

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
            </div>

            <div className="right-column">
                <div className="section-title">Experience (cont.)</div>
                <EducationList />
            </div>
        </div>
    </>
  )
}