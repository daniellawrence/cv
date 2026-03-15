import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"

export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            proto: "./gen/ts"
        },
    },
    optimizeDeps: {
        include: [
            "@bufbuild/protobuf",
            "@bufbuild/protobuf/codegenv2",
        ],
    },
})