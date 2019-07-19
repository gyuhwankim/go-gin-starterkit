package auth_test

import (
	"net/http"
	"testing"

	"github.com/gghcode/go-gin-starterkit/app/api/auth"
	"github.com/gghcode/go-gin-starterkit/app/api/testutil"
	"github.com/gghcode/go-gin-starterkit/app/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gghcode/go-gin-starterkit/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerIntegration struct {
	suite.Suite

	ginEngine *gin.Engine
	dbConn    *db.Conn

	testUser user.User
}

func TestOAuth2ControllerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(controllerIntegration))
}

func (suite *controllerIntegration) SetupSuite() {
	gin.SetMode(gin.TestMode)

	conf, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	dbConn, err := db.NewConn(conf)
	require.NoError(suite.T(), err)

	suite.ginEngine = gin.New()
	suite.dbConn = dbConn

	userRepo := user.NewRepository(dbConn)
	passport := service.NewPassport()

	authController := auth.NewController(conf, userRepo, passport)
	authController.RegisterRoutes(suite.ginEngine)

	testUser := user.User{
		UserName: "username",
	}
	testUser.PasswordHash, _ = passport.HashPassword("password")

	suite.testUser, err = userRepo.CreateUser(testUser)
	suite.NoError(err)
}

func (suite *controllerIntegration) TestGetToken() {
	testCases := []struct {
		description    string
		reqPayload     *auth.CreateAccessTokenRequest
		expectedStatus int
	}{
		{
			description: "ShouldBeTokenGenerated",
			reqPayload: &auth.CreateAccessTokenRequest{
				UserName: suite.testUser.UserName,
				Password: "password",
			},
			expectedStatus: http.StatusOK,
		},
		{
			description:    "ShouldBeBadReqest",
			reqPayload:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "ShouldBeUnauthorized",
			reqPayload: &auth.CreateAccessTokenRequest{
				UserName: "NOT_EXISTS_USER_NAME",
				Password: "password",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			description: "ShouldBeUnauthorized_WhenIncorrectPassword",
			reqPayload: &auth.CreateAccessTokenRequest{
				UserName: suite.testUser.UserName,
				Password: "INCORRECT_PASSWORD",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			reqBody := testutil.ReqBodyFromInterface(suite.T(), tc.reqPayload)

			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"POST",
				"/token",
				reqBody,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)
		})
	}
}