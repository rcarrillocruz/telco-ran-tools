FROM mirror.gcr.io/library/golang:1.18 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

COPY hack/ hack/

# Build
RUN hack/build-binaries.sh

# Create image
FROM registry.ci.openshift.org/ocp/4.12:tools
RUN yum install -y skopeo gdisk && \
    yum clean -y all
RUN curl -sL https://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/oc-mirror.tar.gz \
    | tar xvzf - -C /usr/local/bin && \
    chmod +x /usr/local/bin/oc-mirror

COPY run.sh /run.sh
COPY help.sh /help.sh

COPY --from=builder /workspace/_output/ /usr/local/bin/

# no tool is more important than others
ENTRYPOINT ["/run.sh"]
CMD ["/help.sh"]
