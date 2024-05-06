package sync

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

func TestSynchronizeReadAndWriteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQLRepository := repository.NewMockCurrencyRepository(ctrl)
	mockMemcachedRepository := repository.NewMockCurrencyRepository(ctrl)
	mockLogger := logrus.New()

	// Mock data
	firstID := uuid.New()
	secondID := uuid.New()
	payments := []*model.CurrencyRepositoryModel{
		{ID: firstID},
		{ID: secondID},
	}

	// Mock expectations
	mockSQLRepository.EXPECT().FindAll(gomock.Any()).Return(payments, nil).AnyTimes()
	mockMemcachedRepository.EXPECT().BulkUpsert(gomock.Any(), payments).Return(payments, nil).AnyTimes()

	// Execute the function
	go SynchronizeReadAndWriteData(mockMemcachedRepository, mockSQLRepository, 1, mockLogger)

	// Sleep for a while to allow the worker to finish
	time.Sleep(3 * time.Second)
}
