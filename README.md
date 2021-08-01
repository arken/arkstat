# arkstat
A Minimal Statistics API for Arken Clusters

## Keyset Config
```yaml
name: Core Keyset
replications: 5
cluster_key: 793bdb68b7cfd2f49071a299711df51f1c60283a047e4a8756a5c3a3d1ab776f
stats_node: 12D3KooWL7hvR7nfQxAWMowgoWXWQwKEkQA8QPZrhKjateRTgcDm
bootstrap_nodes:
    - /dns4/bs1.arken.io/tcp/4001/ipfs/12D3KooWSmosHZtDBbepxWwVgo8HyXSgNCUgs2GGD2qnQPbA3KhD
    - /dns4/bs2.arken.io/tcp/4001/ipfs/12D3KooWSmosHZtDBbepxWwVgo8HyXSgNCUgs2BbepxWwVgBbepx
mirrors:
    - https://gitlab.com/arkenproject/core-keyset.git
```

## Node Config
```yaml
general:
  pool_size: 50 GB
  network_limit: 50 GB
stats:
  email: human@example.com
  enabled: true
database:
  path: ~/.arken/files.db
storage:
  path: ~/.arken/storage/

```