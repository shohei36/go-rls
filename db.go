package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type TenantDB interface {
	Conn(ctx context.Context, tenantID string) (*sql.Conn, error)
}

type TenantDBInput struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	SchemaName string
}

func NewTenantDB(i TenantDBInput) (TenantDB, error) {
	db, err := openDB(
		i.DBName, i.Port, i.User, i.Password, i.DBName, i.SchemaName,
	)
	if err != nil {
		return nil, err
	}
	return &tenantDB{
		db: db,
	}, nil
}

type tenantDB struct {
	db *sql.DB
}

func (td *tenantDB) Conn(ctx context.Context, tenantID string) (*sql.Conn, error) {
	conn, err := td.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db connection for tenant-id %s: %w", tenantID, err)
	}
	_, err = conn.ExecContext(ctx, fmt.Sprintf("set app.current_tenant_id='%s'", tenantID))
	if err != nil {
		return nil, fmt.Errorf("failed to set tenant-id to env: %s, %w", tenantID, err)
	}
	return conn, nil
}

func openDB(host, port, user, password, dbName, schemaName string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host, port, user, password, dbName, schemaName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	return db, nil
}
