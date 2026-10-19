import { render, screen } from '@testing-library/react';
import { useK8sWatchResource } from '@openshift-console/dynamic-plugin-sdk';
import WorkloadReloaderTab from './WorkloadReloaderTab';
import { PodTemplateWorkload } from '../utils/reloader';

const mockWatch = useK8sWatchResource as jest.Mock;

const workload = (annotations: Record<string, string>, podAnnotations: Record<string, string> = {}): PodTemplateWorkload => ({
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: { name: 'app', namespace: 'demo', annotations },
  spec: {
    template: {
      metadata: { annotations: podAnnotations },
      spec: { containers: [{ name: 'app', env: [{ name: 'STAKATER_SETTINGS_CONFIGMAP', value: 'abc123' }] }] },
    },
  },
});

describe('WorkloadReloaderTab', () => {
  beforeEach(() => mockWatch.mockReturnValue([[], true, undefined]));

  it('explains how to opt in when the workload is not managed', () => {
    render(<WorkloadReloaderTab obj={workload({})} />);
    expect(screen.getByText('Not managed by Reloader')).toBeInTheDocument();
    expect(screen.getByText('Reloader configuration not found')).toBeInTheDocument();
  });

  it('shows named triggers, the last reload and injected env vars', () => {
    const lastReload = JSON.stringify({
      kind: 'Secret',
      name: 'db',
      namespace: 'demo',
      hash: 'feed',
      containers: ['app'],
      reloadedAt: '2026-10-01T12:00:00Z',
    });
    render(
      <WorkloadReloaderTab
        obj={workload(
          { 'configmap.reloader.stakater.com/reload': 'settings' },
          { 'reloader.stakater.com/last-reloaded-from': lastReload },
        )}
      />,
    );
    expect(screen.getByText('settings')).toBeInTheDocument();
    expect(screen.getByText('Secret/db')).toBeInTheDocument();
    expect(screen.getByText('2026-10-01T12:00:00Z')).toBeInTheDocument();
    expect(screen.getByText('STAKATER_SETTINGS_CONFIGMAP')).toBeInTheDocument();
  });

  it('uses annotation keys from the reloader-meta-info ConfigMap', () => {
    mockWatch.mockReturnValue([
      [
        {
          metadata: { name: 'reloader-meta-info', namespace: 'reloader' },
          data: { config: JSON.stringify({ annotations: { auto: 'example.com/auto' } }) },
        },
      ],
      true,
      undefined,
    ]);
    render(<WorkloadReloaderTab obj={workload({ 'example.com/auto': 'true' })} />);
    expect(screen.queryByText('Not managed by Reloader')).not.toBeInTheDocument();
    expect(screen.queryByText('Reloader configuration not found')).not.toBeInTheDocument();
    expect(screen.getAllByText('Any referenced')).toHaveLength(3);
  });
});
