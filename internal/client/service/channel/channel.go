package channel

type Channels struct {
	SyncProgressBar     chan float64
	SyncProgressBarQuit chan bool
}
