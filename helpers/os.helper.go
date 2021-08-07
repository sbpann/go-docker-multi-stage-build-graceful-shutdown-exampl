package helpers

import (
	"log"
	"os"
	"os/signal"
)

type SignalNotifyArgs struct {
	OSSignal        chan os.Signal
	CatchingSignals []os.Signal
	Callback       func()
}
func SignalNotify(args *SignalNotifyArgs) {
	signal.Notify(args.OSSignal, args.CatchingSignals...)
	log.Printf("Received %s signal", (<- args.OSSignal).String())
	args.Callback()
}