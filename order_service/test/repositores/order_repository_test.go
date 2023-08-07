package repositores_test

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/assert/v2"
	_ "github.com/lib/pq"
	drivenadapters "github.com/manoamaro/microservices-store/order_service/internal/adapters/driven"
	drivenports "github.com/manoamaro/microservices-store/order_service/internal/core/ports/driven"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type orderESRepositorySuite struct {
	suite.Suite
	DB         *gorm.DB
	repository drivenports.OrderRepository
}

func (s *orderESRepositorySuite) SetupSuite() {
	dbName := "order_service_test"
	connInfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
	db, err := sql.Open("postgres", connInfo)
	require.NoError(s.T(), err)

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\"", dbName))
	require.NoError(s.T(), err)
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=postgres password=postgres dbname=%s host=127.0.0.1 sslmode=disable", dbName),
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.repository, err = drivenadapters.NewOrderESRepository(s.DB)
}

func (s *orderESRepositorySuite) TestGetOrder() {
	orderId := 1

	order, err := s.repository.Get(uint64(orderId))

	require.NoError(s.T(), err)
	assert.Equal(s.T(), order.ID, uint64(orderId))
}

func TestCartRepositorySuite(t *testing.T) {
	suite.Run(t, &orderESRepositorySuite{})
}
