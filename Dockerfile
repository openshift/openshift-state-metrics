FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.22-openshift-4.17 AS builder
WORKDIR /go/src/github.com/openshift/openshift-state-metrics
COPY . .
RUN make build

FROM registry.ci.openshift.org/ocp/4.17:base-rhel9
LABEL io.k8s.display-name="openshift-state-metrics" \
      io.k8s.description="This is a component that exposes metrics about OpenShift objects." \
      io.openshift.tags="OpenShift" \
      maintainer="OpenShift Monitoring Team <team-monitoring@redhat.com>"

ARG FROM_DIRECTORY=/go/src/github.com/openshift/openshift-state-metrics
COPY --from=builder ${FROM_DIRECTORY}/openshift-state-metrics  /usr/bin/openshift-state-metrics

USER nobody
ENTRYPOINT ["/usr/bin/openshift-state-metrics"]
