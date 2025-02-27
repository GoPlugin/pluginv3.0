#!/bin/bash

# ASCII Spinner with exit control
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

    local pipe="/tmp/spinner_fifo_$$"  # Unique name using PID
    mkfifo "$pipe"
    exec 3<> "$pipe"

    trap "rm -f $pipe; exec 3>&-" EXIT  # Cleanup on exit

    while :; do
        for frame in "${frames[@]}"; do
            echo -ne "\r[$frame] $SPINNER_MSG"
            sleep 0.2
            if read -t 0.1 <&3; then break 2; fi
        done
    done
    exec 3>&-
    rm -f "$pipe"
}


# Function to run a command with spinner
run_with_spinner() {
    SPINNER_MSG="$1"
    shift
    spin_animation &  # Start spinner in the background
    local spinner_pid=$!
    
    "$@"  # Execute the command
    local cmd_status=$?

    if kill -0 $spinner_pid 2>/dev/null; then
        kill $spinner_pid 2>/dev/null
        wait $spinner_pid 2>/dev/null
    fi

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
echo "Starting Plugin build process... Total steps: 5"

# Step 1: Ensure required tools are installed
if ! command -v go &>/dev/null; then
    echo "Error: Go is not installed. Please install it and retry."
    exit 1
fi

if ! command -v pnpm &>/dev/null; then
    echo "Error: pnpm is not installed. Please install it and retry."
    exit 1
fi

# Step 2: Setup Go environment variables
echo "Step 1/5: Setting up Go environment variables..."
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

# Step 3: Install contracts dependencies
run_with_spinner "Installing contracts dependencies..." bash -c "cd contracts && pnpm i && cd ../"
check_status $? "Contracts dependencies installation"

# Step 4: Install Operator UI using vendor
run_with_spinner "Installing Operator UI..." go run -mod=vendor operator_ui/install.go .
check_status $? "Operator UI installation"

# Step 5: Build Plugin binary using vendor
run_with_spinner "Building Plugin binary..." \
    go build -mod=vendor -ldflags "-X github.com/goplugin/pluginv3.0/v2/core/static.Version=2.4.0 -X github.com/goplugin/pluginv3.0/v2/core/static.Sha=b1245c440825ebbb342c9bfa1b0cfa9da54dae53" -o plugin
check_status $? "Plugin binary build"

# Step 6: Move binary to Go bin folder
GOP=$(go env GOPATH)
if [[ -z "$GOP" ]]; then
    echo "Error: GOPATH is not set. Exiting."
    exit 1
fi
run_with_spinner "Moving Plugin binary to Go bin folder..." mv plugin "$GOP/bin/"
check_status $? "Binary move"

echo "Plugin build process completed successfully!"

