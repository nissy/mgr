# mgr
mgr is redis migrations command.

### config

mgr.sample.yaml

```yaml
to_redis:
- source_file: "./testdata/dump.rdb"
  address: "localhost:6379"
  migrates:
  - source_db: 0
    to_db: 1
  - source_db: 1
    to_db: 2
  - source_db: 2
    to_db: 3
```
