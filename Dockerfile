# Use an official Ubuntu as a parent image
FROM ubuntu:24.04

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl wget bash && \
    rm -rf /var/lib/apt/lists/*

# Copy the rancher-project binary to /usr/local/bin
COPY rancher-projects /usr/local/bin

# Make the script executable
RUN chmod +x /usr/local/bin/rancher-projects

# Define entrypoint
ENTRYPOINT ["/usr/local/bin/rancher-projects"]

# Default command to run if no arguments are provided
CMD ["--help"]
