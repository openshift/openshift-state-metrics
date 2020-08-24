FROM registry.svc.ci.openshift.org/ocp/builder:rhel-8-golang-openshift-4.6 AS builder
WORKDIR /go/src/github.com/openshift/openshift-state-metrics
COPY . .
RUN make build

FROM  registry.svc.ci.openshift.org/ocp/4.6:base
LABEL io.k8s.display-name="openshift-state-metrics" \
      io.k8s.description="This is a component that exposes metrics about OpenShift objects." \
      io.openshift.tags="OpenShift" \
      maintainer="Haoran Wang <haowang@redhat.com>"

ARG FROM_DIRECTORY=/go/src/github.com/openshift/openshift-state-metrics
COPY --from=builder ${FROM_DIRECTORY}/openshift-state-metrics  /usr/bin/openshift-state-metrics

USER nobody
ENTRYPOINT ["/usr/bin/openshift-state-metrics"]
