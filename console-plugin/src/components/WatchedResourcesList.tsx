import { Label, LabelGroup } from '@patternfly/react-core';
import { useTranslation } from 'react-i18next';
import type { FC } from 'react';
import { WatchedResources } from '../utils/reloader';

type WatchedResourcesListProps = {
  watched: WatchedResources;
};

// WatchedResourcesList renders how one resource type (ConfigMaps, Secrets, ...) triggers reloads.
export const WatchedResourcesList: FC<WatchedResourcesListProps> = ({ watched }) => {
  const { t } = useTranslation('plugin__reloader-console-plugin');

  const labels = [
    ...(watched.auto ? [<Label key="auto" color="blue">{t('Any referenced')}</Label>] : []),
    ...(watched.search ? [<Label key="search" color="purple">{t('Search: match annotation')}</Label>] : []),
    ...(watched.autoReloadAll ? [<Label key="all" color="blue">{t('--auto-reload-all')}</Label>] : []),
    ...watched.named.map((name) => (
      <Label key={`named-${name}`} color="green">
        {name}
      </Label>
    )),
    ...watched.excluded.map((name) => (
      <Label key={`excluded-${name}`} color="red" variant="outline">
        {t('excluded: {{name}}', { name })}
      </Label>
    )),
  ];

  if (labels.length === 0) {
    return <>{t('Not watched')}</>;
  }
  return <LabelGroup numLabels={6}>{labels}</LabelGroup>;
};
