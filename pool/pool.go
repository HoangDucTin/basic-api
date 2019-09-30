package pool

type (
	// Pool receives func()
	// to execute and to
	// control the number of
	// functions can be run
	// as the same time.
	pool chan func()
)

var (
	// worker instance keep
	// the channel of the pool
	// for attached functions
	// executing.
	worker *pool
)

// NewPool creates a worker
// pool with the numbers of
// worker 'numb'.
// It also start executing
// the functions pushed in
// the pool.
func NewPool(numb int) {
	pool := make(pool, numb)
	worker = &pool
	go worker.draw()
}

func (p *pool) draw() {
	for {
		if f, ok := <-*p; ok {
			f()
		}
	}
}

func (p *pool) close() {
	if p.available() {
		close(*p)
		*p = nil
		p = nil
	}
}

// Close closes the channel
// of the pool.
func Close() {
	worker.close()
}

func (p *pool) available() bool {
	if p != nil && *p != nil {
		return true
	}
	return false
}

func (p *pool) push(f func()) {
	if p.available() {
		*p <- f
	}
}

// Push pushes functions
// to the pool for executing.
func Push(f func()) {
	worker.push(f)
}
