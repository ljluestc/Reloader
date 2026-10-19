/*
 * The dynamic plugin SDK components and hooks are provided by the console at runtime through
 * module federation, so tests use these minimal stubs instead.
 */
import type * as SDK from '@openshift-console/dynamic-plugin-sdk';

export const ListPageHeader: typeof SDK.ListPageHeader = ({ title }) => <h1>{title}</h1>;

export const DocumentTitle: typeof SDK.DocumentTitle = () => null;

export const NamespaceBar: typeof SDK.NamespaceBar = () => null;

export const ResourceLink: typeof SDK.ResourceLink = ({ groupVersionKind, name }) => (
  <span data-test="resource-link">{`${groupVersionKind?.kind}/${name}`}</span>
);

export const Timestamp: typeof SDK.Timestamp = ({ timestamp }) => <span>{timestamp}</span>;

export const useActiveNamespace = jest.fn(() => ['#ALL_NS#', jest.fn()]);

export const useK8sWatchResource = jest.fn(() => [[], true, undefined]);
