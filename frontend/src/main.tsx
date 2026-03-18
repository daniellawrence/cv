// Global error handling for tests - log errors without throwing
window.errorLog = window.errorLog || [];
const originalOnError = window.onerror || (() => false);
window.onerror = (msg, source, line, column, error) => {
  const errorMsg = msg || String(error?.message);
  if (!errorMsg.includes('Failed to load resource')) {
    window.errorLog.push(errorMsg);
  }
  return originalOnError(msg, source, line, column, error);
};

import "./services/tracing"
import React from "react"
import ReactDOM from "react-dom/client"
import App from "./App"
import "./styles/main.css"

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
