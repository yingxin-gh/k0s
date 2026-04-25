# k0s directories

This page describes the directories k0s reads and writes on the host, what
they are used for, and the permissions k0s expects on them.

The two parent directories are configurable via CLI flags:

- `--data-dir` (default `/var/lib/k0s` on Linux, `C:\var\lib\k0s` on Windows) —
  persistent state.
- `--run-dir` (default `/run/k0s` on Linux) — runtime state such as PID files
  and unix sockets.

Pass these flags to every `k0s` invocation that should share the same state
(controller, worker, status, reset, …).

## Data directory layout

Inside `<data-dir>` (default `/var/lib/k0s`):

| Path                           | Purpose                                                                                | Mode   |
| ------------------------------ | -------------------------------------------------------------------------------------- | ------ |
| `<data-dir>/`                  | Top-level data directory.                                                              | `0755` |
| `<data-dir>/bin/`              | Embedded binaries (kubelet, etcd, containerd, …) extracted on first start.             | `0755` |
| `<data-dir>/pki/`              | Cluster CA and component certificates managed by k0s.                                  | `0751` |
| `<data-dir>/pki/etcd/`         | etcd-specific certificates.                                                            | `0711` |
| `<data-dir>/etcd/`             | etcd data store. Backed up and restored by `k0s backup` / `k0s restore`.               | `0700` |
| `<data-dir>/manifests/`        | Manifest Deployer drop-in directory. Files placed here are applied to the cluster.     | `0755` |
| `<data-dir>/images/`           | Image bundles imported by the embedded container runtime on start-up.                  | `0755` |
| `<data-dir>/kubelet/`          | kubelet root directory (volumes, plugins, pod state). CSI drivers must be aware of it. | varies |
| `<data-dir>/kine/`             | Kine database files when k0s is configured to use Kine instead of etcd.                | `0750` |

Individual certificate files inside `pki/` are written with mode `0644`
(public certs) or `0640` (private keys), and admin/kubelet config files are
written with mode `0600`.

## Run directory layout

Inside `<run-dir>` (default `/run/k0s`):

| Path                  | Purpose                                                              | Mode   |
| --------------------- | -------------------------------------------------------------------- | ------ |
| `<run-dir>/`          | Top-level run directory.                                             | `0755` |
| `<run-dir>/*.pid`     | PID files for supervised processes (kubelet, etcd, containerd, …).   | `0644` |
| `<run-dir>/status.sock` | Unix socket exposing local status to `k0s status` and similar commands. | varies |

The run directory is intentionally outside the data directory so that on
systems where `/run` is a tmpfs, runtime state is cleared on reboot while
persistent state under `<data-dir>` survives.

## Notes

- All paths above are derived from the same `--data-dir` / `--run-dir` flags
  passed to `k0s` and propagated to its supervised components, so changing
  the flags moves the whole tree consistently.
- The CSI guide in [storage.md](storage.md) documents the kubelet root path
  caveat, since the default differs from upstream Kubernetes
  (`<data-dir>/kubelet` rather than `/var/lib/kubelet`).
- Backup/restore covers `pki/`, `etcd/`, `manifests/`, and `images/`. See
  [backup.md](backup.md) for the authoritative list.
