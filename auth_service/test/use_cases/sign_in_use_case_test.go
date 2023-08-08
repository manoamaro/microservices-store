package use_cases

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/manoamaro/microservices-store/auth_service/test/mocks"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type SignInUseCaseTestSuite struct {
	suite.Suite
	authRepository *mocks.AuthRepository
	useCase        use_cases.SignInUseCase
}

func TestSignInUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(SignInUseCaseTestSuite))
}

func (suite *SignInUseCaseTestSuite) SetupTest() {
	suite.authRepository = new(mocks.AuthRepository)
	suite.useCase = use_cases.NewSignInUseCase(suite.authRepository)
}

func (suite *SignInUseCaseTestSuite) TestSignInWithValidCredentials() {
	suite.authRepository.On("Authenticate", validAuth.Email, validAuth.Password).Return(&validAuth, true).Once()
	args := use_cases.SignInDTO{
		Email:         validAuth.Email,
		PlainPassword: validAuth.Password,
	}
	var domains jwt.ClaimStrings = validAuth.DomainArray()

	resultDTO, err := suite.useCase.SignIn(args)
	suite.NoError(err)
	suite.NotEmpty(resultDTO.Token)
	suite.NotEmpty(resultDTO.RefreshToken)
	userClaims, err := helpers.GetClaimsFromToken(resultDTO.Token)
	suite.NoError(err)
	suite.Equal(strconv.Itoa(int(validAuth.ID)), userClaims.ID)
	suite.Equal(domains, userClaims.Audience)
	suite.Equal(validAuth.FlagsArray(), userClaims.Flags)
	suite.authRepository.AssertExpectations(suite.T())
}

func (suite *SignInUseCaseTestSuite) TestSignInWithInvalidCredentials() {
	suite.authRepository.On("Authenticate", "invalidEmail@example.com", "password").Return(nil, false).Once()
	args := use_cases.SignInDTO{
		Email:         "invalidEmail@example.com",
		PlainPassword: "password",
	}

	_, err := suite.useCase.SignIn(args)
	suite.Error(err, use_cases.ErrUserNotFound)
	suite.authRepository.AssertExpectations(suite.T())
}
