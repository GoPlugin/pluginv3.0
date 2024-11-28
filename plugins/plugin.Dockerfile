# Build image: Plugin binary
FROM golang:1.22-bullseye as buildgo
RUN go version
WORKDIR /plugin

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for plugin build
ARG COMMIT_SHA

COPY . .

RUN apt-get update && apt-get install -y jq

# Build the golang binaries
RUN make install-plugin

# Install medianpoc binary
RUN make install-medianpoc

# Install ocr3-capability binary
RUN make install-ocr3-capability

# Link LOOP Plugin source dirs with simple names
RUN go list -m -f "{{.Dir}}" github.com/goplugin/plugin-feeds | xargs -I % ln -s % /plugin-feeds
RUN go list -m -f "{{.Dir}}" github.com/goplugin/plugin-data-streams | xargs -I % ln -s % /plugin-data-streams
RUN go list -m -f "{{.Dir}}" github.com/goplugin/plugin-solana | xargs -I % ln -s % /plugin-solana
RUN mkdir /plugin-starknet
RUN go list -m -f "{{.Dir}}" github.com/goplugin/plugin-starknet/relayer | xargs -I % ln -s % /plugin-starknet/relayer

# Build image: Plugins
FROM golang:1.22-bullseye as buildplugins
RUN go version

WORKDIR /plugin-feeds
COPY --from=buildgo /plugin-feeds .
RUN go install ./cmd/plugin-feeds

WORKDIR /plugin-data-streams
COPY --from=buildgo /plugin-data-streams .
RUN go install ./mercury/cmd/plugin-mercury

WORKDIR /plugin-solana
COPY --from=buildgo /plugin-solana .
RUN go install ./pkg/solana/cmd/plugin-solana

WORKDIR /plugin-starknet/relayer
COPY --from=buildgo /plugin-starknet/relayer .
RUN go install ./pkg/plugin/cmd/plugin-starknet

# Final image: ubuntu with plugin binary
FROM ubuntu:20.04

ARG PLUGIN_USER=root
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && apt-get clean all

COPY --from=buildgo /go/bin/plugin /usr/local/bin/
COPY --from=buildgo /go/bin/plugin-medianpoc /usr/local/bin/
COPY --from=buildgo /go/bin/plugin-ocr3-capability /usr/local/bin/

COPY --from=buildplugins /go/bin/plugin-feeds /usr/local/bin/
ENV CL_MEDIAN_CMD plugin-feeds
COPY --from=buildplugins /go/bin/plugin-mercury /usr/local/bin/
ENV CL_MERCURY_CMD plugin-mercury
COPY --from=buildplugins /go/bin/plugin-solana /usr/local/bin/
ENV CL_SOLANA_CMD plugin-solana
COPY --from=buildplugins /go/bin/plugin-starknet /usr/local/bin/
ENV CL_STARKNET_CMD plugin-starknet

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
