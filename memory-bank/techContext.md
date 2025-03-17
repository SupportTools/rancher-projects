# Tech Context: Rancher Projects

## Development Environment
- **Language**: Go (v1.23.0)
- **Build System**: Standard Go toolchain
- **Version Control**: Git
- **Dependencies**:
  - github.com/sirupsen/logrus v1.9.3 (logging)
  - github.com/stretchr/testify v1.10.0 (testing)

## Deployment Models

### Shell Script Deployment
- Direct installation via curl or wget
- Script placed in system path (/usr/local/bin)
- Executable permissions set for direct command-line execution

### Docker Container
- Containerized deployment via Dockerfile
- Enables consistent execution environment
- Facilitates integration with container orchestration

### GitHub Action
- Custom action for GitHub workflows
- Enables integration into CI/CD pipelines
- Provides standardized interface for GitHub-based automation

## Build Process
- Standard Go build process
- Version information injected via linker flags
- Build metadata captured at compile time:
  - Version number
  - Git commit hash
  - Build timestamp

## CI/CD Pipeline
- **Platform**: Drone CI (drone.support.tools)
- **Trigger**: Commits to main branch
- **Process**:
  - Code checkout
  - Go dependency validation
  - Tests execution
  - Build artifact generation
  - Docker image creation and publication
  - GitHub release creation

## Runtime Requirements
- Network access to Rancher API
- Valid Rancher API credentials
- Sufficient permissions for cluster/project/namespace operations

## Operational Considerations
- **Security**: Requires API secrets that should be protected
- **Performance**: Minimal resource requirements
- **Scalability**: Can operate on multiple clusters in sequence
- **Monitoring**: Basic stdout/stderr logging
- **Networking**: Requires outbound connectivity to Rancher API

## Deployment Artifacts
1. **Shell Script**: Direct execution in any shell environment
2. **Go Binary**: Compiled executable for specific platforms
3. **Docker Image**: Containerized deployment
4. **GitHub Action**: Integrated component for GitHub workflows
