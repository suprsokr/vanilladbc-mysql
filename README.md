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

### Environment Variables (Recommended)

For security, use environment variables instead of command-line flags:

```bash
# Set credentials in environment
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=secret
export MYSQL_DATABASE=wow_vanilla

# Convert DBC to MySQL table
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-table spell

# Import MySQL table back to DBC
vanilladbc import Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-table spell \
  --output Spell.dbc
```

### With Command-Line Flags

```bash
# Convert DBC to MySQL table (not recommended - password visible in process list)
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-host localhost \
  --mysql-port 3306 \
  --mysql-user root \
  --mysql-password secret \
  --mysql-database wow_vanilla \
  --mysql-table spell
```

### With Config File

```bash
# Create config.json (add to .gitignore!)
cat > mysql-config.json << EOF
{
  "host": "localhost",
  "port": 3306,
  "user": "wow_user",
  "password": "secure_password",
  "database": "wow_vanilla"
}
EOF

# Convert using config file
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-config mysql-config.json \
  --mysql-table spell
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
    // Option 1: Load from environment variables (recommended)
    plugin, err := mysqlplugin.NewFromEnv("spell")
    if err != nil {
        panic(err)
    }
    defer plugin.Close()
    
    // Option 2: Load from JSON config file
    plugin2, err := mysqlplugin.NewFromFile("mysql-config.json", "spell")
    if err != nil {
        panic(err)
    }
    defer plugin2.Close()
    
    // Option 3: Configure manually
    config := mysqlplugin.Config{
        Host:     "localhost",
        Port:     3306,
        User:     "root",
        Password: "secret",
        Database: "wow_vanilla",
    }
    plugin3, err := mysqlplugin.New(config, "spell")
    if err != nil {
        panic(err)
    }
    defer plugin3.Close()
    
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

## Security Best Practices

### Environment Variables

**Always use environment variables for credentials** instead of command-line flags:

```bash
# Good - credentials not visible in process list
export MYSQL_PASSWORD=secret
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --plugin mysql

# Bad - password visible in ps aux, shell history
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --mysql-password secret
```

### Supported Environment Variables

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `MYSQL_HOST` | localhost | No | MySQL server hostname |
| `MYSQL_PORT` | 3306 | No | MySQL server port |
| `MYSQL_USER` | root | No | MySQL username |
| `MYSQL_PASSWORD` | - | **Yes** | MySQL password |
| `MYSQL_DATABASE` | - | **Yes** | Database name |

### .env File Support

Create a `.env` file (add to `.gitignore`!):

```bash
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=wow_user
MYSQL_PASSWORD=secure_password_here
MYSQL_DATABASE=wow_vanilla
```

Load it before running commands:

```bash
# Load .env file
export $(cat .env | xargs)

# Run conversion
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --plugin mysql --mysql-table spell
```

### JSON Config File Format

Create a `mysql-config.json` file (**add to `.gitignore`!**):

```json
{
  "host": "localhost",
  "port": 3306,
  "user": "wow_user",
  "password": "secure_password",
  "database": "wow_vanilla"
}
```

**Security Note:** The config file is created with `0600` permissions (owner read/write only) to protect credentials.

Use the config file:

```bash
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 \
  --plugin mysql \
  --mysql-config mysql-config.json \
  --mysql-table spell
```

### Configuration Priority

The plugin loads configuration in this order (first found wins):

1. **Config file** (if `--mysql-config` flag provided)
2. **Environment variables**
3. **Command-line flags**
4. **Defaults** (host=localhost, port=3306, user=root)

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
