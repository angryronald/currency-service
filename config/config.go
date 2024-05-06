package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	return log
}

func GetValue(key string) string {
	return os.Getenv(key)
}

func GetAllowedClients() map[string]string {
	allowedClientsInString := GetValue(ALLOWED_CLIENTS)
	var allowedClients map[string]string

	if err := json.Unmarshal([]byte(allowedClientsInString), &allowedClients); err != nil {
		panic(fmt.Sprintf("error while fetching env vars: %v (stack trace %v)", err, errors.WithStack(err)))
	}
	return allowedClients
}
