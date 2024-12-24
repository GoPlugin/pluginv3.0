#!/bin/bash

echo "Starting Plugin build process... Total steps: 6"

# Step 1: Setup Go env variables
echo "Step 1/6: Installing contracts dependencies..."
export GONOPROXY='github.com/goplugin' || {echo "Failed to setup GONOPROXY"; exit 1;}
export GONOSUMDB='github.com/goplugin' || {echo "Failed to setup GONOSUMDB"; exit 1;}
export GO111MODULE='on' || {echo "Failed to setup GO111MODULE"; exit 1;}
export GOPRIVATE='github.com/goplugin' || {echo "Failed to setup GOPRIVATE"; exit 1;}
export GOPROXY='direct' || {echo "Failed to setup GOPROXY"; exit 1;}
export CGO_ENABLED='1' || {echo "Failed to setup CGO_ENABLED"; exit 1;}
source ~/.bashrc || {echo "Failed to run bashrc"; exit 1;}

# Step 2: Install contracts dependencies
echo "Step 2/6: Installing contracts dependencies..."
cd contracts && pnpm i && cd ../ || { echo "Failed to install contracts dependencies"; exit 1; }

# Step 3: Download Go modules
echo "Step 3/6: Downloading Go modules..."
go mod download || { echo "Failed to download Go modules"; exit 1; }

# Step 4: Install Operator UI
echo "Step 4/6: Installing Operator UI..."
go run operator_ui/install.go . || { echo "Failed to install Operator UI"; exit 1; }

# Step 5: Build the Plugin binary
echo "Step 5/6: Building Plugin binary..."
go build -ldflags "-X github.com/goplugin/pluginv3.0/v2/core/static.Version=2.4.0 -X github.com/goplugin/pluginv3.0/v2/core/static.Sha=b1245c440825ebbb342c9bfa1b0cfa9da54dae53" -o plugin || { echo "Failed to build Plugin binary"; exit 1; }

# Step 6: Move binary to Go bin folder
echo "Step 6/6: Moving Plugin binary to Go bin folder..."
mv plugin $GOPATH/bin/ || { echo "Failed to move Plugin binary"; exit 1; }

echo "Plugin build process completed successfully!"

