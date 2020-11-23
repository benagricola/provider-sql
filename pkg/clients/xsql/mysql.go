package xsql

import (
	"context"
	"database/sql"
	"fmt"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
)

// A MySQLDB client.
type MySQLDB struct {
	dsn      string
	endpoint string
	port     string
}

// NewMySQLDB returns a new MySQL database client.
func NewMySQLDB(creds map[string][]byte) DB {
	// TODO(negz): Support alternative connection secret formats?
	endpoint := string(creds[runtimev1alpha1.ResourceCredentialsSecretEndpointKey])
	port := string(creds[runtimev1alpha1.ResourceCredentialsSecretPortKey])
	return MySQLDB{
		dsn: fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			creds[runtimev1alpha1.ResourceCredentialsSecretUserKey],
			creds[runtimev1alpha1.ResourceCredentialsSecretPasswordKey],
			endpoint,
			port),
		endpoint: endpoint,
		port:     port,
	}
}

// Exec the supplied query.
func (c MySQLDB) Exec(ctx context.Context, q Query) error {
	d, err := sql.Open("mysql", c.dsn)
	if err != nil {
		return err
	}
	defer d.Close() //nolint:errcheck

	_, err = d.ExecContext(ctx, q.String, q.Parameters...)
	return err
}

// Scan the results of the supplied query into the supplied destination.
func (c MySQLDB) Scan(ctx context.Context, q Query, dest ...interface{}) error {
	db, err := sql.Open("mysql", c.dsn)
	if err != nil {
		return err
	}
	defer db.Close() //nolint:errcheck

	return db.QueryRowContext(ctx, q.String, q.Parameters...).Scan(dest...)
}

// GetConnectionDetails returns the connection details for a user of this DB
func (c MySQLDB) GetConnectionDetails(username, password string) managed.ConnectionDetails {
	return managed.ConnectionDetails{
		runtimev1alpha1.ResourceCredentialsSecretUserKey:     []byte(username),
		runtimev1alpha1.ResourceCredentialsSecretPasswordKey: []byte(password),
		runtimev1alpha1.ResourceCredentialsSecretEndpointKey: []byte(c.endpoint),
		runtimev1alpha1.ResourceCredentialsSecretPortKey:     []byte(c.port),
	}
}