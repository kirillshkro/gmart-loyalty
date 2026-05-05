package repository

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type OrderTestSuite struct {
	suite.Suite
	repo IOrderRepository
}

func (o *OrderTestSuite) SetupSuite() {
	l := logger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), logger.Config{
		Colorful:             true,
		ParameterizedQueries: false,
		SlowThreshold:        1000 * time.Millisecond,
		LogLevel:             logger.Info,
	})
	opts := &gorm.Config{
		PrepareStmt:    false,
		TranslateError: true,
		Logger:         l,
	}
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), opts)
	o.Require().NoError(err)
	o.repo = NewOrderRepository(db)
}

func (o *OrderTestSuite) TearDownSuite() {

}

func (o *OrderTestSuite) Test_Validate() {
	testCases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Valid number",
			input: "6011111111111117",
			want:  true,
		},
		{
			name:  "Invalid number",
			input: "4532148803436477",
			want:  false,
		},
	}

	for _, tc := range testCases {
		o.T().Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, utils.Valid(tc.input))
		})
	}
}

func Test_OrderSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
