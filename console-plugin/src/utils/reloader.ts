import type { K8sResourceCommon } from '@openshift-console/dynamic-plugin-sdk';

// Mirrors pkg/config.AnnotationConfig. Reloader lets every key be overridden by a flag, so the
// running values are read from the reloader-meta-info ConfigMap whenever it is visible.
export type AnnotationConfig = {
  auto: string;
  configmapAuto: string;
  secretAuto: string;
  secretProviderClassAuto: string;
  configmapReload: string;
  secretReload: string;
  secretProviderClassReload: string;
  configmapExclude: string;
  secretExclude: string;
  secretProviderClassExclude: string;
  ignore: string;
  search: string;
  match: string;
  rolloutStrategy: string;
  pausePeriod: string;
  pausedAt: string;
  lastReloadedFrom: string;
};

export const DEFAULT_ANNOTATIONS: AnnotationConfig = {
  auto: 'reloader.stakater.com/auto',
  configmapAuto: 'configmap.reloader.stakater.com/auto',
  secretAuto: 'secret.reloader.stakater.com/auto',
  secretProviderClassAuto: 'secretproviderclass.reloader.stakater.com/auto',
  configmapReload: 'configmap.reloader.stakater.com/reload',
  secretReload: 'secret.reloader.stakater.com/reload',
  secretProviderClassReload: 'secretproviderclass.reloader.stakater.com/reload',
  configmapExclude: 'configmaps.exclude.reloader.stakater.com/reload',
  secretExclude: 'secrets.exclude.reloader.stakater.com/reload',
  secretProviderClassExclude: 'secretproviderclasses.exclude.reloader.stakater.com/reload',
  ignore: 'reloader.stakater.com/ignore',
  search: 'reloader.stakater.com/search',
  match: 'reloader.stakater.com/match',
  rolloutStrategy: 'reloader.stakater.com/rollout-strategy',
  pausePeriod: 'deployment.reloader.stakater.com/pause-period',
  pausedAt: 'deployment.reloader.stakater.com/paused-at',
  lastReloadedFrom: 'reloader.stakater.com/last-reloaded-from',
};

export const META_INFO_LABEL_KEY = 'reloader.stakater.com/meta-info';
export const META_INFO_LABEL_VALUE = 'reloader';

const ENV_VAR_PREFIX = 'STAKATER_';
const ENV_VAR_POSTFIXES: Record<string, ResourceKind> = {
  CONFIGMAP: 'ConfigMap',
  SECRET: 'Secret',
  SECRETPROVIDERCLASS: 'SecretProviderClass',
};

export type ResourceKind = 'ConfigMap' | 'Secret' | 'SecretProviderClass';

export type ReloaderMetaInfo = {
  releaseVersion?: string;
  edition?: string;
  reloadStrategy?: string;
  autoReloadAll: boolean;
  namespace?: string;
  deploymentName?: string;
  annotations: AnnotationConfig;
};

export type WatchedResources = {
  // source is where Reloader read the annotations for this resource type from.
  source: 'workload' | 'podTemplate' | 'none';
  auto: boolean;
  search: boolean;
  named: string[];
  excluded: string[];
  // autoReloadAll is true when the --auto-reload-all flag is the only reason this type triggers.
  autoReloadAll: boolean;
};

export type ReloadConfig = {
  // managed is true when a change to at least one resource type can restart the workload.
  managed: boolean;
  configMaps: WatchedResources;
  secrets: WatchedResources;
  secretProviderClasses: WatchedResources;
  pausePeriod?: string;
  pausedAt?: string;
  rolloutStrategy?: string;
};

export type LastReload = {
  kind: string;
  name: string;
  namespace?: string;
  hash?: string;
  containers: string[];
  reloadedAt?: string;
};

export type EnvVarTrigger = {
  container: string;
  envVar: string;
  kind: ResourceKind;
  // resourceName is the env var form of the resource name (upper case, '_' separators), not the real name.
  resourceName: string;
  hash: string;
};

type Container = {
  name: string;
  env?: { name: string; value?: string }[];
};

export type PodTemplateWorkload = K8sResourceCommon & {
  spec?: {
    template?: {
      metadata?: { annotations?: Record<string, string> };
      spec?: { containers?: Container[]; initContainers?: Container[] };
    };
  };
};

