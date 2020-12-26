package core

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rooster/rsql"
)

type RockDB struct {
	namespace string
	Name      string
	*rsql.DB
}

func (d *RockDB) GetNamespace() string {
	if d.namespace == "" {
		return meta.DefaultNamespace
	}
	return d.namespace
}

func NewDB(ds *meta.DataSource) (*RockDB, error) {
	var driverName, dsn string
	if strings.ToLower(ds.Driver) == "mysql" {
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", ds.User, ds.Password, ds.Host, ds.Port, ds.Database)
	}
	db, err := rsql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &RockDB{
		namespace: ds.Namespace,
		Name:      ds.Name,
		DB:        db,
	}, nil
}
