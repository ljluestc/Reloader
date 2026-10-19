import {
  DocumentTitle,
  K8sGroupVersionKind,
  ListPageHeader,
  NamespaceBar,
  ResourceLink,
  Timestamp,
  useActiveNamespace,
  useK8sWatchResource,
} from '@openshift-console/dynamic-plugin-sdk';
import {
  Alert,
  Bullseye,
  Content,
  EmptyState,
  EmptyStateBody,
  PageSection,
  SearchInput,
  Spinner,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
} from '@patternfly/react-core';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { FC } from 'react';
import { useReloaderMetaInfo } from '../hooks/useReloaderMetaInfo';
import { PodTemplateWorkload, getEnvVarTriggers, getLastReload, getReloadConfig } from '../utils/reloader';
import { WatchedResourcesList } from './WatchedResourcesList';

const ALL_NAMESPACES_KEY = '#ALL_NS#';

const WORKLOAD_KINDS: K8sGroupVersionKind[] = [
  { group: 'apps', version: 'v1', kind: 'Deployment' },
  { group: 'apps', version: 'v1', kind: 'StatefulSet' },
  { group: 'apps', version: 'v1', kind: 'DaemonSet' },
];

const useWorkloads = (namespace?: string): [PodTemplateWorkload[], boolean, unknown] => {
  const watch = (groupVersionKind: K8sGroupVersionKind) => ({
    groupVersionKind,
    isList: true,
    namespaced: true,
    namespace,
  });
  // Hooks cannot be called in a loop, so each kind gets its own watch.
  const [deployments, deploymentsLoaded, deploymentsError] = useK8sWatchResource<PodTemplateWorkload[]>(
    watch(WORKLOAD_KINDS[0]),
  );
  const [statefulSets, statefulSetsLoaded, statefulSetsError] = useK8sWatchResource<PodTemplateWorkload[]>(
    watch(WORKLOAD_KINDS[1]),
  );
  const [daemonSets, daemonSetsLoaded, daemonSetsError] = useK8sWatchResource<PodTemplateWorkload[]>(
    watch(WORKLOAD_KINDS[2]),
  );

  return [
    useMemo(
      () => [
        ...(deployments ?? []).map((w) => ({ ...w, kind: 'Deployment' })),
        ...(statefulSets ?? []).map((w) => ({ ...w, kind: 'StatefulSet' })),
        ...(daemonSets ?? []).map((w) => ({ ...w, kind: 'DaemonSet' })),
      ],
      [deployments, statefulSets, daemonSets],
    ),
    deploymentsLoaded && statefulSetsLoaded && daemonSetsLoaded,
    deploymentsError || statefulSetsError || daemonSetsError,
  ];
};

const gvkFor = (kind?: string): K8sGroupVersionKind =>
  WORKLOAD_KINDS.find((gvk) => gvk.kind === kind) ?? WORKLOAD_KINDS[0];

