package gorm_repositories

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type RepositorySuite struct {
	suite.Suite
	SqlDB   *sql.DB
	DB      *gorm.DB
	SqlMock sqlmock.Sqlmock
	repo    GORMRepository
}

func (s *RepositorySuite) SetupTest() {
	var err error
	s.SqlDB, s.SqlMock, err = sqlmock.New()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.SqlDB)
	require.NotNil(s.T(), s.SqlMock)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 s.SqlDB,
		PreferSimpleProtocol: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(s.T(), err)
	s.repo = New(s.DB)
}

func (s *RepositorySuite) TearDownTest() {
	s.SqlDB.Close()
}

func (s *RepositorySuite) SQL(sql string, a ...interface{}) string {
	return strings.Join(strings.Fields(fmt.Sprintf(sql, a...)), " ")
}

func MockRows(objs ...interface{}) *sqlmock.Rows {
	if reflect.ValueOf(objs[0]).Kind() == reflect.Slice {
		if len(objs) > 1 {
			panic(fmt.Errorf("objs must have one element if first element is slice"))
		}

		s := reflect.ValueOf(objs[0])
		if s.IsNil() {
			return nil
		}

		objs = make([]interface{}, s.Len())

		for i := 0; i < s.Len(); i++ {
			objs[i] = s.Index(i).Interface()
		}
	}

	s, err := schema.Parse(objs[0], &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic("failed to create schema")
	}

	columns := make([]string, 0)
	for _, field := range s.Fields {
		if len(field.DBName) == 0 {
			continue
		}
		columns = append(columns, field.DBName)
	}

	rows := sqlmock.NewRows(columns)

	for _, obj := range objs {
		row := make([]driver.Value, 0)

		for _, field := range s.Fields {
			if len(field.DBName) == 0 {
				continue
			}
			r := reflect.ValueOf(obj)
			f := reflect.Indirect(r).FieldByName(field.Name)
			row = append(row, f.Interface())
		}

		rows = rows.AddRow(row...)
	}

	return rows
}
