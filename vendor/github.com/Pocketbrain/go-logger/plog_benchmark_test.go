package plog

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkLogging(b *testing.B) {
	b.Logf("Logging benchmarking test")

	b.Run("global", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Info("Info test logging, will use realistic message length")
			}
		})
	})

	b.Run("instantiated", func(b *testing.B) {
		l := New()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info("Info test logging, will use realistic message length")
			}
		})
	})

	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := logrus.New()
		logger.Level = logrus.DebugLevel
		logger.Out = ioutil.Discard
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("Info test logging, will use realistic message length")
			}
		})
	})
}
