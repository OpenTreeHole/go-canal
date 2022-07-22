# go-canal
a tool to receive mysql binlog and push data to elastic search

## Usage

### Flags
| Flag    | Default Value | Description                           |
|---------|---------------|---------------------------------------|
| -config | config.yaml   | config file                           |
| -dump   | false         | dump all data before receiving binlog |

### Config.yaml

Note: each table MUST have an `id` field.

```yaml
address: localhost:3306
user: root
password: root
elastic_url: http://localhost:9200

schemas:
  my_db:
    tables:
      my_table_1:
        columns:
          - id
          - my_column
      my_table_2:
        columns:
          - id
          - my_column
  my_db_2:
    tables:
      my_table_3:
        columns:
          - id
          - my_column
```

## Build

need c++ build tools and have mysqldump binary installed in PATH
```shell
go build . -o go-canal
```

Or, you can use docker.
```shell
docker run -v config.yaml:/etc/canal/config.taml shi2002/go-canal:latest
```

## Contributors

This project exists thanks to all the people who contribute.

<a href="https://github.com/OpenTreeHole/go-canal/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=OpenTreeHole/go-canal"  alt="contributors"/>
</a>

## Licence

[![license](https://img.shields.io/github/license/OpenTreeHole/go-canal)](https://github.com/OpenTreeHole/go-canal/blob/master/LICENSE)
© OpenTreeHole