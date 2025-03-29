# Quick Fix for PostCSS and Vite Setup Issues

If you're encountering issues with PostCSS configuration in your Vite setup, follow these steps to fix them:

## Using the Fix Script

1. Make the script executable:
   ```bash
   chmod +x ./fix-vite-setup.sh
   ```

2. Run the script:
   ```bash
   ./fix-vite-setup.sh
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

## What the Script Does

The script fixes several common issues:

1. Adds `"type": "module"` to package.json to enable ES module syntax
2. Creates a proper PostCSS configuration using ES module syntax
3. Updates the Vite configuration to use the correct PostCSS file
4. Ensures Tailwind is configured correctly with ES module syntax
5. Installs necessary dependencies

## Manual Fix

If you prefer to fix the issues manually:

1. Add `"type": "module"` to package.json
   ```json
   {
     "name": "sentinel-desktop",
     "version": "0.1.0",
     "private": true,
     "type": "module",
     // rest of your package.json
   }
   ```

2. Create/update postcss.config.js with ES module syntax:
   ```js
   export default {
     plugins: {
       tailwindcss: {},
       autoprefixer: {},
     }
   }
   ```

3. Make sure your vite.config.ts references the correct file:
   ```ts
   css: {
     postcss: './postcss.config.js',
   }
   ```

4. Verify tailwind.config.js uses ES module syntax:
   ```js
   export default {
     // Tailwind configuration
   }
   ```

After making these changes, restart your development server with `npm run dev`.
