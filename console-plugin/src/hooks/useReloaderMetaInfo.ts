import { useMemo } from 'react';
import { K8sResourceCommon, useK8sWatchResource } from '@openshift-console/dynamic-plugin-sdk';
import { META_INFO_LABEL_KEY, META_INFO_LABEL_VALUE, ReloaderMetaInfo, parseMetaInfo } from '../utils/reloader';

type ConfigMap = K8sResourceCommon & { data?: Record<string, string> };

export type UseReloaderMetaInfoResult = {
  metaInfo: ReloaderMetaInfo;
  // found is false when no reloader-meta-info ConfigMap is visible to the user, in which case
  // metaInfo holds the default annotation keys.
  found: boolean;
  loaded: boolean;
};

// useReloaderMetaInfo watches the reloader-meta-info ConfigMap Reloader publishes in its own
// namespace. Users who cannot list ConfigMaps there fall back to the default annotation keys.
export const useReloaderMetaInfo = (): UseReloaderMetaInfoResult => {
  const [configMaps, loaded, loadError] = useK8sWatchResource<ConfigMap[]>({
    groupVersionKind: { version: 'v1', kind: 'ConfigMap' },
    isList: true,
    selector: { matchLabels: { [META_INFO_LABEL_KEY]: META_INFO_LABEL_VALUE } },
  });

  return useMemo(() => {
    const configMap = !loadError && configMaps?.length ? configMaps[0] : undefined;
    return {
      metaInfo: parseMetaInfo(configMap?.data),
      found: !!configMap,
      loaded: loaded || !!loadError,
    };
  }, [configMaps, loaded, loadError]);
};
