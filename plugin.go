package mysqlplugin

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/suprsokr/vanilladbc/pkg/dbc"
	"github.com/suprsokr/vanilladbc/pkg/dbd"
)

// Config holds MySQL connection configuration
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Plugin implements both plugin.Writer and plugin.Reader interfaces for MySQL
type Plugin struct {
	config     Config
	db         *sql.DB
	tableName  string
	
	// For writing
	columnOrder []string
	versionDef  *dbd.VersionDefinition
	columns     map[string]dbd.ColumnDefinition
	
	// For reading
	rows       *sql.Rows
	columnNames []string
}

// New creates a new MySQL plugin with the given configuration
func New(config Config, tableName string) (*Plugin, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Database)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}
	
	return &Plugin{
		config:    config,
		db:        db,
		tableName: tableName,
	}, nil
}

// WriteHeader is called once before any records are written
func (p *Plugin) WriteHeader(versionDef *dbd.VersionDefinition, columns map[string]dbd.ColumnDefinition) error {
	p.versionDef = versionDef
	p.columns = columns
	
	// Build column order from version definition
	p.columnOrder = make([]string, 0, len(versionDef.Definitions))
	for _, def := range versionDef.Definitions {
		p.columnOrder = append(p.columnOrder, def.Column)
	}
	
	// Create table if it doesn't exist
	createTableSQL := p.buildCreateTableSQL()
	_, err := p.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	
	// Truncate existing data
	truncateSQL := fmt.Sprintf("TRUNCATE TABLE %s", p.tableName)
	_, err = p.db.Exec(truncateSQL)
	if err != nil {
		return fmt.Errorf("failed to truncate table: %w", err)
	}
	
	return nil
}

// WriteRecord is called for each record in the DBC file
func (p *Plugin) WriteRecord(record dbc.Record) error {
	// Build INSERT statement
	placeholders := make([]string, len(p.columnOrder))
	values := make([]interface{}, len(p.columnOrder))
	
	for i, colName := range p.columnOrder {
		placeholders[i] = "?"
		values[i] = record[colName]
	}
	
	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		p.tableName,
		strings.Join(p.columnOrder, ", "),
		strings.Join(placeholders, ", "))
	
	_, err := p.db.Exec(insertSQL, values...)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	
	return nil
}

// WriteFooter is called once after all records are written
func (p *Plugin) WriteFooter() error {
	// Nothing to do for MySQL
	return nil
}

// ReadHeader is called once before reading records
func (p *Plugin) ReadHeader() (*dbd.VersionDefinition, map[string]dbd.ColumnDefinition, error) {
	// Query all columns
	querySQL := fmt.Sprintf("SELECT * FROM %s", p.tableName)
	
	rows, err := p.db.Query(querySQL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query table: %w", err)
	}
	
	p.rows = rows
	
	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get column names: %w", err)
	}
	
	p.columnNames = columnNames
	
	// MySQL doesn't store schema info, so we return what was set
	return p.versionDef, p.columns, nil
}

// ReadRecord is called repeatedly to read records
func (p *Plugin) ReadRecord() (dbc.Record, error) {
	if !p.rows.Next() {
		return nil, nil // No more records
	}
	
	// Create value holders
	values := make([]interface{}, len(p.columnNames))
	valuePtrs := make([]interface{}, len(p.columnNames))
	for i := range values {
		valuePtrs[i] = &values[i]
	}
	
	// Scan row
	if err := p.rows.Scan(valuePtrs...); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	
	// Build record
	record := make(dbc.Record)
	for i, colName := range p.columnNames {
		record[colName] = values[i]
	}
	
	return record, nil
}

// Close is called to cleanup resources
func (p *Plugin) Close() error {
	if p.rows != nil {
		p.rows.Close()
	}
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// SetSchema allows the caller to set the schema information needed for reading
func (p *Plugin) SetSchema(versionDef *dbd.VersionDefinition, columns map[string]dbd.ColumnDefinition) {
	p.versionDef = versionDef
	p.columns = columns
}

// buildCreateTableSQL generates CREATE TABLE SQL from the DBD definition
func (p *Plugin) buildCreateTableSQL() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", p.tableName))
	
	for i, def := range p.versionDef.Definitions {
		if i > 0 {
			sb.WriteString(",\n")
		}
		
		sb.WriteString("  ")
		sb.WriteString(def.Column)
		sb.WriteString(" ")
		
		// Map DBD type to MySQL type
		colDef := p.columns[def.Column]
		switch colDef.Type {
		case dbd.TypeInt:
			if def.Size <= 8 {
				sb.WriteString("TINYINT")
			} else if def.Size <= 16 {
				sb.WriteString("SMALLINT")
			} else if def.Size <= 32 {
				sb.WriteString("INT")
			} else {
				sb.WriteString("BIGINT")
			}
		case dbd.TypeUInt:
			if def.Size <= 8 {
				sb.WriteString("TINYINT UNSIGNED")
			} else if def.Size <= 16 {
				sb.WriteString("SMALLINT UNSIGNED")
			} else if def.Size <= 32 {
				sb.WriteString("INT UNSIGNED")
			} else {
				sb.WriteString("BIGINT UNSIGNED")
			}
		case dbd.TypeFloat:
			sb.WriteString("FLOAT")
		case dbd.TypeString:
			sb.WriteString("TEXT")
		case dbd.TypeLocString:
			sb.WriteString("TEXT") // Store as JSON or concatenated
		default:
			sb.WriteString("TEXT")
		}
		
		// Mark ID columns as primary key
		if def.IsID {
			sb.WriteString(" PRIMARY KEY")
		}
	}
	
	sb.WriteString("\n)")
	
	return sb.String()
}
