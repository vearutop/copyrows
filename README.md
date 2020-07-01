# copyrows

Copyrows copies data between databases (MySQL and PostgreSQL are supported).

## Installation

```
go get -u github.com/vearutop/copyrows
```

## Usage

```bash
copyrows \
  -src "postgres://foo_bar:pass@foo-bar-service-db000.live.baz.io:5432/foo_bar?sslmode=disable" \
  -dst "postgres://qux_service_admin:pass@qux-service-db000.live.baz.io:5432/qux_service?sslmode=disable" \
  -query "SELECT * FROM allocation WHERE experiment_id = 1 OFFSET 1000 LIMIT 1000" \
  -page-size 100 \
  -table foo_allocation

2020/07/01 17:38:10 rows affected: 100, total inserted: 100
2020/07/01 17:38:10 rows affected: 100, total inserted: 200
2020/07/01 17:38:10 rows affected: 100, total inserted: 300
2020/07/01 17:38:10 rows affected: 100, total inserted: 400
2020/07/01 17:38:10 rows affected: 100, total inserted: 500
2020/07/01 17:38:10 rows affected: 100, total inserted: 600
2020/07/01 17:38:10 rows affected: 100, total inserted: 700
2020/07/01 17:38:10 rows affected: 100, total inserted: 800
2020/07/01 17:38:10 rows affected: 100, total inserted: 900
2020/07/01 17:38:10 rows affected: 100, total inserted: 1000
```