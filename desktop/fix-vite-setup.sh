#!/bin/bash

# Make sure we're in the correct directory
cd "$(dirname "$0")"

echo "Fixing Vite + Tailwind setup..."

# Make sure package.json has the proper type
if ! grep -q '"type": "module"' package.json; then
  echo "Adding 'type: module' to package.json..."
  sed -i '' 's/"private": true,/"private": true,\n  "type": "module",/g' package.json
fi

# Fix PostCSS config for ES modules
echo "Creating proper PostCSS config (ES module)..."
cat > postcss.config.js << EOF
// Using ES Module format since package.json has "type": "module"
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  }
}
EOF

# Remove the old CJS version
rm -f postcss.config.cjs

# Fix Vite config to use the JS version
echo "Updating Vite config to use postcss.config.js..."
sed -i '' 's/postcss: .*postcss.config.cjs.*,/postcss: ".\/postcss.config.js",/g' vite.config.ts

# Fix Tailwind config if needed
echo "Checking tailwind.config.js..."
if ! grep -q 'export default' tailwind.config.js; then
  echo "Fixing tailwind.config.js to use ES module syntax..."
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
        gray: {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',
          800: '#1f2937',
          900: '#111827',
          950: '#030712',
        },
      },
    },
  },
  plugins: [],
}
EOF
fi

# Install dependencies if needed
echo "Installing required dependencies..."
npm install --save react-router-dom @heroicons/react recharts react-hot-toast
npm install --save-dev @vitejs/plugin-react vite postcss autoprefixer tailwindcss

echo "Setup complete! Try running 'npm run dev' now."
