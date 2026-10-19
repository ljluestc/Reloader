import {
  DEFAULT_ANNOTATIONS,
  PodTemplateWorkload,
  getEnvVarTriggers,
  getLastReload,
  getReloadConfig,
  parseMetaInfo,
} from './reloader';

const deployment = (
  annotations: Record<string, string> = {},
  podAnnotations: Record<string, string> = {},
): PodTemplateWorkload => ({
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: { name: 'app', namespace: 'demo', annotations },
  spec: { template: { metadata: { annotations: podAnnotations }, spec: { containers: [] } } },
});

describe('getReloadConfig', () => {
  it('is unmanaged without Reloader annotations', () => {
    const config = getReloadConfig(deployment());
    expect(config.managed).toBe(false);
    expect(config.configMaps.source).toBe('none');
  });

  it('enables every type with reloader.stakater.com/auto', () => {
    const config = getReloadConfig(deployment({ 'reloader.stakater.com/auto': 'true' }));
    expect(config.managed).toBe(true);
    expect(config.configMaps.auto).toBe(true);
    expect(config.secrets.auto).toBe(true);
    expect(config.secretProviderClasses.auto).toBe(true);
  });

  it('only enables the typed auto annotation for its type', () => {
    const config = getReloadConfig(deployment({ 'secret.reloader.stakater.com/auto': 'true' }));
    expect(config.secrets.auto).toBe(true);
    expect(config.configMaps.auto).toBe(false);
  });

  it('requires the exact value "true"', () => {
    const config = getReloadConfig(deployment({ 'reloader.stakater.com/auto': 'True' }));
    expect(config.managed).toBe(false);
  });

  it('splits named and excluded resources', () => {
    const config = getReloadConfig(
      deployment({
        'configmap.reloader.stakater.com/reload': 'cm-a, cm-b,',
        'configmaps.exclude.reloader.stakater.com/reload': 'cm-c',
      }),
    );
    expect(config.configMaps.named).toEqual(['cm-a', 'cm-b']);
    expect(config.configMaps.excluded).toEqual(['cm-c']);
    expect(config.secrets.named).toEqual([]);
  });

  it('falls back to pod template annotations per resource type, without merging', () => {
    const config = getReloadConfig(
      deployment(
        { 'secret.reloader.stakater.com/reload': 'db' },
        {
          'configmap.reloader.stakater.com/reload': 'settings',
          'secret.reloader.stakater.com/reload': 'ignored-because-workload-wins',
        },
      ),
    );
    expect(config.secrets.source).toBe('workload');
    expect(config.secrets.named).toEqual(['db']);
    expect(config.configMaps.source).toBe('podTemplate');
    expect(config.configMaps.named).toEqual(['settings']);
  });

  it('reads excludes only from the selected annotation source', () => {
    const config = getReloadConfig(
      deployment(
        { 'reloader.stakater.com/auto': 'true' },
        { 'configmaps.exclude.reloader.stakater.com/reload': 'cm-c' },
      ),
    );
    expect(config.configMaps.excluded).toEqual([]);
  });

  it('reports search mode', () => {
    const config = getReloadConfig(deployment({ 'reloader.stakater.com/search': 'true' }));
    expect(config.managed).toBe(true);
    expect(config.configMaps.search).toBe(true);
  });

  it('marks workloads managed only through --auto-reload-all', () => {
    const config = getReloadConfig(deployment(), DEFAULT_ANNOTATIONS, true);
    expect(config.managed).toBe(true);
    expect(config.configMaps.autoReloadAll).toBe(true);
  });

  it('honours auto "false" opt-outs from --auto-reload-all', () => {
    const all = getReloadConfig(deployment({ 'reloader.stakater.com/auto': 'false' }), DEFAULT_ANNOTATIONS, true);
    expect(all.managed).toBe(false);

    const typed = getReloadConfig(
      deployment({ 'secret.reloader.stakater.com/auto': 'false' }),
      DEFAULT_ANNOTATIONS,
      true,
    );
    expect(typed.secrets.autoReloadAll).toBe(false);
    expect(typed.configMaps.autoReloadAll).toBe(true);
  });

  it('uses custom annotation keys', () => {
    const keys = { ...DEFAULT_ANNOTATIONS, auto: 'example.com/auto' };
    expect(getReloadConfig(deployment({ 'example.com/auto': 'true' }), keys).managed).toBe(true);
    expect(getReloadConfig(deployment({ 'reloader.stakater.com/auto': 'true' }), keys).managed).toBe(false);
  });

  it('reads pause annotations', () => {
    const config = getReloadConfig(
      deployment({
        'reloader.stakater.com/auto': 'true',
        'deployment.reloader.stakater.com/pause-period': '5m',
        'deployment.reloader.stakater.com/paused-at': '2026-10-01T00:00:00Z',
      }),
    );
    expect(config.pausePeriod).toBe('5m');
    expect(config.pausedAt).toBe('2026-10-01T00:00:00Z');
  });
});

