#!/bin/bash
set -e

echo "🎨 Fixing SentinelStacks UI styling issues..."

# Create the public directory if it doesn't exist
mkdir -p public

# Initialize MSW service worker if needed
if [ ! -f "public/mockServiceWorker.js" ]; then
  echo "📝 Setting up Mock Service Worker..."
  npx msw init public/ --save
fi

# Apply additional styling fix: Set body background color directly
cat > public/additional-styles.css << 'EOF'
body {
  background-color: #0a1017;
  color: white;
}
EOF

# Add the additional stylesheet to the HTML
if ! grep -q "additional-styles.css" "../index.html"; then
  echo "📄 Adding additional stylesheet to index.html..."
  sed -i '' 's/<\/head>/<link rel="stylesheet" href="\/additional-styles.css" \/><\/head>/' ../index.html
fi

# Run the application
echo "🚀 Starting the application..."
cd ..
./scripts/dev.sh
