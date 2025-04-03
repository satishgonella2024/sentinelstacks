# Deploying to GitHub Pages

This document explains how to deploy the SentinelStacks documentation, including the API documentation, to GitHub Pages.

## Automatic Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The deployment is handled by a GitHub Actions workflow defined in `.github/workflows/deploy-docs.yml`.

The workflow:
1. Builds the API documentation using the `scripts/generate_api_docs.sh` script
2. Builds the MkDocs site
3. Deploys the site to GitHub Pages

## Manual Deployment

If you need to deploy the documentation manually, follow these steps:

1. Generate the API documentation:
   ```bash
   ./scripts/generate_api_docs.sh
   ```

2. Build the MkDocs site:
   ```bash
   mkdocs build
   ```

3. Deploy the site to GitHub Pages:
   ```bash
   mkdocs gh-deploy --force
   ```

## Structure on GitHub Pages

The documentation is deployed with the following structure:

- `/`: The main documentation site
- `/api/`: The API documentation page with Swagger UI
- `/api-reference.yaml`: The OpenAPI specification
- `/api-usage-guide/`: The API usage guide

## Troubleshooting

If you encounter issues with the deployment:

1. Check that GitHub Pages is enabled for the repository and set to deploy from the `gh-pages` branch
2. Verify that the GitHub Actions workflow has the necessary permissions
3. Ensure that the API documentation is correctly generated
4. Check the GitHub Actions logs for errors
5. Test locally by running `mkdocs serve` before pushing changes

## Customizing the Deployment

To customize the deployment:

1. Edit the `.github/workflows/deploy-docs.yml` file to modify the GitHub Actions workflow
2. Update the `mkdocs.yml` file to change the site configuration
3. Modify the `scripts/sync_api_docs.sh` script to adjust how the API documentation is synced with MkDocs 