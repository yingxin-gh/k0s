# Version skew policy

k0s ships a curated set of Kubernetes components (kube-apiserver,
kube-controller-manager, kube-scheduler, kubelet, kube-proxy, …) for each
release, so the version skew between those components inside a single k0s
release is always supported by upstream Kubernetes by construction.

What you still control as an operator is the skew **between releases of
k0s** that you mix in one cluster — for example, controllers running k0s
`v1.31.x+k0s.0` while workers still run `v1.30.x+k0s.0` during a rolling
upgrade.

## Supported skew

k0s follows the upstream Kubernetes
[version skew policy](https://kubernetes.io/releases/version-skew-policy/).
In practice, that means:

- **Control plane components** (kube-apiserver, kube-controller-manager,
  kube-scheduler) must all be within **one minor** of each other in a
  multi-controller cluster.
- **kubelet** on a worker may be **up to three minors older** than the
  oldest kube-apiserver in the cluster, but **must never be newer**.
- **kube-proxy** on a worker must match the kubelet's minor version.

Patch-version differences inside the same minor are always supported and
are the expected state during a phased upgrade.

## Upgrades

When upgrading k0s in place:

1. Upgrade controllers first, **one minor at a time** (e.g. v1.30 → v1.31,
   not v1.30 → v1.32).
2. Upgrade workers after the control plane is fully on the new minor.
3. Stay on the latest patch release of the old minor before stepping up to
   the next minor.

[Autopilot](autopilot.md) automates this with multi-step plans; if you
upgrade manually or via [k0sctl](k0sctl-install.md), make the same staging
explicit in your plan.

## Downgrades

Downgrades are **not** supported by upstream Kubernetes and are not tested
by k0s. If a release needs to be rolled back, restore from a
[backup](backup.md) taken before the upgrade rather than running an older
k0s binary against a newer cluster state.

## Major versions

There is currently no Kubernetes v2 and no k0s v2, so major-version skew is
not a configuration that exists in practice.
