package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
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

func (s *dbAuthRepository) Create(email string, plainPassword string, audience []string, flags []string) (auth *models.Auth, err error) {
	s.ormDB.Where(&models.Auth{Email: email}).First(&auth)
	if auth != nil && auth.ID > 0 {
		return nil, errors.New("user already exists")
	}

	var domains []models.Domain
	s.ormDB.Where("domain IN ?", audience).Find(&domains)

	var dbFlags []models.Flag
	s.ormDB.Where("name IN ?", flags).Find(&flags)

	salt := strconv.FormatInt(time.Now().UnixNano(), 16)
	auth = &models.Auth{
		Email:    email,
		Password: CalculatePasswordHash(plainPassword, salt),
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
	} else if passwordHash := CalculatePasswordHash(plainPassword, auth.Salt); passwordHash != auth.Password {
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

func CalculatePasswordHash(plainPassword string, salt string) string {
	h := sha256.New()
	h.Write([]byte(plainPassword))
	h.Write([]byte(salt))
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
}
