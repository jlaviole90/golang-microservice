package resources

import (
	"context"
	"database/sql"
	"employee-worklog-service/utils"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

var database *sql.DB

// NOTE: This connection is set up for Postgres on Google Cloud.
// It can be modified to fit all needs. Most of the code here is
// only applicable to gcloud.
// 
// Connecting to a gcloud DB on a container running locally can
// be a PITA but it's doable.
func GetDB() (*sql.DB, error) {
    if database != nil {
        return database, nil
    }

    var (
        dbUser = utils.ReqEnvs("DB_USER")
        dbPwd = utils.ReqEnvs("DB_PASS")
        dbName = utils.ReqEnvs("DB_NAME")
        instanceConnectionName = utils.ReqEnvs("INSTANCE_CONNECTION_NAME")
        usePrivate = os.Getenv("PRIVATE_IP")
    )

    dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
    config, err := pgx.ParseConfig(dsn)
    if err != nil {
        return nil, err
    }

    var opts []cloudsqlconn.Option
    if usePrivate != "" {
        opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
    }
    d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
    if err != nil {
        return nil, err
    }
    config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
        return d.Dial(ctx, instanceConnectionName)
    }
    dbURI := stdlib.RegisterConnConfig(config)
    database, err = sql.Open("pgx", dbURI)
    if err != nil {
        return nil, fmt.Errorf("sql.Open: %v", err)
    }
    return database, nil
}
