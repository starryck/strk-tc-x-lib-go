package gbspvs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gblog"
)

var mSupervisor *Supervisor

type Supervisor struct {
	daemons       []*Daemon
	exitCode      int
	waitGroup     *sync.WaitGroup
	waitChannel   chan struct{}
	rootContext   context.Context
	rootCanceller context.CancelFunc
	signalChannel chan os.Signal

	gracefulTimeout   time.Duration
	heartbeatInterval time.Duration
}

type Process interface {
	Setup() error
	Start(ctx context.Context) error
}

type Operate func(args ...any)

func GetSupervisor(options *Options) *Supervisor {
	if mSupervisor == nil {
		mSupervisor = newSupervisor(options)
	}
	return mSupervisor
}

func newSupervisor(options *Options) *Supervisor {
	supervisor := (&builder{options: options}).
		initialize().
		setExitCode().
		setWaitGroup().
		setRootContext().
		setSignalChannel().
		setGracefulTimeout().
		setHeartbeatInterval().
		build()
	return supervisor
}

func GetWaitGroup() *sync.WaitGroup {
	if mSupervisor == nil {
		panic("Supervisor hasn't been created.")
	}
	return mSupervisor.waitGroup
}

func GetRootContext() context.Context {
	if mSupervisor == nil {
		panic("Supervisor hasn't been created.")
	}
	return mSupervisor.rootContext
}

func WithWaitGroup(operate Operate, args ...any) {
	wg := GetWaitGroup()
	wg.Add(1)
	defer wg.Done()
	operate(args...)
	return
}

func (supervisor *Supervisor) Handle(process Process) {
	supervisor.daemons = append(supervisor.daemons, &Daemon{
		process:  process,
		typeName: reflect.TypeOf(process).String(),
		isActive: true,
	})
}

func (supervisor *Supervisor) RunForever() {
	supervisor.setupDaemons()
	supervisor.startDaemons()
	supervisor.serveDaemons()
}

func (supervisor *Supervisor) setupDaemons() {
	for _, daemon := range supervisor.daemons {
		daemon.setup()
	}
}

func (supervisor *Supervisor) startDaemons() {
	for _, daemon := range supervisor.daemons {
		go func(daemon *Daemon) {
			supervisor.waitGroup.Add(1)
			defer supervisor.waitGroup.Done()
			defer supervisor.rootCanceller()
			daemon.start(supervisor.rootContext)
			daemon.stamp()
		}(daemon)
	}
}

func (supervisor *Supervisor) serveDaemons() {
	timer := time.NewTicker(supervisor.heartbeatInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			supervisor.emitBeatInfo()
		case <-supervisor.signalChannel:
			supervisor.rootCanceller()
		case <-supervisor.rootContext.Done():
			supervisor.emitShutInfo()
			supervisor.waitDaemons()
			supervisor.emitExitInfo()
			return
		}
	}
}

func (supervisor *Supervisor) emitBeatInfo() {
	fields := supervisor.makeLoggerFields()
	gblog.WithFields(fields).Info("Supervisor heartbeats.")
}

func (supervisor *Supervisor) emitShutInfo() {
	fields := supervisor.makeLoggerFields()
	gblog.WithFields(fields).Infof("Supervisor will gracefully shut down in %v.", supervisor.gracefulTimeout)
}

func (supervisor *Supervisor) waitDaemons() {
	timer, timeout := time.NewTicker(1*time.Second), time.After(supervisor.gracefulTimeout)
	go func() {
		supervisor.waitGroup.Wait()
		supervisor.waitChannel <- struct{}{}
	}()
	for {
		select {
		case <-timer.C:
			supervisor.emitBeatInfo()
		case <-timeout:
			supervisor.exitCode = exitCodeFailure
			return
		case <-supervisor.waitChannel:
			supervisor.exitCode = exitCodeSuccess
			return
		}
	}
}

