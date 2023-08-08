package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/redis/go-redis/v9"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository interface {
	Get(id uint) (auth models.Auth, err error)
	Create(email string, plainPassword string, audience []string, flags []string) (auth *models.Auth, err error)
	Authenticate(email string, plainPassword string) (auth *models.Auth, found bool)
	InvalidateToken(token *models.UserClaims, rawToken string) error
	IsInvalidatedToken(rawToken string) bool
}

type dbAuthRepository struct {
	context     context.Context
	redisClient *redis.Client
	db          *sql.DB
	ormDB       *gorm.DB
}

func NewDBAuthRepository(db *sql.DB, redisClient *redis.Client) AuthRepository {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return &dbAuthRepository{
		context:     context.Background(),
		redisClient: redisClient,
		db:          db,
		ormDB:       gormDB,
	}
}
func (s *dbAuthRepository) Get(id uint) (auth models.Auth, err error) {
	if tx := s.ormDB.Preload(clause.Associations).First(&auth, id); tx.Error != nil {
		return auth, tx.Error
	} else {
		return auth, nil
	}
}

var ErrUserExists = errors.New("user already exists")

func (s *dbAuthRepository) Create(email string, plainPassword string, audience []string, flags []string) (auth *models.Auth, err error) {
	s.ormDB.Where(&models.Auth{Email: email}).First(&auth)
	if auth != nil && auth.ID > 0 {
		return nil, ErrUserExists
	}

	var domains []models.Domain
	s.ormDB.Where("domain IN ?", audience).Find(&domains)

	var dbFlags []models.Flag
	s.ormDB.Where("name IN ?", flags).Find(&flags)

	salt := createSalt()

	auth = &models.Auth{
		Email:    email,
		Password: calculatePasswordHash(plainPassword, salt),
		Salt:     salt,
		Flags:    dbFlags,
		Domains:  domains,
	}

	result := s.ormDB.Create(auth)
	if result.Error != nil {
		return nil, result.Error
	}
	return auth, nil
}

func (s *dbAuthRepository) Authenticate(email string, plainPassword string) (auth *models.Auth, found bool) {
	s.ormDB.Preload(clause.Associations).Where(&models.Auth{Email: email}).First(&auth)
	if auth.ID == 0 {
		return nil, false
	} else if passwordHash := calculatePasswordHash(plainPassword, auth.Salt); passwordHash != auth.Password {
		return nil, false
	} else {
		return auth, true
	}
}

func (s *dbAuthRepository) InvalidateToken(token *models.UserClaims, rawToken string) error {
	key := s.getRedisInvalidTokenKey(rawToken)
	expiration := time.Now().Sub(token.ExpiresAt.Time)
	err := s.redisClient.Set(s.context, key, true, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *dbAuthRepository) IsInvalidatedToken(rawToken string) bool {
	key := s.getRedisInvalidTokenKey(rawToken)
	val, _ := s.redisClient.Get(s.context, key).Bool()
	return val
}

func (s *dbAuthRepository) getRedisInvalidTokenKey(rawToken string) string {
	return fmt.Sprintf("token.invalid.%s", rawToken)
}

func calculatePasswordHash(plainPassword string, salt string) string {
	h := sha256.New()
	h.Write([]byte(plainPassword))
	h.Write([]byte(salt))
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
}

func createSalt() string {
	salt := helpers.GetEnv("AUTH_SALT", "salt")
	salt = fmt.Sprintf("%s%s", salt, time.Now().String())
	salt = fmt.Sprintf("%x", sha256.Sum256([]byte(salt)))
	return salt
}
