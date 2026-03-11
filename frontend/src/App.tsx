import EducationList from "./components/educationList"
import ExperienceList from "./components/expereienceList"
import Interest from "./components/interestList"
import Header from "./components/header"
import GeometricBackground from "./components/wave"
import TraceDebugPanel from "./components/traceDebugPanel"

export default function App() {


  return (
    <>
        <GeometricBackground />
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
                <Interest />
            </div>
        </div>
        <TraceDebugPanel />
    </>
  )
}