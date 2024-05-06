package worker

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func RunFuncEveryGivenPeriod(
	fn func(ctx context.Context) error,
	periodInSecond int,
	log *logrus.Logger,
) {
	for {
		ctx := context.WithoutCancel(context.Background())
		if err := fn(ctx); err != nil {
			log.Warnf("failed to run %s: %v", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), err)
		}

		time.Sleep(time.Duration(periodInSecond * int(time.Second)))
	}
}
