package test

import (
	"errors"
	"testing"
	"time"

	"github.com/chuck7000/gotch"
)

func SetupManager(d bool, dtime int, q bool, qtime int, cf gotch.ChangeProcessor) (m *gotch.Manager) {
	m = gotch.NewManager(&gotch.ManagerOptions{
		PreRunDelay:     d,
		PreRunDelayTime: dtime,
		QuietPeriod:     q,
		QuietPeriodTime: qtime,
		ChangeFunc:      cf,
	})
	return
}

//func SetupChangeProcessor(changes int) gotch.ChangeProcessor {
//  return
//}

func TestBasicManager(t *testing.T) {
	wait := make(chan struct{})

	m := SetupManager(true, 1, true, 1, func() error {
		wait <- struct{}{}
		return nil
	})

	go func() {
		m.Change <- struct{}{}
	}()

	select {
	case <-wait:
	case <-time.After(time.Second * 3):
		t.Fail()
	}

}

func TestManagerWithMultipleCallsInRestPeriod(t *testing.T) {
	i := 0
	m := SetupManager(true, 1, false, 0, func() error {
		i++
		return nil
	})

	go func() {
		for l := 0; l < 2; l++ {
			m.Change <- struct{}{}
		}
	}()

	time.Sleep(2 * time.Second)

	if i != 1 {
		t.Fail()
	}
}

func TestManagerWithCallsInQuietTime(t *testing.T) {
	i := 0
	m := SetupManager(false, 0, true, 3, func() error {
		i++
		return nil
	})

	m.Change <- struct{}{}

	time.Sleep(1 * time.Second)

	m.Change <- struct{}{}
	m.Change <- struct{}{}

	time.Sleep(5 * time.Second)

	if i != 2 {
		t.Fail()
	}
}

func TestManagerWithChangeFuncReturningError(t *testing.T) {
	wait := make(chan struct{})

	m := SetupManager(false, 0, false, 0, func() error {
		wait <- struct{}{}
		return errors.New("a test error")
	})

	go func() {
		m.Change <- struct{}{}
	}()

	select {
	case <-wait:
	case <-time.After(time.Second * 3):
		t.Fail()
	}

}
