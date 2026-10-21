package forwarder

import (
	"github.com/sirupsen/logrus"
	"io"
)

func quietLogger() *logrus.Logger { log := logrus.New(); log.SetOutput(io.Discard); return log }
