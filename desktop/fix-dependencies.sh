#!/bin/bash

# Make sure we're in the correct directory
cd "$(dirname "$0")"

# Install necessary dependencies
echo "Installing dependencies..."
npm install --save react-router-dom @heroicons/react recharts react-hot-toast
npm install --save-dev @vitejs/plugin-react vite postcss-cli postcss autoprefixer tailwindcss

# Fix PostCSS config
echo "Fixing PostCSS configuration..."
cat > postcss.config.cjs << EOF
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
EOF

# Fix Tailwind config
echo "Checking Tailwind configuration..."
if [ ! -f tailwind.config.js ]; then
  echo "Creating tailwind.config.js..."
  cat > tailwind.config.js << EOF
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          200: '#bae6fd',
          300: '#7dd3fc',
          400: '#38bdf8',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
          800: '#075985',
          900: '#0c4a6e',
          950: '#082f49',
        },
      },
    },
  },
  plugins: [],
}
EOF
fi

# Make sure the root package has the correct postcss config
cd ..
cat > postcss.config.cjs << EOF
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
EOF

# Return to desktop directory
cd desktop

echo "Dependencies fixed! You can now run 'npm run dev' to start the development server."