const ReloaderPage: FC = () => {
  const { t } = useTranslation('plugin__reloader-console-plugin');
  const [activeNamespace] = useActiveNamespace();
  const namespace = activeNamespace === ALL_NAMESPACES_KEY ? undefined : activeNamespace;
  const [filter, setFilter] = useState('');

  const { metaInfo, found } = useReloaderMetaInfo();
  const [workloads, loaded, loadError] = useWorkloads(namespace);

  const rows = useMemo(
    () =>
      workloads
        .map((workload) => ({
          workload,
          config: getReloadConfig(workload, metaInfo.annotations, metaInfo.autoReloadAll),
          lastReload: getLastReload(workload, metaInfo.annotations.lastReloadedFrom),
          envVarTriggers: getEnvVarTriggers(workload).length,
        }))
        .filter(({ config }) => config.managed)
        .filter(({ workload }) => !filter || workload.metadata?.name?.includes(filter))
        .sort(
          (a, b) =>
            (a.workload.metadata?.namespace ?? '').localeCompare(b.workload.metadata?.namespace ?? '') ||
            (a.workload.metadata?.name ?? '').localeCompare(b.workload.metadata?.name ?? ''),
        ),
    [workloads, metaInfo, filter],
  );

  const columns = {
    name: t('Name'),
    namespace: t('Namespace'),
    configMaps: t('ConfigMaps'),
    secrets: t('Secrets'),
    lastReload: t('Last reload'),
  };

  return (
    <>
      <DocumentTitle>{t('Reloader')}</DocumentTitle>
      <NamespaceBar />
      <ListPageHeader title={t('Reloader')} />
      <PageSection>
        <Content component="p">
          {t('Workloads that Reloader restarts when their ConfigMaps, Secrets or SecretProviderClasses change.')}
        </Content>
        {found ? (
          <Content component="small" className="pf-v6-u-mb-md">
            {t('Reloader {{version}} in namespace {{namespace}}, {{strategy}} strategy', {
              version: metaInfo.releaseVersion ?? '',
              namespace: metaInfo.namespace ?? '',
              strategy: metaInfo.reloadStrategy ?? '',
            })}
            {metaInfo.autoReloadAll && `, ${t('auto-reload-all enabled')}`}
          </Content>
        ) : (
          <Alert isInline isPlain variant="info" title={t('Reloader configuration not found')} className="pf-v6-u-mb-md">
            {t(
              'The reloader-meta-info ConfigMap is not visible to you, so the default annotation keys are assumed.',
            )}
          </Alert>
        )}

        <Toolbar>
          <ToolbarContent>
            <ToolbarItem>
              <SearchInput
                aria-label={t('Filter by name')}
                placeholder={t('Filter by name')}
                value={filter}
                onChange={(_event, value) => setFilter(value)}
                onClear={() => setFilter('')}
              />
            </ToolbarItem>
          </ToolbarContent>
        </Toolbar>

        {loadError ? (
          <Alert isInline variant="danger" title={t('Unable to load workloads')}>
            {String((loadError as Error)?.message ?? loadError)}
          </Alert>
        ) : !loaded ? (
          <Bullseye>
            <Spinner aria-label={t('Loading')} />
          </Bullseye>
        ) : rows.length === 0 ? (
          <EmptyState titleText={t('No workloads managed by Reloader')} headingLevel="h2">
            <EmptyStateBody>
              {t('Annotate a Deployment, StatefulSet or DaemonSet with {{annotation}} to have it listed here.', {
                annotation: `${metaInfo.annotations.auto}: "true"`,
              })}
            </EmptyStateBody>
          </EmptyState>
        ) : (
          <Table aria-label={t('Reloader managed workloads')} variant="compact">
            <Thead>
              <Tr>
                <Th>{columns.name}</Th>
                {!namespace && <Th>{columns.namespace}</Th>}
                <Th>{columns.configMaps}</Th>
                <Th>{columns.secrets}</Th>
                <Th>{columns.lastReload}</Th>
              </Tr>
            </Thead>
            <Tbody>
              {rows.map(({ workload, config, lastReload, envVarTriggers }) => (
                <Tr key={workload.metadata?.uid ?? `${workload.kind}/${workload.metadata?.namespace}/${workload.metadata?.name}`}>
                  <Td dataLabel={columns.name}>
                    <ResourceLink
                      groupVersionKind={gvkFor(workload.kind)}
                      name={workload.metadata?.name}
                      namespace={workload.metadata?.namespace}
                    />
                  </Td>
                  {!namespace && (
                    <Td dataLabel={columns.namespace}>
                      <ResourceLink groupVersionKind={{ version: 'v1', kind: 'Namespace' }} name={workload.metadata?.namespace} />
                    </Td>
                  )}
                  <Td dataLabel={columns.configMaps}>
                    <WatchedResourcesList watched={config.configMaps} />
                  </Td>
                  <Td dataLabel={columns.secrets}>
                    <WatchedResourcesList watched={config.secrets} />
                  </Td>
                  <Td dataLabel={columns.lastReload}>
                    {lastReload ? (
                      <>
                        {lastReload.kind}/{lastReload.name}
                        {lastReload.reloadedAt && <Timestamp timestamp={lastReload.reloadedAt} simple />}
                      </>
                    ) : envVarTriggers > 0 ? (
                      t('Env var triggers: {{count}}', { count: envVarTriggers })
                    ) : (
                      '-'
                    )}
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        )}
      </PageSection>
    </>
  );
};

export default ReloaderPage;
