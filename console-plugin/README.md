# Reloader OpenShift Console Plugin

An [OpenShift Console dynamic plugin](https://github.com/openshift/console/tree/main/frontend/packages/console-dynamic-plugin-sdk)
that shows what Reloader does to your workloads without reading annotations by hand.

## What it adds

- **Workloads → Reloader** page in the Administrator perspective. It lists every Deployment, StatefulSet
  and DaemonSet in the selected namespace (or all namespaces) that Reloader will restart, which
  ConfigMaps and Secrets trigger it, and what triggered the last reload.
- **Reloader** tab on Deployment, StatefulSet and DaemonSet details pages. It shows:
    - the reload triggers per resource type (auto, named, search, excluded, `--auto-reload-all`) and whether they come from the workload or the pod template annotations
    - the pause period
    - the last reload recorded by the `annotations` strategy (resource, time, containers, hash)
    - the `STAKATER_*` variables injected by the `env-vars` strategy

The plugin reads the `reloader-meta-info` ConfigMap that Reloader publishes to pick up the running
reload strategy, `--auto-reload-all`, and any annotation keys overridden with flags. Users who can't
list that ConfigMap see an info notice, and the plugin falls back to the default annotation keys.

The plugin only reads resources with the logged in user's own permissions. It never changes workloads.

## Requirements

- OpenShift 4.23 or later
- Reloader v2, for `reloader-meta-info` and the v2 `last-reloaded-from` format. v1 annotations are also understood.

## Deploy

Build and push the image, then apply the manifests to the namespace Reloader runs in (`reloader` by default):

```bash
docker build -t <registry>/reloader-console-plugin:<tag> console-plugin
docker push <registry>/reloader-console-plugin:<tag>

cd console-plugin/deploy
kustomize edit set image reloader-console-plugin=<registry>/reloader-console-plugin:<tag>
oc apply -k .
```

To use a namespace other than `reloader`, change `namespace` in both `kustomization.yaml` and
`consoleplugin.yaml`.

Enable the plugin on the cluster console:

```bash
oc patch consoles.operator.openshift.io cluster --type=json \
  -p '[{"op": "add", "path": "/spec/plugins/-", "value": "reloader-console-plugin"}]'
```

If `spec.plugins` doesn't exist yet, use `-p '[{"op": "add", "path": "/spec/plugins", "value": ["reloader-console-plugin"]}]'`.

## Development

```bash
cd console-plugin
npm ci
npm test          # unit tests
npm run typecheck
npm run build     # production bundle in dist/
npm start         # dev server on http://localhost:9001
```

To try it against a cluster, run a local console with the plugin proxied to the dev server as
described in the [console plugin template](https://github.com/openshift/console-plugin-template#development).

`src/utils/reloader.ts` mirrors the annotation rules in `pkg/matcher`. Update both together.
