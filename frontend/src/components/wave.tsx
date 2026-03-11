// This was inspired by https://www.vantajs.com/
import { useEffect, useState } from "react"

export default function GeometricBackground() {
  const [t, setT] = useState(0)

  useEffect(() => {
    let frame: number

    const animate = () => {
      setT(v => v + 0.002)
      frame = requestAnimationFrame(animate)
    }

    animate()
    return () => cancelAnimationFrame(frame)
  }, [])

  const xs = [0, 160, 320, 480, 640, 800, 960, 1120, 1280, 1440]

  const buildPolygon = (base: number, amp: number, phase: number) => {
    const points = xs
      .map((x, i) => {
        const y = base + Math.sin(t * 0.8 + i * 0.6 + phase) * amp
        return `${x},${y}`
      })
      .join(" ")

    return `${points} 1440,900 0,900`
  }

  const poly1 = buildPolygon(450, 35, 0)
  const poly2 = buildPolygon(620, 40, 1.5)

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        zIndex: -1,
        background: "#0b1220",
      }}
    >
      <svg
        viewBox="0 0 1440 900"
        preserveAspectRatio="none"
        style={{ width: "100%", height: "100%" }}
      >
        <polygon fill="#1f2b4d" points={poly1} />
        <polygon fill="#26345a" opacity="0.8" points={poly2} />
      </svg>
    </div>
  )
}