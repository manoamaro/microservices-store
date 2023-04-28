package repositores

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/assert/v2"
	_ "github.com/lib/pq"
	driven2 "github.com/manoamaro/microservices-store/order_service/internal/adapters/driven"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type cartRepositorySuite struct {
	suite.Suite
	DB         *gorm.DB
	repository driven_ports.CartRepository
}

func (s *cartRepositorySuite) SetupSuite() {
	dbName := "order_service_test"
	connInfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
	db, err := sql.Open("postgres", connInfo)
	require.NoError(s.T(), err)

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
	require.NoError(s.T(), err)
	_, err = db.Exec("CREATE DATABASE " + dbName)
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=postgres password=postgres dbname=%s host=127.0.0.1 sslmode=disable", dbName),
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.repository = driven2.NewCartDBRepository(s.DB)
}

func (s *cartRepositorySuite) TestGetOrCreateCart() {
	userId := "USER_ID"

	cart, err := s.repository.GetOrCreateByUserId(userId)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), cart.UserId, userId)
}

func TestCartRepositorySuite(t *testing.T) {
	suite.Run(t, new(cartRepositorySuite))
}
