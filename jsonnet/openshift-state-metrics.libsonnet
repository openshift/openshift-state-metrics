local k = import 'ksonnet/ksonnet.beta.4/k.libsonnet';

{
  _config+:: {
    namespace: 'openshift-monitoring',

    openshiftStateMetrics+:: {
      collectors: '',  // empty string gets a default set
      scrapeInterval: '30s',
      scrapeTimeout: '30s',

      baseCPU: '100m',
      baseMemory: '150Mi',
      cpuPerNode: '2m',
      memoryPerNode: '30Mi',
    },

    versions+:: {
      openshiftStateMetrics: '4.2',
      kubeRbacProxy: '4.2',
    },

    imageRepos+:: {
      openshiftStateMetrics: 'quay.io/openshift/origin-openshift-state-metrics',
      kubeRbacProxy: 'quay.io/openshift/origin-kube-rbac-proxy',
    },

    tlsCipherSuites: [
      'TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256',  // required by h2: http://golang.org/cl/30721
      'TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256',  // required by h2: http://golang.org/cl/30721

      // 'TLS_RSA_WITH_RC4_128_SHA',            // insecure: https://access.redhat.com/security/cve/cve-2013-2566
      // 'TLS_RSA_WITH_3DES_EDE_CBC_SHA',       // insecure: https://access.redhat.com/articles/2548661
      // 'TLS_RSA_WITH_AES_128_CBC_SHA',        // disabled by h2
      // 'TLS_RSA_WITH_AES_256_CBC_SHA',        // disabled by h2
      'TLS_RSA_WITH_AES_128_CBC_SHA256',
      // 'TLS_RSA_WITH_AES_128_GCM_SHA256',     // disabled by h2
      // 'TLS_RSA_WITH_AES_256_GCM_SHA384',     // disabled by h2
      // 'TLS_ECDHE_ECDSA_WITH_RC4_128_SHA',    // insecure: https://access.redhat.com/security/cve/cve-2013-2566
      // 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA',// disabled by h2
      // 'TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA',// disabled by h2
      // 'TLS_ECDHE_RSA_WITH_RC4_128_SHA',      // insecure: https://access.redhat.com/security/cve/cve-2013-2566
      // 'TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA', // insecure: https://access.redhat.com/articles/2548661
      // 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA',  // disabled by h2
      // 'TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA',  // disabled by h2
      'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256',
      'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256',

      // disabled by h2 means: https://github.com/golang/net/blob/e514e69ffb8bc3c76a71ae40de0118d794855992/http2/ciphers.go

      // 'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384',   // TODO: Might not work with h2
      // 'TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384', // TODO: Might not work with h2
      // 'TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305',    // TODO: Might not work with h2
      // 'TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305',  // TODO: Might not work with h2
    ],

  },

  openshiftStateMetrics+:: {
    clusterRoleBinding:
      local clusterRoleBinding = k.rbac.v1.clusterRoleBinding;

      clusterRoleBinding.new() +
      clusterRoleBinding.mixin.metadata.withName('openshift-state-metrics') +
      clusterRoleBinding.mixin.roleRef.withApiGroup('rbac.authorization.k8s.io') +
      clusterRoleBinding.mixin.roleRef.withName('openshift-state-metrics') +
      clusterRoleBinding.mixin.roleRef.mixinInstance({ kind: 'ClusterRole' }) +
      clusterRoleBinding.withSubjects([{ kind: 'ServiceAccount', name: 'openshift-state-metrics', namespace: $._config.namespace }]),

    clusterRole:
      local clusterRole = k.rbac.v1.clusterRole;
      local rulesType = clusterRole.rulesType;


      local appsRule = rulesType.new() +
                       rulesType.withApiGroups(['apps.openshift.io']) +
                       rulesType.withResources([
                         'deploymentconfigs',
                       ]) +
                       rulesType.withVerbs(['list', 'watch']);

      local buildRule = rulesType.new() +
                        rulesType.withApiGroups(['build.openshift.io']) +
                        rulesType.withResources([
                          'buildconfigs',
                          'builds',
                        ]) +
                        rulesType.withVerbs(['list', 'watch']);

      local quotaRule = rulesType.new() +
                              rulesType.withApiGroups(['quota.openshift.io']) +
                              rulesType.withResources([
                                'clusterresourcequotas',
                              ]) +
                              rulesType.withVerbs(['list', 'watch']);


      local routeRule = rulesType.new() +
                         rulesType.withApiGroups(['route.openshift.io']) +
                         rulesType.withResources([
                           'routes',
                         ]) +
                         rulesType.withVerbs(['list', 'watch']);

      local authenticationRole = rulesType.new() +
                                 rulesType.withApiGroups(['authentication.k8s.io']) +
                                 rulesType.withResources([
                                   'tokenreviews',
                                 ]) +
                                 rulesType.withVerbs(['create']);

      local authorizationRole = rulesType.new() +
                                rulesType.withApiGroups(['authorization.k8s.io']) +
                                rulesType.withResources([
                                  'subjectaccessreviews',
                                ]) +
                                rulesType.withVerbs(['create']);

      local rules = [appsRule, buildRule, quotaRule, routeRule, authenticationRole, authorizationRole];

      clusterRole.new() +
      clusterRole.mixin.metadata.withName('openshift-state-metrics') +
      clusterRole.withRules(rules),
    deployment:
      local deployment = k.apps.v1.deployment;
      local container = k.apps.v1.deployment.mixin.spec.template.spec.containersType;
      local volume = k.apps.v1.deployment.mixin.spec.template.spec.volumesType;
      local containerPort = container.portsType;
      local containerVolumeMount = container.volumeMountsType;
      local podSelector = deployment.mixin.spec.template.spec.selectorType;

      local podLabels = { 'k8s-app': 'openshift-state-metrics' };
      local privateVolumeName = 'openshift-state-metrics-tls';
      local privateVolume = volume.fromSecret('openshift-state-metrics-tls', 'openshift-state-metrics-tls');
      local privateVolumeMount = containerVolumeMount.new(privateVolumeName, '/etc/tls/private');


      local proxyClusterMetrics =
        container.new('kube-rbac-proxy-main', $._config.imageRepos.kubeRbacProxy + ':' + $._config.versions.kubeRbacProxy) +
        container.withArgs([
          '--logtostderr',
          '--secure-listen-address=:8443',
          '--tls-cipher-suites=' + std.join(',', $._config.tlsCipherSuites),
          '--upstream=http://127.0.0.1:8081/',
          '--tls-cert-file=/etc/tls/private/tls.crt',
          '--tls-private-key-file=/etc/tls/private/tls.key',
        ]) +
        container.withPorts(containerPort.newNamed('https-main', 8443)) +
        container.mixin.resources.withRequests({ cpu: '10m', memory: '20Mi' }) +
        container.mixin.resources.withLimits({ cpu: '20m', memory: '40Mi' }) +
        container.withVolumeMounts([privateVolumeMount]);

      local proxySelfMetrics =
        container.new('kube-rbac-proxy-self', $._config.imageRepos.kubeRbacProxy + ':' + $._config.versions.kubeRbacProxy) +
        container.withArgs([
          '--logtostderr',
          '--secure-listen-address=:9443',
          '--tls-cipher-suites=' + std.join(',', $._config.tlsCipherSuites),
          '--upstream=http://127.0.0.1:8082/',
          '--tls-cert-file=/etc/tls/private/tls.crt',
          '--tls-private-key-file=/etc/tls/private/tls.key',
        ]) +
        container.withVolumeMounts([privateVolumeMount]) +
        container.withPorts(containerPort.newNamed('https-self', 9443)) +
        container.mixin.resources.withRequests({ cpu: '10m', memory: '20Mi' }) +
        container.mixin.resources.withLimits({ cpu: '20m', memory: '40Mi' });

      local openshiftStateMetrics =
        container.new('openshift-state-metrics', $._config.imageRepos.openshiftStateMetrics + ':' + $._config.versions.openshiftStateMetrics) +
        container.withArgs([
          '--host=127.0.0.1',
          '--port=8081',
          '--telemetry-host=127.0.0.1',
          '--telemetry-port=8082',
        ] + if $._config.openshiftStateMetrics.collectors != '' then ['--collectors=' + $._config.openshiftStateMetrics.collectors] else []) +
        container.mixin.resources.withRequests({ cpu: $._config.openshiftStateMetrics.baseCPU, memory: $._config.openshiftStateMetrics.baseMemory }) +
        container.mixin.resources.withLimits({ cpu: $._config.openshiftStateMetrics.baseCPU, memory: $._config.openshiftStateMetrics.baseMemory });

      local c = [proxyClusterMetrics, proxySelfMetrics, openshiftStateMetrics];

      deployment.new('openshift-state-metrics', 1, c, podLabels) +
      deployment.mixin.metadata.withNamespace($._config.namespace) +
      deployment.mixin.metadata.withLabels(podLabels) +
      deployment.mixin.spec.selector.withMatchLabels(podLabels) +
      deployment.mixin.spec.template.spec.withNodeSelector({ 'beta.kubernetes.io/os': 'linux' }) +
      deployment.mixin.spec.template.spec.withVolumes([privateVolume]) +
      deployment.mixin.spec.template.spec.withServiceAccountName('openshift-state-metrics') +
      deployment.mixin.spec.template.spec.withPriorityClassName('system-cluster-critical'),

    roleBinding:
      local roleBinding = k.rbac.v1.roleBinding;

      roleBinding.new() +
      roleBinding.mixin.metadata.withName('openshift-state-metrics') +
      roleBinding.mixin.metadata.withNamespace($._config.namespace) +
      roleBinding.mixin.roleRef.withApiGroup('rbac.authorization.k8s.io') +
      roleBinding.mixin.roleRef.withName('openshift-state-metrics') +
      roleBinding.mixin.roleRef.mixinInstance({ kind: 'Role' }) +
      roleBinding.withSubjects([{ kind: 'ServiceAccount', name: 'openshift-state-metrics' }]),

    role:
      local role = k.rbac.v1.role;
      local rulesType = role.rulesType;

      local coreRule = rulesType.new() +
                       rulesType.withApiGroups(['']) +
                       rulesType.withResources([
                         'pods',
                       ]) +
                       rulesType.withVerbs(['get']);

      local extensionsRule = rulesType.new() +
                             rulesType.withApiGroups(['extensions']) +
                             rulesType.withResources([
                               'deployments',
                             ]) +
                             rulesType.withVerbs(['get', 'update']) +
                             rulesType.withResourceNames(['openshift-state-metrics']);

      local appsRule = rulesType.new() +
                       rulesType.withApiGroups(['apps']) +
                       rulesType.withResources([
                         'deployments',
                       ]) +
                       rulesType.withVerbs(['get', 'update']) +
                       rulesType.withResourceNames(['openshift-state-metrics']);

      local rules = [coreRule, extensionsRule, appsRule];

      role.new() +
      role.mixin.metadata.withName('openshift-state-metrics') +
      role.mixin.metadata.withNamespace($._config.namespace) +
      role.withRules(rules),

    serviceAccount:
      local serviceAccount = k.core.v1.serviceAccount;

      serviceAccount.new('openshift-state-metrics') +
      serviceAccount.mixin.metadata.withNamespace($._config.namespace),

    service:
      local service = k.core.v1.service;
      local servicePort = k.core.v1.service.mixin.spec.portsType;

      local ksmServicePortMain = servicePort.newNamed('https-main', 8443, 'https-main');
      local ksmServicePortSelf = servicePort.newNamed('https-self', 9443, 'https-self');

      service.new('openshift-state-metrics', $.openshiftStateMetrics.deployment.spec.selector.matchLabels, [ksmServicePortMain, ksmServicePortSelf]) +
      service.mixin.metadata.withNamespace($._config.namespace) +
      service.mixin.metadata.withLabels({ 'k8s-app': 'openshift-state-metrics' }) +
      service.mixin.metadata.withAnnotations({
        'service.alpha.openshift.io/serving-cert-secret-name': 'openshift-state-metrics-tls',
      }),

    serviceMonitor:
      {
        apiVersion: 'monitoring.coreos.com/v1',
        kind: 'ServiceMonitor',
        metadata: {
          name: 'openshift-state-metrics',
          namespace: $._config.namespace,
          labels: {
            'k8s-app': 'openshift-state-metrics',
          },
        },
        spec: {
          jobLabel: 'k8s-app',
          selector: {
            matchLabels: {
              'k8s-app': 'openshift-state-metrics',
            },
          },
          endpoints: [
            {
              bearerTokenFile: '/var/run/secrets/kubernetes.io/serviceaccount/token',
              honorLabels: true,
              interval: '2m',
              scrapeTimeout: '2m',
              port: 'https-main',
              scheme: 'https',
              tlsConfig: {
                caFile: '/etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt',
                serverName: 'openshift-state-metrics.openshift-monitoring.svc',
              },
            },
            {
              bearerTokenFile: '/var/run/secrets/kubernetes.io/serviceaccount/token',
              interval: '2m',
              scrapeTimeout: '2m',
              port: 'https-self',
              scheme: 'https',
              tlsConfig: {
                caFile: '/etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt',
                serverName: 'openshift-state-metrics.openshift-monitoring.svc',
              },
            },
          ],
        },
      },
  },
}
