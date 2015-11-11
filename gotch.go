package gotch

import (
	"time"

	"github.com/chuck7000/gotch/log"
)

// ChangeProcessor any function fitting this interface can be passed as a callback to a Manager
type ChangeProcessor func() (err error)

// Manager is used to trigger a chance processor with a set amount of action buffering
type Manager struct {

	// state applys a mutex to the various manager states to keep sync'ed things synced
	State

	//Change channel used to inform the manager that somethign changed
	Change chan struct{}

	// PreRunDelay boolean indicating whether the manager should wait for more changes before fireing
	// the supplied ChangeFunc
	PreRunDelay bool

	// PreRunDelayTime is the number of seconds the manager should wait before sarting a run
	// of ChangeFunc
	PreRunDelayTime int

	// QuietPeriod is a boolean indicating whether the manager shoudl enforce a quite period
	// after completing a run of ChangeFunc
	QuietPeriod bool

	// QuietPeriodTime is how long the quite period should last in seconds
	QuietPeriodTime int

	// ChangeFunc is the ChangeProcessor function passed in when the manager is created
	ChangeFunc ChangeProcessor
}

// ManagerOptions is used to set the options for configuring a new Manager
type ManagerOptions struct {
	PreRunDelay     bool
	PreRunDelayTime int
	QuietPeriod     bool
	QuietPeriodTime int
	ChangeFunc      ChangeProcessor
}

// NewManager is used to create and fully initalize a Manager
func NewManager(mo *ManagerOptions) (m *Manager) {
	m = &Manager{
		Change:          make(chan struct{}),
		PreRunDelay:     mo.PreRunDelay,
		PreRunDelayTime: mo.PreRunDelayTime,
		QuietPeriod:     mo.QuietPeriod,
		QuietPeriodTime: mo.QuietPeriodTime,
		ChangeFunc:      mo.ChangeFunc,
	}

	m.initAndListen()

	return
}

// initAndListen will start a go func to listen on the Manager.Change channel and update the
// ChangeHappened state anytime a chance occurs.  It will then check to see if the ChangeFunc
// can be ran, and if so - will run it.
func (m *Manager) initAndListen() {
	go func() {
		log.Info("starting listener...")
		for {
			// listen to the change channel
			<-m.Change
			log.Debug("recieved message from change channel")

			// set changed happened when a message is recieved
			log.Debug("setting change happend to true")
			m.SetChangeHappened(true)

			// check to see if we can run a cycle of the ChangeFunc
			if m.canRunChangeFunc() {
				log.Debug("can run change func, so calling as go routine")
				// if so, launch it in a goroutine.
				go m.runChangeFunc(false)
			}
		}
	}()
}

func (m *Manager) canRunChangeFunc() bool {
	// the ChangeFunc can run if we are not in a quite period, an dif we are not currently running
	if m.Quiet() {
		log.Debug("ChangeFunc can't run due to quiet period...")
	}
	if m.Running() {
		log.Debug("ChangeFunc is already running")
	}
	return !m.Quiet() && !m.Running()
}

// runChangeFunc will manage the state of the manager and enforce PreRunDelay and QuietPeriod
func (m *Manager) runChangeFunc(skipQuiet bool) {
	// mark that the manger is currently running to block re-fireing the run based on
	// changes that come in during a run.
	m.SetRunning(true)
	log.Debug("Setting running status to true")

	// if we have a PreRunDelay, sleep for that period of time before firing the ChangeFunc
	if m.PreRunDelay {
		log.Debugf("pre-run flag is set to true, pausing for %v seconds to wait for more changes",
			m.PreRunDelay)
		time.Sleep(time.Duration(m.PreRunDelayTime) * time.Second)
	}

	// run the change func and log the result if there is an error
	err := m.ChangeFunc()
	if err != nil {
		log.Infof("Running change func resulted in error: %v", err)
	}

	// clear the runnign flag
	m.SetRunning(false)
	log.Debug("Clearing running state")

	// if we need to have a quiet period, and the skipQuiet period override is not set
	// run the quiet period
	if m.QuietPeriod && !skipQuiet {
		// set that we are in a quiet period
		m.SetQuiet(true)
		log.Debugf("entering quiet period of %v", m.QuietPeriodTime)
		// sleep for the amount of time in QuietPeriodTime
		time.Sleep(time.Duration(m.QuietPeriodTime) * time.Second)
		// clear the quiet state
		m.SetQuiet(false)
		log.Debug("exiting quiet period")

		// check to see if more changes happened while we were sleeping
		if m.ChangeHappened() {
			log.Info("More changes happened during process run, so running again")
			// if so, run the change func with the quiet period delayed.
			m.runChangeFunc(true)
		}
	}
}
