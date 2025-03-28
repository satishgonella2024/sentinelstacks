package tools

import (
	"fmt"
	"strings"
)

// TerraformTool provides Terraform infrastructure as code generation capabilities
type TerraformTool struct{}

// ID returns the unique identifier for the Terraform tool
func (t *TerraformTool) ID() string {
	return "terraform"
}

// Name returns a user-friendly name
func (t *TerraformTool) Name() string {
	return "Terraform Generator"
}

// Description returns a detailed description
func (t *TerraformTool) Description() string {
	return "Generates Terraform infrastructure as code based on resource type and configuration"
}

// Version returns the semantic version
func (t *TerraformTool) Version() string {
	return "0.1.0"
}

// ParameterSchema returns the JSON schema for parameters
func (t *TerraformTool) ParameterSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"provider": map[string]interface{}{
				"type": "string",
				"enum": []string{"aws", "azure", "gcp", "digitalocean", "kubernetes"},
				"description": "Infrastructure provider",
			},
			"resource_type": map[string]interface{}{
				"type": "string",
				"description": "Type of resource to generate (e.g., aws_instance, azurerm_virtual_machine)",
			},
			"config": map[string]interface{}{
				"type": "object",
				"description": "Resource configuration parameters",
			},
			"format": map[string]interface{}{
				"type": "string",
				"enum": []string{"tf", "json", "hcl"},
				"default": "tf",
				"description": "Output format (tf for standard Terraform, json for JSON, hcl for raw HCL)",
			},
		},
		"required": []string{"provider", "resource_type"},
	}
}

