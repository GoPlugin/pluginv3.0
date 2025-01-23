#!/bin/bash

# ASCII Spinner
spin_animation() {
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
            echo -ne "\r[$frame] $SPINNER_MSG"
            sleep 0.2
        done
    done
}

# Function to run a command with spinner
run_with_spinner() {
    SPINNER_MSG="$1"
    shift
    spin_animation &  # Start spinner in the background
    local spinner_pid=$!
    trap "kill $spinner_pid 2>/dev/null" EXIT  # Ensure spinner stops on exit
    "$@"  # Execute the command
    local cmd_status=$?
    kill $spinner_pid 2>/dev/null  # Stop spinner
    trap - EXIT
    echo -ne "\r"  # Clear spinner line
    return $cmd_status
}

# Check command status and exit on failure
check_status() {
    if [[ $1 -ne 0 ]]; then
        echo "Error: $2 failed. Exiting."
        exit 1
    fi
}

# Main script logic
echo "Starting Plugin build process... Total steps: 6"

# Step 1: Setup Go environment variables
echo "Step 1/6: Setting up Go environment variables..."
export GONOPROXY='github.com/goplugin'
export GONOSUMDB='github.com/goplugin'
export GO111MODULE='on'
export GOPRIVATE='github.com/goplugin'
export GOPROXY='direct'
export CGO_ENABLED='1'

# Optional: Append these variables to ~/.bashrc for future sessions
ENV_VARS=(
    "export GONOPROXY='github.com/goplugin'"
    "export GONOSUMDB='github.com/goplugin'"
    "export GO111MODULE='on'"
    "export GOPRIVATE='github.com/goplugin'"
    "export GOPROXY='direct'"
    "export CGO_ENABLED='1'"
)
for VAR in "${ENV_VARS[@]}"; do
    if ! grep -qxF "$VAR" ~/.bashrc; then
        echo "$VAR" >> ~/.bashrc
        echo "Added: $VAR"
    else
        echo "Already exists: $VAR"
    fi
done

# Step 2: Install contracts dependencies
run_with_spinner "Installing contracts dependencies..." bash -c "cd contracts && pnpm i && cd ../"
check_status $? "Contracts dependencies installation"

# Step 3: Build Go modules
run_with_spinner "Building Go modules..." go mod tidy
check_status $? "Go modules build"

# Step 4: Install Operator UI
run_with_spinner "Installing Operator UI..." go run operator_ui/install.go .
check_status $? "Operator UI installation"

# Step 5: Build Plugin binary
run_with_spinner "Building Plugin binary..." \
    go build -ldflags "-X github.com/goplugin/pluginv3.0/v2/core/static.Version=2.4.0 -X github.com/goplugin/pluginv3.0/v2/core/static.Sha=b1245c440825ebbb342c9bfa1b0cfa9da54dae53" -o plugin
check_status $? "Plugin binary build"

# Step 6: Move binary to Go bin folder
GOP=$(grep 'GOPATH=' ~/.profile | cut -d '=' -f 2)
if [[ -z "$GOP" ]]; then
    echo "Error: GOPATH is not set in ~/.profile. Exiting."
    exit 1
fi
run_with_spinner "Moving Plugin binary to Go bin folder..." mv plugin "$GOP/bin/"
check_status $? "Binary move"

echo "Plugin build process completed successfully!"
