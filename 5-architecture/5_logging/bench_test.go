package main

import (
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// go test -bench=. -benchmem

var (
	benchMethod = "GET"
	benchRemote = "127.0.0.1"
	benchPath   = "/users/42"
	benchWork   = 12 * time.Microsecond
)

// --- fmt / stdlib log ---

func BenchmarkFmtPrintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "[%s] %s, %s %s\n",
			benchMethod, benchRemote, benchPath, benchWork)
	}
}

func BenchmarkStdLog(b *testing.B) {
	l := log.New(io.Discard, "", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Printf("[%s] %s, %s %s\n",
			benchMethod, benchRemote, benchPath, benchWork)
	}
}

// --- logrus ---

func BenchmarkLogrusFields(b *testing.B) {
	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithFields(logrus.Fields{
			"method":      benchMethod,
			"remote_addr": benchRemote,
			"work_time":   benchWork,
		}).Info(benchPath)
	}
}

func BenchmarkLogrusNoFields(b *testing.B) {
	l := logrus.New()
	l.SetOutput(io.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info(benchPath)
	}
}

// --- zap: поля в каждом вызове ---

func BenchmarkZapTyped(b *testing.B) {
	l := zap.NewNop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info(benchPath,
			zap.String("method", benchMethod),
			zap.String("remote_addr", benchRemote),
			zap.Duration("work_time", benchWork),
		)
	}
}

// --- zap: поля приклеены один раз через With (контекстный логер) ---

func BenchmarkZapTypedWith(b *testing.B) {
	l := zap.NewNop().With(
		zap.String("method", benchMethod),
		zap.String("remote_addr", benchRemote),
		zap.Duration("work_time", benchWork),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info(benchPath)
	}
}

// --- zap: синтаксический сахар ---

func BenchmarkZapSugar(b *testing.B) {
	l := zap.NewNop().Sugar()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Infow(benchPath,
			"method", benchMethod,
			"remote_addr", benchRemote,
			"work_time", benchWork,
		)
	}
}

// --- zap: реальное кодирование в JSON (в discard) ---

func BenchmarkZapTypedJSON(b *testing.B) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)
	l := zap.New(core)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info(benchPath,
			zap.String("method", benchMethod),
			zap.String("remote_addr", benchRemote),
			zap.Duration("work_time", benchWork),
		)
	}
}
