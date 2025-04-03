# API Documentation

This page provides interactive documentation for the SentinelStacks API.

## Interactive API Explorer

<div id="swagger-ui"></div>

<script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js"></script>
<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />

<script>
  window.onload = function() {
    window.ui = SwaggerUIBundle({
      url: "../api-reference.yaml",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
      ],
      layout: "StandaloneLayout",
      defaultModelsExpandDepth: -1,
      displayRequestDuration: true,
      defaultModelRendering: 'model',
      showExtensions: true,
      showCommonExtensions: true,
      tagsSorter: 'alpha',
      operationsSorter: 'alpha'
    });
  };
</script>

<style>
  .swagger-ui .topbar { display: none }
</style>

## API Reference

The full OpenAPI specification is available [here](../api-reference.yaml).

### Memory Management API

SentinelStacks provides a comprehensive memory API for storing, retrieving, and searching data:

- **[Store Memory](../api-usage-guide.md#storing-values-in-memory)**: Store key-value pairs in agent memory
- **[Retrieve Memory](../api-usage-guide.md#retrieving-values-from-memory)**: Retrieve values by key
- **[Search Memory](../api-usage-guide.md#semantic-search-in-memory)**: Perform semantic search across stored data
- **[Delete Memory](../api-usage-guide.md#deleting-memory-entries)**: Remove stored data by key

For more details, see the [API Usage Guide](../api-usage-guide.md). 