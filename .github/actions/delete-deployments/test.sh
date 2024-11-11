#!/bin/sh
export NUM_OF_PAGES=all
export ENVIRONMENT=integration
export DRY_RUN=false
export REPOSITORY=goplugin/pluginv3.0
export REF=fix/golint
export GITHUB_ACTION=true

pnpm start
