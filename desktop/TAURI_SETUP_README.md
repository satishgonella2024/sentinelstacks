# SentinelStacks Desktop Tauri 2.0 Setup Guide

## Issues Found

When setting up the SentinelStacks desktop application with Tauri 2.0, we encountered the following issues:

1. **Plugin Feature Flags**: The feature flags specified for some Tauri plugins don't exist in the current versions
2. **Configuration Structure**: The configuration structure needed updates to match Tauri 2.0 requirements
3. **Plugin Initialization**: Plugins needed to be properly initialized in the Rust code

## Changes Made

We've made the following changes to fix these issues:

### 1. Updated Cargo.toml

Removed non-existent feature flags from the dependencies:
```toml
[dependencies]
tauri = { version = "2.0.0", features = [] }
tauri-plugin-shell = "2.0.0"  # Removed features like execute, open, sidecar
tauri-plugin-dialog = "2.0.0"
tauri-plugin-fs = "2.0.0"     # Removed the "all" feature
tauri-plugin-process = "2.0.0"
```

### 2. Updated tauri.conf.json

Modified the plugin configuration to match available capabilities:
```json
"plugins": {
  "shell": {
    "open": true
  },
  "dialog": {
    "all": true
  },
  "fs": {
    "scope": {
      "allow": ["$APP/*", "$CONFIG/*"]
    }
  }
}
```

### 3. Updated lib.rs

Ensured that plugins are properly initialized:
```rust
tauri::Builder::default()
    .plugin(tauri_plugin_shell::init())
    .plugin(tauri_plugin_dialog::init())
    .plugin(tauri_plugin_fs::init())
    .plugin(tauri_plugin_process::init())
```

## How to Apply the Fix

1. Run the updated fix script:
   ```bash
   chmod +x ./fix-tauri.sh
   ./fix-tauri.sh
   ```

2. Start the Tauri application:
   ```bash
   npm run tauri dev
   ```

## Tauri 2.0 Changes

Tauri 2.0 introduced significant changes from previous versions:

1. **Plugin System**: Now uses separate plugin crates with different API structures
2. **Configuration Format**: Changed from `allowlist` to `plugins` section
3. **Feature Flags**: Many feature flags were removed or changed
4. **Mobile Support**: Added support for mobile development

## Additional Resources

- [Tauri 2.0 Documentation](https://tauri.app/v2/guides/)
- [Tauri 2.0 Plugin API](https://docs.rs/tauri-plugin/2.0.0/tauri_plugin/)
- [Tauri 2.0 Migration Guide](https://tauri.app/v2/migration/)