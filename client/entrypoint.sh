#!/bin/bash
set -e

# Create new Svelte project if package.json doesn't exist
if [ ! -f "/app/package.json" ]; then
    echo "Creating new Svelte project..."
    # Create temporary directory
    mkdir -p /tmp/svelte-temp
    cd /tmp/svelte-temp
    
    # Create new project
    create-vite . --template svelte --skip-git
    
    # Copy files to app directory
    cp -a . /app/
    
    # Clean up
    cd /app
    rm -rf /tmp/svelte-temp
fi

# Install dependencies
echo "Installing dependencies..."
npm install

# Start development server
echo "Starting development server..."
npm run dev -- --host