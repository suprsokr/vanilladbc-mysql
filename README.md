# vanilladbc-mysql

MySQL conversion plugin for [vanilladbc-cli](https://github.com/suprsokr/vanilladbc-cli).

Provides bidirectional conversion between World of Warcraft Vanilla DBC files and MySQL database tables.

## Features

- **DBC to MySQL** - Export DBC files to MySQL tables
- **MySQL to DBC** - Import MySQL tables back to DBC format
- **Auto Table Creation** - Automatically creates tables with correct column types
- **Configurable Connection** - Supports custom MySQL host, port, user, password, and database
- **Type Mapping** - Intelligent mapping between DBC types and MySQL column types
- **Primary Keys** - Automatically sets ID columns as primary keys

## Installation

```bash
go get github.com/suprsokr/vanilladbc-mysql
```

## Usage

### With vanilladbc-cli

```bash
# Convert DBC to MySQL table
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-host localhost \
  --mysql-port 3306 \
  --mysql-user root \
  --mysql-password secret \
  --mysql-database wow_vanilla \
  --mysql-table spell

# Import MySQL table back to DBC
vanilladbc import Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-host localhost \
  --mysql-port 3306 \
  --mysql-user root \
  --mysql-password secret \
  --mysql-database wow_vanilla \
  --mysql-table spell \
  --output Spell.dbc
```

### As a Library

```go
package main

import (
    mysqlplugin "github.com/suprsokr/vanilladbc-mysql"
    "github.com/suprsokr/vanilladbc/pkg/dbd"
    "github.com/suprsokr/vanilladbc/pkg/dbc"
)

func main() {
    // Configure MySQL connection
    config := mysqlplugin.Config{
        Host:     "localhost",
        Port:     3306,
        User:     "root",
        Password: "secret",
        Database: "wow_vanilla",
    }
    
    // Create plugin
    plugin, err := mysqlplugin.New(config, "spell")
    if err != nil {
        panic(err)
    }
    defer plugin.Close()
    
    // ... use plugin.WriteHeader(), WriteRecord(), WriteFooter()
    // ... or plugin.ReadHeader(), ReadRecord()
}
```

## Type Mapping

The plugin automatically maps DBC types to appropriate MySQL column types:

| DBC Type | Size | MySQL Type |
|----------|------|------------|
| int | 8-bit | TINYINT |
| int | 16-bit | SMALLINT |
| int | 32-bit | INT |
| int | 64-bit | BIGINT |
| uint | 8-bit | TINYINT UNSIGNED |
| uint | 16-bit | SMALLINT UNSIGNED |
| uint | 32-bit | INT UNSIGNED |
| uint | 64-bit | BIGINT UNSIGNED |
| float | 32-bit | FLOAT |
| string | - | TEXT |
| locstring | - | TEXT |

## Requirements

- MySQL 5.7+ or MariaDB 10.2+
- Valid MySQL user with CREATE TABLE and INSERT/SELECT permissions

## Dependencies

- [vanilladbc](https://github.com/suprsokr/vanilladbc) - Core DBC/DBD library
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver

## Related Projects

- [vanilladbc-cli](https://github.com/suprsokr/vanilladbc-cli) - Command-line tool
- [vanilladbc-json](https://github.com/suprsokr/vanilladbc-json) - JSON plugin
- [vanilladbc-csv](https://github.com/suprsokr/vanilladbc-csv) - CSV plugin
- [VanillaDBDefs](https://github.com/suprsokr/VanillaDBDefs) - Database definitions

## License

MIT License - See LICENSE file for details.
