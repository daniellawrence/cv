import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import path from "node:path"
import { fileURLToPath } from "node:url"

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            proto: path.resolve(__dirname, "../gen/ts")
        },
        preserveSymlinks: true,
    },
    optimizeDeps: {
        include: [
            "@bufbuild/protobuf",
            "@bufbuild/protobuf/codegenv2",
        ],
    },
    server: {
        fs: {
            allow: [".."]
        }
    }
})