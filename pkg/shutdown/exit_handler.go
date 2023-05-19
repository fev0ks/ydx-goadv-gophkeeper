package shutdown

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"ydx-goadv-gophkeeper/pkg/logger"
)

var (
	mu = &sync.Mutex{}
)

type ExitHandler struct {
	mainCtxCanceler   context.CancelFunc
	log               *zap.SugaredLogger
	httpServer        *http.Server
	grpcServer        *grpc.Server
	ToCancel          []context.CancelFunc
	ToStop            []chan struct{}
	ToClose           []io.Closer
	ToExecute         []func(ctx context.Context) error
	funcsInProcessing sync.WaitGroup
	newFuncAllowed    bool
}

func NewExitHandlerWithCtx(mainCtxCanceler context.CancelFunc) *ExitHandler {
	return &ExitHandler{
		log:               logger.NewLogger("exit-hdr"),
		mainCtxCanceler:   mainCtxCanceler,
		newFuncAllowed:    true,
		funcsInProcessing: sync.WaitGroup{},
	}
}

func (eh *ExitHandler) IsNewFuncExecutionAllowed() bool {
	mu.Lock()
	defer mu.Unlock()
	return eh.newFuncAllowed
}

func (eh *ExitHandler) setNewFuncExecutionAllowed(value bool) {
	mu.Lock()
	defer mu.Unlock()
	eh.newFuncAllowed = value
}

func (eh *ExitHandler) ShutdownHTTPServerBeforeExit(httpServer *http.Server) {
	eh.httpServer = httpServer
}

func (eh *ExitHandler) ShutdownGrpcServerBeforeExit(grpcServer *grpc.Server) {
	eh.grpcServer = grpcServer
}

func (eh *ExitHandler) AddFuncInProcessing(alias string) {
	mu.Lock()
	defer mu.Unlock()
	eh.log.Infof("'%s' func is started and added to exit handler", alias)
	eh.funcsInProcessing.Add(1)
}

func (eh *ExitHandler) FuncFinished(alias string) {
	mu.Lock()
	defer mu.Unlock()
	eh.log.Infof("'%s' func is finished and removed from exit handler", alias)
	eh.funcsInProcessing.Add(-1)
}

func ProperExitDefer(exitHandler *ExitHandler) chan struct{} {
	exitHandler.log.Info("Graceful exit handler is activated")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	exit := make(chan struct{})
	go func() {
		s := <-signals
		exitHandler.log.Infof("Received a signal '%s'", s)
		exit <- struct{}{}
		exitHandler.setNewFuncExecutionAllowed(false)
		exitHandler.shutdown()
	}()
	return exit
}

func (eh *ExitHandler) shutdown() {
	successfullyFinished := make(chan struct{})
	go func() {
		eh.waitForShutdownServer()
		eh.waitForFinishFunc()
		eh.endHeldObjects()
		successfullyFinished <- struct{}{}
	}()
	select {
	case <-successfullyFinished:
		log.Println("System finished work, graceful shutdown")
		eh.mainCtxCanceler()
	case <-time.After(1 * time.Minute):
		log.Println("System has not shutdown in time '1m', shutdown with interruption")
		os.Exit(1)
	}
}

func (eh *ExitHandler) waitForFinishFunc() {
	log.Println("Waiting for functions finish work...")
	eh.funcsInProcessing.Wait()
	log.Println("All functions finished work successfully")
}

func (eh *ExitHandler) waitForShutdownServer() {
	if eh.httpServer != nil {
		log.Println("Waiting for shutdown http server...")
		err := eh.httpServer.Shutdown(context.Background())
		log.Println("Http Server shutdown complete")
		if err != nil {
			eh.log.Infof("failed to shutdown server: %v", err)
		}
	}
	if eh.grpcServer != nil {
		log.Println("Waiting for shutdown proto server...")
		eh.grpcServer.GracefulStop()
		log.Println("Grpc Server shutdown complete")
	}
}

func (eh *ExitHandler) endHeldObjects() {
	if len(eh.ToExecute) > 0 {
		log.Println("ToExecute final funcs")
		for _, execute := range eh.ToExecute {
			err := execute(context.Background())
			if err != nil {
				eh.log.Infof("func error: %v", err)
			}
		}
	}
	if len(eh.ToCancel) > 0 {
		log.Println("ToCancel active contexts")
		for _, cancel := range eh.ToCancel {
			cancel()
		}
	}
	if len(eh.ToStop) > 0 {
		log.Println("ToStop active goroutines")
		for _, toStop := range eh.ToStop {
			close(toStop)
		}
	}
	if len(eh.ToClose) > 0 {
		log.Println("ToClose active resources")
		for _, toClose := range eh.ToClose {
			err := toClose.Close()
			if err != nil {
				eh.log.Infof("failed to close an resource: %v", err)
			}
		}
	}
	log.Println("Success end final work")
}
