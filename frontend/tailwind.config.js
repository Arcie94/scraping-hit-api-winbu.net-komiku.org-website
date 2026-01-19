/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                background: '#0a0a0c', // Deep dark
                primary: '#3b82f6',    // Neon Blue
                secondary: '#1e293b',  // Card bg
                accent: '#facc15',     // Gold/Yellow for badges
            },
        },
    },
    plugins: [],
}
