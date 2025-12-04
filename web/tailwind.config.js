/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "#0f172a", // Slate 900
        surface: "#1e293b",    // Slate 800
        primary: "#3b82f6",    // Blue 500
        secondary: "#64748b",  // Slate 500
        accent: {
          wakatime: "#8b5cf6", // Violet 500
          googlefit: "#10b981", // Emerald 500
          calendar: "#0ea5e9", // Sky 500
          activity: "#f59e0b", // Amber 500
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
