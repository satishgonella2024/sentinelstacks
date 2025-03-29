# Complete Tailwind CSS Troubleshooting Guide

## The Issues

You're encountering two main issues with Tailwind CSS v4 in your SentinelStacks desktop application:

1. **PostCSS Plugin Error**: Tailwind CSS v4 has moved its PostCSS plugin to a separate package.
2. **Missing Color Utilities**: Tailwind CSS v4 has different default color palettes.

## Complete Fix Instructions

### Step 1: Install the Required Package

The most critical step is to install the separate `@tailwindcss/postcss` package:

```bash
# Navigate to the desktop directory
cd /Users/subrahmanyagonella/SentinelStacks/desktop

# Install the required package
npm install -D @tailwindcss/postcss
```

Or use the provided script:

```bash
chmod +x ./fix-tailwind.sh
./fix-tailwind.sh
```

### Step 2: Verify Configuration Files

We've updated these files for you, but double-check that they contain the correct configurations:

1. **postcss.config.cjs**:
```js
module.exports = {
  plugins: {
    '@tailwindcss/postcss': {},
    autoprefixer: {},
  }
}
```

2. **tailwind.config.js**:
```js
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
          // ... other primary colors
        },
        gray: {
          50: '#f9fafb',
          // ... other gray colors
        },
      },
    },
  },
  plugins: [],
}
```

### Step 3: Clear Cache and Reinstall (If Needed)

If you still encounter issues after the above steps, try the following:

```bash
# Navigate to desktop directory
cd /Users/subrahmanyagonella/SentinelStacks/desktop

# Remove node_modules and package-lock.json
rm -rf node_modules
rm package-lock.json

# Reinstall dependencies
npm install

# Start the development server
npm run dev
```

## Troubleshooting Specific Errors

### Error: "tailwindcss directly as a PostCSS plugin"
- Make sure you've installed `@tailwindcss/postcss`
- Verify that your postcss.config.cjs is using CommonJS syntax (`module.exports`)
- Try deleting any other PostCSS config files (like postcss.config.js)

### Error: "Cannot apply unknown utility class"
- Check if the utility is part of Tailwind CSS v4
- Add the missing color or utility to your tailwind.config.js
- Update your CSS to use the new Tailwind CSS v4 syntax

## Additional Notes

- Tailwind CSS v4 has significant changes from v3
- Always check the Tailwind CSS v4 documentation for updated syntax
- Consider using the Tailwind CSS v4 migration guide for more complex projects

## Final Steps

After making all these changes, restart your development server:

```bash
npm run dev
```