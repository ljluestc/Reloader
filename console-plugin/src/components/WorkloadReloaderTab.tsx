import { ResourceLink, Timestamp } from '@openshift-console/dynamic-plugin-sdk';
import {
  Alert,
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  EmptyState,
  EmptyStateBody,
  PageSection,
  Title,
} from '@patternfly/react-core';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { useTranslation } from 'react-i18next';
import type { FC } from 'react';
import { useReloaderMetaInfo } from '../hooks/useReloaderMetaInfo';
import {
  PodTemplateWorkload,
  WatchedResources,
  getEnvVarTriggers,
  getLastReload,
  getReloadConfig,
} from '../utils/reloader';
import { WatchedResourcesList } from './WatchedResourcesList';

type WorkloadReloaderTabProps = {
  obj: PodTemplateWorkload;
};

const resourceGVK = (kind: string) =>
  kind === 'SecretProviderClass'
    ? { group: 'secrets-store.csi.x-k8s.io', version: 'v1', kind }
    : { version: 'v1', kind };

const WorkloadReloaderTab: FC<WorkloadReloaderTabProps> = ({ obj }) => {
  const { t } = useTranslation('plugin__reloader-console-plugin');
  const { metaInfo, found, loaded } = useReloaderMetaInfo();
  const keys = metaInfo.annotations;

  const config = getReloadConfig(obj, keys, metaInfo.autoReloadAll);
  const lastReload = getLastReload(obj, keys.lastReloadedFrom);
  const envVarTriggers = getEnvVarTriggers(obj);

  const sourceLabel = (w: WatchedResources): string => {
    switch (w.source) {
      case 'workload':
        return t('Workload annotations');
      case 'podTemplate':
        return t('Pod template annotations');
      default:
        return t('None');
    }
  };

  const types: [string, WatchedResources][] = [
    [t('ConfigMaps'), config.configMaps],
    [t('Secrets'), config.secrets],
    [t('SecretProviderClasses'), config.secretProviderClasses],
  ];

  return (
    <PageSection>
      {loaded && !found && (
        <Alert
          isInline
          isPlain
          variant="info"
          title={t('Reloader configuration not found')}
          className="pf-v6-u-mb-md"
        >
          {t(
            'The reloader-meta-info ConfigMap is not visible to you, so the default annotation keys are assumed.',
          )}
        </Alert>
      )}

      {!config.managed ? (
        <EmptyState titleText={t('Not managed by Reloader')} headingLevel="h2">
          <EmptyStateBody>
            {t('Add the {{annotation}} annotation to restart this workload when its ConfigMaps or Secrets change.', {
              annotation: `${keys.auto}: "true"`,
            })}
          </EmptyStateBody>
        </EmptyState>
      ) : (
        <>
          <Title headingLevel="h2" className="pf-v6-u-mb-md">
            {t('Reload triggers')}
          </Title>
          <DescriptionList isHorizontal>
            {types.map(([label, watched]) => (
              <DescriptionListGroup key={label}>
                <DescriptionListTerm>{label}</DescriptionListTerm>
                <DescriptionListDescription>
                  <WatchedResourcesList watched={watched} />
                  {watched.source !== 'none' && (
                    <Content component="small">{t('Read from: {{source}}', { source: sourceLabel(watched) })}</Content>
                  )}
                </DescriptionListDescription>
              </DescriptionListGroup>
            ))}
            {config.pausePeriod && (
              <DescriptionListGroup>
                <DescriptionListTerm>{t('Pause period')}</DescriptionListTerm>
                <DescriptionListDescription>
                  {config.pausePeriod}
                  {config.pausedAt && (
                    <>
                      {' '}
                      ({t('paused')} <Timestamp timestamp={config.pausedAt} />)
                    </>
                  )}
                </DescriptionListDescription>
              </DescriptionListGroup>
            )}
            {config.rolloutStrategy && (
              <DescriptionListGroup>
                <DescriptionListTerm>{t('Rollout strategy')}</DescriptionListTerm>
                <DescriptionListDescription>{config.rolloutStrategy}</DescriptionListDescription>
              </DescriptionListGroup>
            )}
          </DescriptionList>
        </>
      )}

      <Title headingLevel="h2" className="pf-v6-u-mt-xl pf-v6-u-mb-md">
        {t('Last reload')}
      </Title>
      {lastReload ? (
        <DescriptionList isHorizontal>
          <DescriptionListGroup>
            <DescriptionListTerm>{t('Triggered by')}</DescriptionListTerm>
            <DescriptionListDescription>
              <ResourceLink
                groupVersionKind={resourceGVK(lastReload.kind)}
                name={lastReload.name}
                namespace={lastReload.namespace ?? obj.metadata?.namespace}
              />
            </DescriptionListDescription>
          </DescriptionListGroup>
          {lastReload.reloadedAt && (
            <DescriptionListGroup>
              <DescriptionListTerm>{t('Reloaded at')}</DescriptionListTerm>
              <DescriptionListDescription>
                <Timestamp timestamp={lastReload.reloadedAt} />
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {lastReload.containers.length > 0 && (
            <DescriptionListGroup>
              <DescriptionListTerm>{t('Containers')}</DescriptionListTerm>
              <DescriptionListDescription>{lastReload.containers.join(', ')}</DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {lastReload.hash && (
            <DescriptionListGroup>
              <DescriptionListTerm>{t('Data hash')}</DescriptionListTerm>
              <DescriptionListDescription>
                <code>{lastReload.hash}</code>
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
        </DescriptionList>
      ) : (
        <Content component="p">
          {envVarTriggers.length > 0
            ? t('Reloader uses the env-vars strategy here, which records hashes but no reload time.')
            : t('No reload has been recorded on this workload.')}
        </Content>
      )}

      {envVarTriggers.length > 0 && (
        <>
          <Title headingLevel="h2" className="pf-v6-u-mt-xl pf-v6-u-mb-md">
            {t('Injected environment variables')}
          </Title>
          <Table aria-label={t('Injected environment variables')} variant="compact">
            <Thead>
              <Tr>
                <Th>{t('Container')}</Th>
                <Th>{t('Resource type')}</Th>
                <Th>{t('Environment variable')}</Th>
                <Th>{t('Data hash')}</Th>
              </Tr>
            </Thead>
            <Tbody>
              {envVarTriggers.map((trigger) => (
                <Tr key={`${trigger.container}/${trigger.envVar}`}>
                  <Td dataLabel={t('Container')}>{trigger.container}</Td>
                  <Td dataLabel={t('Resource type')}>{trigger.kind}</Td>
                  <Td dataLabel={t('Environment variable')}>
                    <code>{trigger.envVar}</code>
                  </Td>
                  <Td dataLabel={t('Data hash')} modifier="truncate">
                    <code>{trigger.hash}</code>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </>
      )}
    </PageSection>
  );
};

export default WorkloadReloaderTab;