func (supervisor *Supervisor) emitExitInfo() {
	fields := supervisor.makeLoggerFields()
	switch supervisor.exitCode {
	case exitCodeSuccess:
		gblog.WithFields(fields).Info("Supervisor gracefully shut down normally.")
	case exitCodeFailure:
		gblog.WithFields(fields).Warn("Supervisor gracefully shut down abnormally.")
	}
}

func (supervisor *Supervisor) makeLoggerFields() gblog.Fields {
	actives := []string{}
	inactives := []string{}
	for _, daemon := range supervisor.daemons {
		if daemon.isActive {
			actives = append(actives, daemon.String())
		} else {
			inactives = append(inactives, daemon.String())
		}
	}
	return gblog.Fields{
		"actives":   actives,
		"inactives": inactives,
	}
}

const (
	exitCodeDefault = -1 + iota
	exitCodeSuccess
	exitCodeFailure

	defaultGracefulTimeout   = 30 * time.Second
	defaultHeartbeatInterval = 5 * time.Minute
)

type builder struct {
	supervisor *Supervisor
	options    *Options
}

type Options struct {
	GracefulTimeout   *time.Duration
	HeartbeatInterval *time.Duration
}

func (builder *builder) build() *Supervisor {
	return builder.supervisor
}

func (builder *builder) initialize() *builder {
	builder.supervisor = &Supervisor{}
	if builder.options == nil {
		builder.options = &Options{}
	}
	return builder
}

func (builder *builder) setExitCode() *builder {
	builder.supervisor.exitCode = exitCodeDefault
	return builder
}

func (builder *builder) setWaitGroup() *builder {
	builder.supervisor.waitGroup = &sync.WaitGroup{}
	builder.supervisor.waitChannel = make(chan struct{}, 1)
	return builder
}

func (builder *builder) setRootContext() *builder {
	ctx, cancel := context.WithCancel(context.Background())
	builder.supervisor.rootContext = ctx
	builder.supervisor.rootCanceller = cancel
	return builder
}

func (builder *builder) setSignalChannel() *builder {
	sigchn := make(chan os.Signal, 1)
	signal.Notify(sigchn, syscall.SIGINT, syscall.SIGTERM)
	builder.supervisor.signalChannel = sigchn
	return builder
}

func (builder *builder) setGracefulTimeout() *builder {
	gracefulTimeout := builder.options.GracefulTimeout
	if gracefulTimeout != nil {
		builder.supervisor.gracefulTimeout = *gracefulTimeout
	} else {
		builder.supervisor.gracefulTimeout = defaultGracefulTimeout
	}
	return builder
}

func (builder *builder) setHeartbeatInterval() *builder {
	heartbeatInterval := builder.options.HeartbeatInterval
	if heartbeatInterval != nil {
		builder.supervisor.heartbeatInterval = *heartbeatInterval
	} else {
		builder.supervisor.heartbeatInterval = defaultHeartbeatInterval
	}
	return builder
}

type Daemon struct {
	process  Process
	typeName string
	isActive bool
}

func (daemon *Daemon) setup() {
	if err := daemon.process.Setup(); err != nil {
		panic(err)
	}
}

func (daemon *Daemon) start(ctx context.Context) {
	fields := gblog.Fields{"daemon": daemon.String()}
	if err := daemon.process.Start(ctx); err == nil {
		gblog.WithFields(fields).Info("Daemon succeeded in exiting process.")
	} else {
		gblog.WithFields(fields).WithError(err).Error("Daemon failed to exit process.")
	}
}

func (daemon *Daemon) stamp() {
	daemon.isActive = false
}

func (daemon *Daemon) String() string {
	return fmt.Sprintf("<Daemon| typeName: `%s`, isActive: `%v`>", daemon.typeName, daemon.isActive)
}

type ServerProcess struct {
	server *http.Server
}

func (process *ServerProcess) Setup() error {
	panic("This method hasn't been implemented.")
}

func (process *ServerProcess) Start(ctx context.Context) error {
	go func() {
		if err := process.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			if err := process.server.Shutdown(ctx); err != nil && err != context.Canceled {
				return err
			}
			return nil
		}
	}
}

func (process *ServerProcess) SetServer(server *http.Server) {
	process.server = server
}