// Execute runs the Terraform tool with the provided parameters
func (t *TerraformTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Extract parameters
	provider, ok := params["provider"].(string)
	if !ok {
		return nil, fmt.Errorf("provider parameter is required")
	}

	resourceType, ok := params["resource_type"].(string)
	if !ok {
		return nil, fmt.Errorf("resource_type parameter is required")
	}

	// Default format is 'tf'
	format := "tf"
	if formatParam, ok := params["format"].(string); ok {
		format = formatParam
	}

	// Resource configuration
	config := make(map[string]interface{})
	if configParam, ok := params["config"].(map[string]interface{}); ok {
		config = configParam
	}

	// Generate the code based on provider and resource type
	switch provider {
	case "aws":
		return t.generateAWSResource(resourceType, config, format)
	case "azure":
		return t.generateAzureResource(resourceType, config, format)
	case "gcp":
		return t.generateGCPResource(resourceType, config, format)
	case "digitalocean":
		return t.generateDigitalOceanResource(resourceType, config, format)
	case "kubernetes":
		return t.generateKubernetesResource(resourceType, config, format)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// generateAWSResource generates AWS Terraform code
func (t *TerraformTool) generateAWSResource(resourceType string, config map[string]interface{}, format string) (string, error) {
	var result strings.Builder

	// Provider block
	result.WriteString("provider \"aws\" {\n")
	result.WriteString("  region = \"us-west-2\" # Change to your preferred region\n")
	result.WriteString("}\n\n")

	// Resource block
	result.WriteString(fmt.Sprintf("resource \"%s\" \"this\" {\n", resourceType))

	// Add config parameters
	for key, value := range config {
		// Format the value based on its type
		valueStr := formatValue(value)
		result.WriteString(fmt.Sprintf("  %s = %s\n", key, valueStr))
	}

	// Add default parameters based on resource type
	switch resourceType {
	case "aws_instance":
		if _, exists := config["ami"]; !exists {
			result.WriteString("  ami           = \"ami-0c55b159cbfafe1f0\" # Amazon Linux 2 AMI\n")
		}
		if _, exists := config["instance_type"]; !exists {
			result.WriteString("  instance_type = \"t2.micro\"\n")
		}
		if _, exists := config["tags"]; !exists {
			result.WriteString("  tags = {\n")
			result.WriteString("    Name = \"example-instance\"\n")
			result.WriteString("  }\n")
		}
	case "aws_s3_bucket":
		if _, exists := config["bucket"]; !exists {
			result.WriteString("  bucket = \"my-tf-test-bucket\"\n")
		}
		if _, exists := config["tags"]; !exists {
			result.WriteString("  tags = {\n")
			result.WriteString("    Name        = \"My bucket\"\n")
			result.WriteString("    Environment = \"Dev\"\n")
			result.WriteString("  }\n")
		}
	}

	result.WriteString("}\n")

	// Convert to the requested format
	if format == "json" {
		// This is a simplified conversion - in a real tool you'd use HCL parser
		return "{\n  \"resource\": {\n    \"" + resourceType + "\": {\n      \"this\": {\n        // JSON format not fully implemented\n      }\n    }\n  }\n}", nil
	}

	return result.String(), nil
}

// generateAzureResource generates Azure Terraform code
func (t *TerraformTool) generateAzureResource(resourceType string, config map[string]interface{}, format string) (string, error) {
	var result strings.Builder

	// Provider block
	result.WriteString("provider \"azurerm\" {\n")
	result.WriteString("  features {}\n")
	result.WriteString("}\n\n")

	// Resource group (commonly needed)
	if !strings.Contains(resourceType, "resource_group") {
		result.WriteString("resource \"azurerm_resource_group\" \"example\" {\n")
		result.WriteString("  name     = \"example-resources\"\n")
		result.WriteString("  location = \"West Europe\"\n")
		result.WriteString("}\n\n")
	}

	// Resource block
	result.WriteString(fmt.Sprintf("resource \"%s\" \"this\" {\n", resourceType))

	// Add config parameters
	for key, value := range config {
		valueStr := formatValue(value)
		result.WriteString(fmt.Sprintf("  %s = %s\n", key, valueStr))
	}

	// Add default parameters based on resource type
	switch resourceType {
	case "azurerm_virtual_machine":
		if _, exists := config["name"]; !exists {
			result.WriteString("  name                  = \"example-vm\"\n")
		}
		if _, exists := config["location"]; !exists {
			result.WriteString("  location              = azurerm_resource_group.example.location\n")
		}
		if _, exists := config["resource_group_name"]; !exists {
			result.WriteString("  resource_group_name   = azurerm_resource_group.example.name\n")
		}
		// Add more default parameters as needed
	}

	result.WriteString("}\n")

	// Convert to the requested format
	if format == "json" {
		// This is a simplified conversion
		return "{\n  \"resource\": {\n    \"" + resourceType + "\": {\n      \"this\": {\n        // JSON format not fully implemented\n      }\n    }\n  }\n}", nil
	}

	return result.String(), nil
}

// generateGCPResource generates GCP Terraform code
func (t *TerraformTool) generateGCPResource(resourceType string, config map[string]interface{}, format string) (string, error) {
	var result strings.Builder

	// Provider block
	result.WriteString("provider \"google\" {\n")
	result.WriteString("  project = \"your-project-id\"\n")
	result.WriteString("  region  = \"us-central1\"\n")
	result.WriteString("  zone    = \"us-central1-c\"\n")
	result.WriteString("}\n\n")

	// Resource block
	result.WriteString(fmt.Sprintf("resource \"%s\" \"this\" {\n", resourceType))

	// Add config parameters
	for key, value := range config {
		valueStr := formatValue(value)
		result.WriteString(fmt.Sprintf("  %s = %s\n", key, valueStr))
	}

	// Add default parameters based on resource type
	switch resourceType {
	case "google_compute_instance":
		if _, exists := config["name"]; !exists {
			result.WriteString("  name         = \"test-instance\"\n")
		}
		if _, exists := config["machine_type"]; !exists {
			result.WriteString("  machine_type = \"e2-medium\"\n")
		}
		if _, exists := config["boot_disk"]; !exists {
			result.WriteString("  boot_disk {\n")
			result.WriteString("    initialize_params {\n")
			result.WriteString("      image = \"debian-cloud/debian-11\"\n")
			result.WriteString("    }\n")
			result.WriteString("  }\n")
		}
		if _, exists := config["network_interface"]; !exists {
			result.WriteString("  network_interface {\n")
			result.WriteString("    network = \"default\"\n")
			result.WriteString("    access_config {\n")
			result.WriteString("      // Ephemeral public IP\n")
			result.WriteString("    }\n")
			result.WriteString("  }\n")
		}
	}

	result.WriteString("}\n")

	return result.String(), nil
}

// generateDigitalOceanResource generates DigitalOcean Terraform code
func (t *TerraformTool) generateDigitalOceanResource(resourceType string, config map[string]interface{}, format string) (string, error) {
	var result strings.Builder

	// Provider block
	result.WriteString("provider \"digitalocean\" {\n")
	result.WriteString("  # Set this to your DigitalOcean API token\n")
	result.WriteString("  # token = var.do_token\n")
	result.WriteString("}\n\n")

	// Resource block
	result.WriteString(fmt.Sprintf("resource \"%s\" \"this\" {\n", resourceType))

	// Add config parameters
	for key, value := range config {
		valueStr := formatValue(value)
		result.WriteString(fmt.Sprintf("  %s = %s\n", key, valueStr))
	}

	// Add default parameters based on resource type
	switch resourceType {
	case "digitalocean_droplet":
		if _, exists := config["name"]; !exists {
			result.WriteString("  name   = \"test-droplet\"\n")
		}
		if _, exists := config["size"]; !exists {
			result.WriteString("  size   = \"s-1vcpu-1gb\"\n")
		}
		if _, exists := config["image"]; !exists {
			result.WriteString("  image  = \"ubuntu-20-04-x64\"\n")
		}
		if _, exists := config["region"]; !exists {
			result.WriteString("  region = \"nyc3\"\n")
		}
	}

	result.WriteString("}\n")

	return result.String(), nil
}

// generateKubernetesResource generates Kubernetes Terraform code
func (t *TerraformTool) generateKubernetesResource(resourceType string, config map[string]interface{}, format string) (string, error) {
	var result strings.Builder

	// Provider block
	result.WriteString("provider \"kubernetes\" {\n")
	result.WriteString("  config_path    = \"~/.kube/config\"\n")
	result.WriteString("  config_context = \"my-context\"\n")
	result.WriteString("}\n\n")

	// Resource block
	result.WriteString(fmt.Sprintf("resource \"%s\" \"this\" {\n", resourceType))

	// Add config parameters
	for key, value := range config {
		valueStr := formatValue(value)
		result.WriteString(fmt.Sprintf("  %s = %s\n", key, valueStr))
	}

	// Add default parameters based on resource type
	switch resourceType {
	case "kubernetes_deployment":
		if _, exists := config["metadata"]; !exists {
			result.WriteString("  metadata {\n")
			result.WriteString("    name = \"nginx\"\n")
			result.WriteString("    labels = {\n")
			result.WriteString("      app = \"nginx\"\n")
			result.WriteString("    }\n")
			result.WriteString("  }\n\n")
		}
		if _, exists := config["spec"]; !exists {
			result.WriteString("  spec {\n")
			result.WriteString("    replicas = 2\n\n")
			result.WriteString("    selector {\n")
			result.WriteString("      match_labels = {\n")
			result.WriteString("        app = \"nginx\"\n")
			result.WriteString("      }\n")
			result.WriteString("    }\n\n")
			result.WriteString("    template {\n")
			result.WriteString("      metadata {\n")
			result.WriteString("        labels = {\n")
			result.WriteString("          app = \"nginx\"\n")
			result.WriteString("        }\n")
			result.WriteString("      }\n\n")
			result.WriteString("      spec {\n")
			result.WriteString("        container {\n")
			result.WriteString("          image = \"nginx:1.21.6\"\n")
			result.WriteString("          name  = \"nginx\"\n\n")
			result.WriteString("          resources {\n")
			result.WriteString("            limits = {\n")
			result.WriteString("              cpu    = \"0.5\"\n")
			result.WriteString("              memory = \"512Mi\"\n")
			result.WriteString("            }\n")
			result.WriteString("            requests = {\n")
			result.WriteString("              cpu    = \"250m\"\n")
			result.WriteString("              memory = \"50Mi\"\n")
			result.WriteString("            }\n")
			result.WriteString("          }\n")
			result.WriteString("        }\n")
			result.WriteString("      }\n")
			result.WriteString("    }\n")
			result.WriteString("  }\n")
		}
	}

	result.WriteString("}\n")

	return result.String(), nil
}

// formatValue formats a value based on its type
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case map[string]interface{}:
		var builder strings.Builder
		builder.WriteString("{\n")
		for k, val := range v {
			builder.WriteString(fmt.Sprintf("    %s = %s\n", k, formatValue(val)))
		}
		builder.WriteString("  }")
		return builder.String()
	case []interface{}:
		var builder strings.Builder
		builder.WriteString("[\n")
		for _, val := range v {
			builder.WriteString(fmt.Sprintf("    %s,\n", formatValue(val)))
		}
		builder.WriteString("  ]")
		return builder.String()
	default:
		return fmt.Sprintf("\"%v\"", v)
	}
}
