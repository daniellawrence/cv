type LoadingProps = {
  text?: string
}

export default function Loading({ text = "loading" }: LoadingProps) {
  return (
    <div
      className="redacted-script-regular"
      style={{
        opacity: 0.6,
        animation: "pulse 1.2s ease-in-out infinite",
      }}
    >
      {text}
    </div>
  )}