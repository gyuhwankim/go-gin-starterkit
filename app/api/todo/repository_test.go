package todo

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gyuhwankim/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type repoTestSuite struct {
	suite.Suite

	mockDB     *sql.DB
	mockSQL    sqlmock.Sqlmock
	mockGormDB *gorm.DB
	mockDbConn *db.Conn

	repo *Repository
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (suite *repoTestSuite) SetupTest() {
	mockDB, mockSQL, err := sqlmock.New()

	require.NoError(suite.T(), err)

	mockGormDB, err := gorm.Open("postgres", mockDB)

	require.NoError(suite.T(), err)

	suite.mockDB = mockDB
	suite.mockSQL = mockSQL
	suite.mockGormDB = mockGormDB
	suite.mockDbConn = db.NewConn(mockGormDB)
	suite.repo = NewRepository(suite.mockDbConn)

	require.NotNil(suite.T(), suite.repo)
}

func (suite *repoTestSuite) TearDownTest() {
	suite.mockDB.Close()
	suite.mockGormDB.Close()
}

func (suite *repoTestSuite) TestShouldGetTodos() {
	expectedTodos := []TodoModel{
		TodoModel{
			ID:       uuid.NewV4(),
			Title:    "FIRST TITLE",
			Contents: "FIRST CONTENTS",
		},
		TodoModel{
			ID:       uuid.NewV4(),
			Title:    "SECOND TITLE",
			Contents: "SECOND CONTENTS",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "contents"})
	for _, todo := range expectedTodos {
		rows.AddRow(todo.ID, todo.Title, todo.Contents)
	}

	suite.mockSQL.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "todo_models"`)).
		WillReturnRows(rows)

	actualTodos, err := suite.repo.getTodos()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(expectedTodos), len(actualTodos))
	assert.Equal(suite.T(), expectedTodos, actualTodos)
}
