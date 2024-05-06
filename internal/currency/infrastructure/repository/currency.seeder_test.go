package repository

import (
	"errors"
	"fmt"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	"github.com/angryronald/currency-service/lib/cast"
	"github.com/angryronald/currency-service/lib/file"
)

func getFilepath(filename string) string {
	projectDir, _ := file.GetProjectDir("currency-service")

	return fmt.Sprintf("%s/seed/%s", projectDir, filename)
}

func getMockData(filepath string) []*model.CurrencyRepositoryModel {
	data, err := os.ReadFile(filepath)
	if err != nil {
		logrus.Errorf("error reading file: %v", err)
		return nil
	}

	var jsonData map[string]string
	if err = cast.FromBytes(data, &jsonData); err != nil {
		logrus.Errorf("error decoding JSON: %v", err)
		return nil
	}

	currencies := []*model.CurrencyRepositoryModel{}
	for code, name := range jsonData {
		currencies = append(currencies, &model.CurrencyRepositoryModel{
			Code: code,
			Name: name,
		})
	}

	return currencies
}

// A successful test case with valid data.
func TestRunMigration_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock a successful database repository.
	databaseRepo := NewMockCurrencyRepository(ctrl)

	// Mock a successful memcached repository.
	memcachedRepo := NewMockCurrencyRepository(ctrl)

	filename := "currency_code.json"
	filepath := getFilepath(filename)
	mockData := getMockData(filepath)

	databaseRepo.EXPECT().FindAll(gomock.Any()).Return(nil, constant.ErrNotFound).Times(1)
	databaseRepo.EXPECT().BulkUpsert(gomock.Any(), gomock.Any()).Return(mockData, nil).Times(1)

	memcachedRepo.EXPECT().FindAll(gomock.Any()).Return(nil, constant.ErrNotFound).Times(1)
	memcachedRepo.EXPECT().BulkUpsert(gomock.Any(), gomock.Any()).Return(mockData, nil).Times(1)

	err := Seeding(memcachedRepo, databaseRepo, filepath)

	// Assert that there are no errors.
	if err != nil {
		t.Errorf("Expected success, but got an error: %v", err)
	}
}

// A test case for a missing JSON file.
func TestRunMigration_MissingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock a successful database repository.
	databaseRepo := NewMockCurrencyRepository(ctrl)

	// Mock a successful memcached repository.
	memcachedRepo := NewMockCurrencyRepository(ctrl)

	err := Seeding(memcachedRepo, databaseRepo, "incorrect_path.json")

	// Assert that an error occurred.
	if err == nil {
		t.Error("Expected an error for missing file, but got nil")
	}
}

// A test case for a JSON decoding error.
func TestRunMigration_JSONDecodingError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock a successful database repository.
	databaseRepo := NewMockCurrencyRepository(ctrl)

	// Mock a successful memcached repository.
	memcachedRepo := NewMockCurrencyRepository(ctrl)

	filename := "invalid.json"
	filepath := getFilepath(filename)

	err := Seeding(memcachedRepo, databaseRepo, filepath)

	// Assert that an error occurred during JSON decoding.
	if err == nil {
		t.Error("Expected an error for JSON decoding, but got nil")
	}
}

// A test case for a database repository error.
func TestRunMigration_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock a successful database repository.
	databaseRepo := NewMockCurrencyRepository(ctrl)

	// Mock a successful memcached repository.
	memcachedRepo := NewMockCurrencyRepository(ctrl)

	filename := "currency_code.json"
	filepath := getFilepath(filename)

	mockerr := errors.New("something went wrong")
	databaseRepo.EXPECT().FindAll(gomock.Any()).Return(nil, constant.ErrNotFound).Times(1)
	databaseRepo.EXPECT().BulkUpsert(gomock.Any(), gomock.Any()).Return(nil, mockerr).Times(1)

	err := Seeding(memcachedRepo, databaseRepo, filepath)

	// Assert that an error occurred during database insertion.
	if err == nil {
		t.Error("Expected an error for database insertion, but got nil")
	}
}

// A test case for a memcached repository error.
func TestRunMigration_MemcachedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock a successful database repository.
	databaseRepo := NewMockCurrencyRepository(ctrl)

	// Mock a successful memcached repository.
	memcachedRepo := NewMockCurrencyRepository(ctrl)

	filename := "currency_code.json"
	filepath := getFilepath(filename)
	mockData := getMockData(filepath)

	databaseRepo.EXPECT().FindAll(gomock.Any()).Return(nil, constant.ErrNotFound).Times(1)
	databaseRepo.EXPECT().BulkUpsert(gomock.Any(), gomock.Any()).Return(mockData, nil).Times(1)

	mockerr := errors.New("something went wrong")
	memcachedRepo.EXPECT().FindAll(gomock.Any()).Return(nil, constant.ErrNotFound).Times(1)
	memcachedRepo.EXPECT().BulkUpsert(gomock.Any(), gomock.Any()).Return(nil, mockerr).Times(1)

	err := Seeding(memcachedRepo, databaseRepo, filepath)

	// Assert that an error occurred during memcached insertion.
	if err == nil {
		t.Error("Expected an error for memcached insertion, but got nil")
	}
}
