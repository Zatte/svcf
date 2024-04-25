package subworker_test

import (
	"testing"
	"time"

	"github.com/zatte/svcf/nullworker"
	"github.com/zatte/svcf/subworker"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestSubWorkerNoWorkers(t *testing.T) {
	baseWorker := nullworker.New()
	baseWorkerWithOrch := subworker.NewOrchestrator(baseWorker)

	level := zap.DebugLevel
	logger := zaptest.NewLogger(t, zaptest.Level(level), zaptest.WrapOptions(zap.AddCaller()))
	defer logger.Sync()

	if err := baseWorkerWithOrch.Init(logger); err != nil {
		t.Fatalf("failed to initialize subworker: %v", err)
	}

	wait := time.Millisecond * 50

	tNow := time.Now()
	go func() {
		time.Sleep(wait)
		baseWorkerWithOrch.Terminate()
	}()
	baseWorkerWithOrch.Run()

	if time.Since(tNow) < wait {
		t.Fatalf("worker terminated too early, expected at least %v got %v", wait, time.Since(tNow))
	}
}

func TestSubWorker(t *testing.T) {
	level := zap.DebugLevel
	logger := zaptest.NewLogger(t, zaptest.Level(level), zaptest.WrapOptions(zap.AddCaller()))
	logger = logger.Named("root")
	defer logger.Sync()

	baseWorkerWithOrch := subworker.NewOrchestrator(nullworker.New())

	sw1 := nullworker.New()
	if err := baseWorkerWithOrch.AddSubWorker("subworker_1", sw1); err != nil {
		t.Fatalf("failed to add subworker: %v", err)
	}
	sw2 := nullworker.New()
	if err := baseWorkerWithOrch.AddSubWorker("subworker_2", sw2); err != nil {
		t.Fatalf("failed to add subworker: %v", err)
	}
	sw3 := nullworker.New()
	if err := baseWorkerWithOrch.AddSubWorker("subworker_3", sw3); err != nil {
		t.Fatalf("failed to add subworker: %v", err)
	}

	if err := baseWorkerWithOrch.Init(logger); err != nil {
		t.Fatalf("failed to initialize subworker: %v", err)
	}

	wait := time.Millisecond * 50
	tNow := time.Now()

	go func() {
		time.Sleep(wait)
		if err := baseWorkerWithOrch.Terminate(); err != nil {
			t.Fatalf("failed to terminate subworker: %v", err)
		}
	}()

	err := baseWorkerWithOrch.Run()
	t.Logf("baseWorker.Run() completed after %v: %v", time.Since(tNow), err)

	if time.Since(tNow) < wait {
		t.Fatalf("worker terminated too early")
	}

	if sw1.Ctx().Err() == nil {
		t.Fatalf("subworker 1 not terminated")
	}
	if sw2.Ctx().Err() == nil {
		t.Fatalf("subworker 2 not terminated")
	}
	if sw3.Ctx().Err() == nil {
		t.Fatalf("subworker 3 not terminated")
	}

}
