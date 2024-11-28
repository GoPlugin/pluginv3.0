# Build image: Plugin binary
FROM golang:1.21-bullseye AS buildgo
RUN go version
WORKDIR /plugin

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for plugin build
ARG COMMIT_SHA

COPY . .

# Build the golang binary
RUN make install-plugin

# Link LOOP Plugin source dirs with simple names
RUN go list -m -f "{{.Dir}}" github.com/goplugin/plugin-feeds | xargs -I % ln -s % /plugin-feeds
RUN go list -m -f "{{.Dir}}" github.com/goplugin/plugin-solana | xargs -I % ln -s % /plugin-solana

# Build image: Plugins
FROM golang:1.21-bullseye AS buildplugins
RUN go version

WORKDIR /plugin-feeds
COPY --from=buildgo /plugin-feeds .
RUN go install ./cmd/plugin-feeds

WORKDIR /plugin-solana
COPY --from=buildgo /plugin-solana .
RUN go install ./pkg/solana/cmd/plugin-solana

# Final image: ubuntu with plugin binary
FROM --platform=linux/amd64 golang:1.21-bullseye

ARG PLUGIN_USER=plugin
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-15 \
  && apt-get clean all

COPY --from=buildgo /go/bin/plugin /usr/local/bin/

# Install (but don't enable) LOOP Plugins
COPY --from=buildplugins /go/bin/plugin-feeds /usr/local/bin/
COPY --from=buildplugins /go/bin/plugin-solana /usr/local/bin/

# Dependency of CosmWasm/wasmd
COPY --from=buildgo /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.*.so /usr/lib/
RUN chmod 755 /usr/lib/libwasmvm.*.so

RUN if [ ${PLUGIN_USER} != root ]; then \
  useradd --uid 14933 --create-home ${PLUGIN_USER}; \
  fi
USER ${PLUGIN_USER}
WORKDIR /home/${PLUGIN_USER}
# explicit set the cache dir. needed so both root and non-root user has an explicit location
ENV XDG_CACHE_HOME /home/${PLUGIN_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

EXPOSE 6688
ENTRYPOINT ["plugin"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]
