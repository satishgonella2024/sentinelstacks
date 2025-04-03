# Import Cycle Resolution

We've made significant progress in resolving import cycles in the codebase by implementing the following changes:

## 1. Common Types Package Structure

### Created the Following Files:
- `pkg/types/common.go`: Base package definition
- `pkg/types/memory.go`: Memory interfaces and types 
- `pkg/types/shim.go`: LLM provider interfaces
- `pkg/types/stack.go`: Stack interfaces and types
- `pkg/types/agent.go`: Agent runtime interfaces
- `pkg/types/registry.go`: Registry interfaces and types
- `pkg/types/runtime.go`: Runtime interfaces for agent execution

### Key Design Approach:
- All interface definitions are moved to the `pkg/types` package
- Implementation packages depend only on these common interfaces
- Direct dependencies between implementation packages are avoided

## 2. Updated Implementation Packages

### Memory Management:
- `internal/memory/factory.go`: Updated to use common types
- `internal/memory/local.go`: Refactored to implement common interfaces
- `internal/memory/plugin.go`: Added to support plugin-based extensions
- `internal/memory/types.go`: Added compatibility helpers for legacy code

### Stack Engine:
- `internal/stack/memory/state_manager.go`: Updated to use common types
- `internal/stack/memory/memory_manager.go`: Refactored to use common interfaces
- `internal/stack/runtime/factory.go`: Updated to implement common runtime interfaces
- `internal/stack/runtime/simple_runtime.go`: Simplified implementation of agent runtime

### Registry System:
- `internal/registry/auth/provider.go`: Implemented file-based token provider
- `internal/registry/package/sentinel_package.go`: Updated to use common registry types
- `cmd/sentinel/registry/registry.go`: Fixed duplicate function declarations

## 3. Verification

We've built the following packages successfully:
- `cmd/sentinel/registry`: The registry CLI command
- `internal/memory`: The memory management system
- `internal/stack/memory`: The stack memory management

## 4. Next Steps

Despite our progress, there may still be some import cycles to resolve:

1. Continue testing other packages for import cycles
2. Update any remaining packages to use common types
3. Add more unit tests to verify the refactored components
4. Complete the remaining registry functionality (search, tags)

## 5. Implementation Strategy

1. **Extract Common Interfaces**: Identify shared interfaces and move them to `pkg/types`
2. **Use Dependency Inversion**: Make implementations depend on interfaces, not other implementations
3. **Introduce Adapters**: For legacy code that can't be directly updated
4. **Create Integration Tests**: To verify everything works together correctly
