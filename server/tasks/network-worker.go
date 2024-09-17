package tasks

import "sync"

type NetworkWorker interface {
	StartNetWorker()
}
type AcceptWorker interface {
	StartNetWorker(g *sync.WaitGroup)
}
