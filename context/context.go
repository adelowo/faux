// Package context is built out of my desire to understand the http context
// library and as an experiement in such a library works.
package context

import (
	"sync"
	"time"
)

//==============================================================================

// Fields defines a map of key:value pairs.
type Fields map[interface{}]interface{}

//==============================================================================

// nilPair defines a nil starting pair.
var nilPair = (*Pair)(nil)

// Pair defines a struct for storing a linked pair of key and values.
type Pair struct {
	prev  *Pair
	key   interface{}
	value interface{}
}

// Append returns a new Pair with the giving key and with the provded Pair set as
// it's previous link.
func Append(p *Pair, key, value interface{}) *Pair {
	return p.Append(key, value)
}

// Fields returns all internal pair data as a map.
func (p *Pair) Fields() Fields {
	var f Fields

	if p.prev == nil {
		f = make(Fields)
		f[p.key] = p.value
		return f
	}

	f = p.prev.Fields()
	f[p.key] = p.value
	return f
}

// Append returns a new pair with the giving key and value and its previous
// set to this pair.
func (p *Pair) Append(key, val interface{}) *Pair {
	return &Pair{
		prev:  p,
		key:   key,
		value: val,
	}
}

// Get collects the value of a key if it exists.
func (p *Pair) Get(key interface{}) (value interface{}, found bool) {
	if p == nil {
		return
	}

	if p.key == key {
		return p.value, true
	}

	if p.prev == nil {
		return
	}

	return p.prev.Get(key)
}

//==============================================================================

// Context defines an interface for a context providers which allows us to
// build passable context around.
type Context interface {

	// IsExpired returns true/false if the context is considered expired.
	IsExpired() bool

	// Get returns the giving value for the provided key if it exists else nil.
	Get(key interface{}) (interface{}, bool)

	// Done returns a channel which gets closed when the given channel
	// expires else closes immediately if its not an expiring context.
	Done() <-chan struct{}

	// Ctx returns a Context which exposes a basic context interface without  the
	// cancellable method.
	Ctx() Context

	// New returns a new context based on the fileds of the context which its
	// called from, it does inherits the lifetime limits of the context its
	// called from.
	New(cancelWithParent bool) CancelableContext

	// WithTimeout returns a new Context from the previous with the given timeout
	// if the timeout is still further than the previous in expiration date else uses
	// the previous expiration date instead since that is still further in the future.
	WithTimeout(timeout time.Duration, cancelWithParent bool) CancelableContext

	// Set adds a key and value pair into the context store.
	Set(key interface{}, value interface{})

	// // SetParent adds a key and value pair into the context store.
	// SetParent(key interface{}, value interface{})

	// WithValue returns a new context then adds the key and value pair into the
	// context's store.
	WithValue(key interface{}, value interface{}) CancelableContext

	// TimeRemaining returns the remaining time for expiring of the context if it
	// indeed has an expiration date set and returns a bool value indicating if it
	// has a timeout.
	TimeRemaining() (remaining time.Duration, hasTimeout bool)
}

// CancelableContext defines a interface which exposes the capability to cancel
// a Context.
type CancelableContext interface {
	Context

	//Cancel cancels out the timer setup to nil out contexts internal store.
	Cancel()
}

// New returns a new context instance.
func New() CancelableContext {
	cl := context{
		fields:    nilPair,
		canceller: make(chan struct{}),
	}
	return &cl
}

//==============================================================================

// context defines a struct for bundling a context against specific
// use cases with a explicitly set duration which clears all its internal
// data after the giving period.
type context struct {
	fields    *Pair
	lifetime  time.Time
	timer     *time.Timer
	duration  time.Duration
	parent    Context
	canceller chan struct{}
	cl        sync.Mutex
	canceled  bool
}

// New returns a new context from with the configuration limits of this one.
func (c *context) New(cancelWithParent bool) CancelableContext {
	if c.timer != nil {
		return c.WithTimeout(c.duration, cancelWithParent)
	}

	return c.newChild(cancelWithParent)
}

// WithTimeout returns a new context whoes internal value expires
// after the giving duration.
func (c *context) WithTimeout(life time.Duration, cancelWithParent bool) CancelableContext {
	child := c.newChild(cancelWithParent)

	var useChild bool

	lifetime := time.Now().Add(life)
	if lifetime.After(child.lifetime) {
		child.duration = life
		child.lifetime = lifetime
		useChild = true
	}

	var to time.Duration

	if useChild {
		to = life
	} else {
		to = c.duration
	}

	child.timer = time.AfterFunc(to, func() {
		child.fields = nilPair
		child.Cancel()
	})

	return child
}

// WithValue returns a new context based on the previos one.
func (c *context) WithValue(key, value interface{}) CancelableContext {
	child := c.newChild(true)
	child.fields = Append(child.fields, key, value)
	return child
}

// TimeRemaining returns the remaining time before expiration.
func (c *context) TimeRemaining() (rem time.Duration, hasTimeout bool) {
	if c.lifetime.IsZero() {
		return
	}

	hasTimeout = true

	now := time.Now()
	if now.Before(c.lifetime) {
		rem = c.lifetime.Sub(now)
		return
	}

	return
}

// Done returns a channel which gets closed when the context
// has expired.
func (c *context) Done() <-chan struct{} {
	if c.IsExpired() {
		return nil
	}

	return c.canceller
}

// IsExpired returns true/false if the context internal data has expired.
func (c *context) IsExpired() bool {
	left, has := c.TimeRemaining()
	if has {
		if left <= 0 {
			return true
		}
	}

	c.cl.Lock()
	{
		if c.canceled {
			c.cl.Unlock()
			return true
		}
	}
	c.cl.Unlock()

	return false
}

// Ctx returns the Context interface  for a giving CancelableContext.
func (c *context) Ctx() Context {
	return c
}

// Cancel cancels the timer if there exists one set to clear context.
func (c *context) Cancel() {
	if c.IsExpired() {
		return
	}

	c.cl.Lock()
	c.canceled = true
	c.cl.Unlock()

	close(c.canceller)

	if c.timer != nil {
		c.timer.Stop()
		return
	}
}

// SetParent adds the giving value using the given key into the map of the
// root parent of the context to have this persist to new context else sets
// the value on itself if it has no parent.
func (c *context) SetParent(key, val interface{}) {
	if c.parent == nil {
		c.parent.Set(key, val)
		return
	}

	c.fields = Append(c.fields, key, val)
}

// Set adds the giving value using the given key into the map.
func (c *context) Set(key, val interface{}) {
	c.fields = Append(c.fields, key, val)
}

// Get returns the value for the necessary key within the context.
func (c *context) Get(key interface{}) (item interface{}, found bool) {
	item, found = c.fields.Get(key)
	return
}

// newChild returns a new fresh context based on the fields of this context.
func (c *context) newChild(cancelWithParent bool) *context {
	canceller := make(chan struct{})

	if c.IsExpired() {
		close(canceller)
	}

	cm := &context{
		parent:    c,
		fields:    c.fields,
		lifetime:  c.lifetime,
		duration:  c.duration,
		canceled:  c.canceled,
		canceller: canceller,
	}

	if cancelWithParent {
		go func() {
			cancel := c.Done()
			if cancel == nil {
				cm.Cancel()
				return
			}

			<-cancel
			cm.Cancel()
		}()
	}

	return cm
}