describe('getLastReload', () => {
  it('returns undefined without the annotation or with invalid JSON', () => {
    expect(getLastReload(deployment())).toBeUndefined();
    expect(
      getLastReload(deployment({}, { 'reloader.stakater.com/last-reloaded-from': '{not json' })),
    ).toBeUndefined();
  });

  it('decodes the v2 format', () => {
    const value = JSON.stringify({
      kind: 'ConfigMap',
      name: 'settings',
      namespace: 'demo',
      hash: 'abc',
      containers: ['app'],
      reloadedAt: '2026-10-01T12:00:00Z',
    });
    expect(getLastReload(deployment({}, { 'reloader.stakater.com/last-reloaded-from': value }))).toEqual({
      kind: 'ConfigMap',
      name: 'settings',
      namespace: 'demo',
      hash: 'abc',
      containers: ['app'],
      reloadedAt: '2026-10-01T12:00:00Z',
    });
  });

  it('decodes the v1 format', () => {
    const value = JSON.stringify({
      type: 'SECRET',
      name: 'db',
      namespace: 'demo',
      hash: 'def',
      containerRefs: ['app'],
      observedAt: 1790000000,
    });
    const last = getLastReload(deployment({}, { 'reloader.stakater.com/last-reloaded-from': value }));
    expect(last?.kind).toBe('Secret');
    expect(last?.containers).toEqual(['app']);
    expect(last?.reloadedAt).toBe(new Date(1790000000 * 1000).toISOString());
  });
});

describe('getEnvVarTriggers', () => {
  it('lists STAKATER_ env vars across containers', () => {
    const workload = deployment();
    workload.spec!.template!.spec = {
      initContainers: [{ name: 'init', env: [{ name: 'STAKATER_DB_CREDS_SECRET', value: 'h1' }] }],
      containers: [
        {
          name: 'app',
          env: [
            { name: 'PATH', value: '/bin' },
            { name: 'STAKATER_APP_SETTINGS_CONFIGMAP', value: 'h2' },
            { name: 'STAKATER_VAULT_SECRETPROVIDERCLASS', value: 'h3' },
            { name: 'STAKATER_UNKNOWN_THING', value: 'x' },
            { name: 'STAKATER_CONFIGMAP', value: 'x' },
          ],
        },
      ],
    };
    expect(getEnvVarTriggers(workload)).toEqual([
      { container: 'init', envVar: 'STAKATER_DB_CREDS_SECRET', kind: 'Secret', resourceName: 'DB_CREDS', hash: 'h1' },
      {
        container: 'app',
        envVar: 'STAKATER_APP_SETTINGS_CONFIGMAP',
        kind: 'ConfigMap',
        resourceName: 'APP_SETTINGS',
        hash: 'h2',
      },
      {
        container: 'app',
        envVar: 'STAKATER_VAULT_SECRETPROVIDERCLASS',
        kind: 'SecretProviderClass',
        resourceName: 'VAULT',
        hash: 'h3',
      },
    ]);
  });
});

describe('parseMetaInfo', () => {
  it('falls back to defaults without data', () => {
    const info = parseMetaInfo(undefined);
    expect(info.annotations).toEqual(DEFAULT_ANNOTATIONS);
    expect(info.autoReloadAll).toBe(false);
  });

  it('reads build info, config and custom annotation keys', () => {
    const info = parseMetaInfo({
      buildInfo: JSON.stringify({ releaseVersion: 'v2.0.0', edition: 'oss' }),
      config: JSON.stringify({
        reloadStrategy: 'annotations',
        autoReloadAll: true,
        annotations: { auto: 'example.com/auto', unknownKey: 'x', secretAuto: '' },
      }),
      deploymentInfo: JSON.stringify({ name: 'reloader', namespace: 'reloader' }),
    });
    expect(info.releaseVersion).toBe('v2.0.0');
    expect(info.reloadStrategy).toBe('annotations');
    expect(info.autoReloadAll).toBe(true);
    expect(info.namespace).toBe('reloader');
    expect(info.annotations.auto).toBe('example.com/auto');
    expect(info.annotations.secretAuto).toBe(DEFAULT_ANNOTATIONS.secretAuto);
    expect(info.annotations).not.toHaveProperty('unknownKey');
  });
});
