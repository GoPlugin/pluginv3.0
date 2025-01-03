#!/bin/bash
# ASCII Spinner
spin_animation() {
#       local frames=(" PLI " "  LIP" "I  LP" "IP  L" "LIP " " PL I" "PLI ")
local frames=(
"   PLI🚂💨      "
"    PLI🚂💨     "
"     PLI🚂💨    "
"      PLI🚂💨   "
"       PLI🚂💨  "
"        PLI🚂💨 "
"         PLI🚂💨"
"        PLI🚂💨 "
"       PLI🚂💨  "
"      PLI🚂💨   "
"     PLI🚂💨    "
"    PLI🚂💨     "
"   PLI🚂💨      "
)
  while :; do
    for frame in "${frames[@]}"; do
      echo -ne "\r[$frame] Executing command... "
      sleep 0.2
    done
  done
}

echo "Starting Plugin build process... Total steps: 6"

# Step 1: Setup Go env variables
echo "Step 1/6: Appending environment variables to ~/.bashrc..."
# Run animation in the background
spin_animation &
spinner_pid=$!

# Define your 
ENV_VARS=(
    "export GONOPROXY='github.com/goplugin'"
    "export GONOSUMDB='github.com/goplugin'"
    "export GO111MODULE='on'"
    "export GOPRIVATE='github.com/goplugin'"
    "export GOPROXY='direct'"
    "export CGO_ENABLED='1'"
)
# Append variables to ~/.bashrc if they are not already present
for VAR in "${ENV_VARS[@]}"; do
    if ! grep -qxF "$VAR" ~/.bashrc; then
        echo "$VAR" >> ~/.bashrc
        echo "Added: $VAR"
    else
        echo "Already exists: $VAR"
    fi
done
# Reload ~/.bashrc to apply changes
source ~/.bashrc

# Kill spinner after command finishes
kill $spinner_pid
wait $spinner_pid 2>/dev/null

echo -e "\nCommand completed!"

# Step 2: Install contracts dependencies
echo "Step 2/6: Installing contracts dependencies..."
spin_animation &
spinner_pid=$!
cd contracts && pnpm i && cd ../ 
# Kill spinner after command finishes
kill $spinner_pid
wait $spinner_pid 2>/dev/null

# Step 3: Download Go modules
echo "Step 3/6: Downloading Go modules..."
spin_animation &
spinner_pid=$!
go mod download 
# Kill spinner after command finishes
kill $spinner_pid
wait $spinner_pid 2>/dev/null

# Step 4: Install Operator UI
echo "Step 4/6: Installing Operator UI..."
spin_animation &
spinner_pid=$!
go run operator_ui/install.go . 
# Kill spinner after command finishes
kill $spinner_pid
wait $spinner_pid 2>/dev/null

# Step 5: Build the Plugin binary
echo "Step 5/6: Building Plugin binary..."
spin_animation &
spinner_pid=$!
go build -ldflags "-X github.com/goplugin/pluginv3.0/v2/core/static.Version=2.4.0 -X github.com/goplugin/pluginv3.0/v2/core/static.Sha=b1245c440825ebbb342c9bfa1b0cfa9da54dae53" -o plugin 
# Kill spinner after command finishes
kill $spinner_pid
wait $spinner_pid 2>/dev/null

# Step 6: Move binary to Go bin folder
echo "Step 6/6: Moving Plugin binary to Go bin folder..."
spin_animation &
spinner_pid=$!
GOP=$(grep 'GOPATH=' ~/.profile |cut -d '=' -f 2)
echo "$GOP";
mv plugin $GOP/bin/;
# Kill spinner after command finishes
kill $spinner_pid
wait $spinner_pid 2>/dev/null

echo "Plugin build process completed successfully!"