const splitList = (value?: string): string[] =>
  (value ?? '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean);

const parseJSON = <T>(value?: string): T | undefined => {
  if (!value) {
    return undefined;
  }
  try {
    return JSON.parse(value) as T;
  } catch {
    return undefined;
  }
};

// parseMetaInfo reads the data of the reloader-meta-info ConfigMap that Reloader publishes on start.
export const parseMetaInfo = (data?: Record<string, string>): ReloaderMetaInfo => {
  const buildInfo = parseJSON<{ releaseVersion?: string; edition?: string }>(data?.buildInfo);
  const config = parseJSON<{
    reloadStrategy?: string;
    autoReloadAll?: boolean;
    annotations?: Partial<AnnotationConfig>;
  }>(data?.config);
  const deploymentInfo = parseJSON<{ name?: string; namespace?: string }>(data?.deploymentInfo);

  const annotations = { ...DEFAULT_ANNOTATIONS };
  for (const [key, value] of Object.entries(config?.annotations ?? {})) {
    if (key in annotations && value) {
      annotations[key as keyof AnnotationConfig] = value;
    }
  }

  return {
    releaseVersion: buildInfo?.releaseVersion,
    edition: buildInfo?.edition,
    reloadStrategy: config?.reloadStrategy,
    autoReloadAll: config?.autoReloadAll ?? false,
    namespace: deploymentInfo?.namespace,
    deploymentName: deploymentInfo?.name,
    annotations,
  };
};

type TypeKeys = { auto: string; reload: string; exclude: string };

// getReloadConfig mirrors pkg/matcher: for each resource type Reloader reads either the workload
// annotations or, when those carry no relevant key, the pod template annotations, never a merge.
export const getReloadConfig = (
  workload: PodTemplateWorkload,
  keys: AnnotationConfig = DEFAULT_ANNOTATIONS,
  autoReloadAll = false,
): ReloadConfig => {
  const workloadAnnotations = workload.metadata?.annotations ?? {};
  const podAnnotations = workload.spec?.template?.metadata?.annotations ?? {};

  const watched = ({ auto, reload, exclude }: TypeKeys): WatchedResources => {
    const relevant = (a: Record<string, string>): boolean =>
      [reload, keys.search, keys.auto, auto].some((k) => k in a);
    let source: WatchedResources['source'] = 'none';
    let annotations: Record<string, string> = workloadAnnotations;
    if (relevant(workloadAnnotations)) {
      source = 'workload';
    } else if (relevant(podAnnotations)) {
      source = 'podTemplate';
      annotations = podAnnotations;
    }

    const annotated =
      annotations[keys.auto] === 'true' ||
      annotations[auto] === 'true' ||
      annotations[keys.search] === 'true' ||
      splitList(annotations[reload]).length > 0;
    const optedOut = annotations[keys.auto] === 'false' || annotations[auto] === 'false';

    return {
      source,
      auto: annotations[keys.auto] === 'true' || annotations[auto] === 'true',
      search: annotations[keys.search] === 'true',
      named: splitList(annotations[reload]),
      excluded: splitList(annotations[exclude]),
      autoReloadAll: autoReloadAll && !annotated && !optedOut,
    };
  };

  const configMaps = watched({
    auto: keys.configmapAuto,
    reload: keys.configmapReload,
    exclude: keys.configmapExclude,
  });
  const secrets = watched({
    auto: keys.secretAuto,
    reload: keys.secretReload,
    exclude: keys.secretExclude,
  });
  const secretProviderClasses = watched({
    auto: keys.secretProviderClassAuto,
    reload: keys.secretProviderClassReload,
    exclude: keys.secretProviderClassExclude,
  });

  const pick = (key: string): string | undefined => workloadAnnotations[key] ?? podAnnotations[key];

  return {
    managed: [configMaps, secrets, secretProviderClasses].some(isWatching),
    configMaps,
    secrets,
    secretProviderClasses,
    pausePeriod: pick(keys.pausePeriod),
    pausedAt: pick(keys.pausedAt),
    rolloutStrategy: pick(keys.rolloutStrategy),
  };
};

export const isWatching = (w: WatchedResources): boolean =>
  w.auto || w.search || w.named.length > 0 || w.autoReloadAll;

// getLastReload decodes the last-reloaded-from annotation written by the annotations strategy.
// It accepts both the v2 shape (kind, containers, reloadedAt) and the v1 shape
// (type, containerRefs, observedAt in unix seconds).
export const getLastReload = (
  workload: PodTemplateWorkload,
  key: string = DEFAULT_ANNOTATIONS.lastReloadedFrom,
): LastReload | undefined => {
  const raw = parseJSON<{
    kind?: string;
    type?: string;
    name?: string;
    namespace?: string;
    hash?: string;
    containers?: string[];
    containerRefs?: string[];
    reloadedAt?: string;
    observedAt?: number;
  }>(workload.spec?.template?.metadata?.annotations?.[key]);
  if (!raw?.name) {
    return undefined;
  }

  let reloadedAt = raw.reloadedAt;
  if (!reloadedAt && raw.observedAt) {
    reloadedAt = new Date(raw.observedAt * 1000).toISOString();
  }

  return {
    kind: raw.kind ?? kindFromV1Type(raw.type),
    name: raw.name,
    namespace: raw.namespace,
    hash: raw.hash,
    containers: raw.containers ?? raw.containerRefs ?? [],
    reloadedAt,
  };
};

const kindFromV1Type = (type?: string): string =>
  (type && ENV_VAR_POSTFIXES[type.toUpperCase()]) || type || '';

// getEnvVarTriggers lists the STAKATER_<NAME>_<TYPE> env vars injected by the env-vars strategy.
export const getEnvVarTriggers = (workload: PodTemplateWorkload): EnvVarTrigger[] => {
  const podSpec = workload.spec?.template?.spec;
  const containers = [...(podSpec?.initContainers ?? []), ...(podSpec?.containers ?? [])];
  const triggers: EnvVarTrigger[] = [];

  for (const container of containers) {
    for (const env of container.env ?? []) {
      const trigger = parseEnvVarName(env.name);
      if (trigger) {
        triggers.push({ container: container.name, envVar: env.name, hash: env.value ?? '', ...trigger });
      }
    }
  }
  return triggers;
};

const parseEnvVarName = (name: string): { kind: ResourceKind; resourceName: string } | undefined => {
  if (!name.startsWith(ENV_VAR_PREFIX)) {
    return undefined;
  }
  const rest = name.slice(ENV_VAR_PREFIX.length);
  const sep = rest.lastIndexOf('_');
  if (sep <= 0) {
    return undefined;
  }
  const kind = ENV_VAR_POSTFIXES[rest.slice(sep + 1)];
  if (!kind) {
    return undefined;
  }
  return { kind, resourceName: rest.slice(0, sep) };
};
