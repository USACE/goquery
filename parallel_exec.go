package goquery

import (
	"context"
	"fmt"
	"sync"
)

// Executable represents any operation that can be executed and return an error
type Executable interface {
	Execute() error
}

// ExecQuery is a function that executes a query
type ExecQuery func() error

// ParallelExecutor executes multiple queries concurrently with a configurable concurrency limit
type ParallelExecutor struct {
	queries     []ExecQuery
	maxParallel int
}

// DeferredExec wraps a DataStore.Exec call for deferred execution
type DeferredExec struct {
	store  DataStore
	tx     *Tx
	stmt   string
	params []interface{}
}

// Execute runs the deferred Exec call
func (d *DeferredExec) Execute() error {
	return d.store.Exec(d.tx, d.stmt, d.params...)
}

// NewDeferredExec creates a deferred Exec call that can be added to ParallelExecutor
// Usage: executor.Add(goquery.NewDeferredExec(store, goquery.NoTx, "UPDATE ...", param1, param2))
func NewDeferredExec(store DataStore, tx *Tx, stmt string, params ...interface{}) *DeferredExec {
	return &DeferredExec{
		store:  store,
		tx:     tx,
		stmt:   stmt,
		params: params,
	}
}

// fluentSelectExecutable wraps FluentSelect to implement Executable
type fluentSelectExecutable struct {
	fs *FluentSelect
}

func (f *fluentSelectExecutable) Execute() error {
	return f.fs.Fetch()
}

// fluentInsertExecutable wraps FluentInsert to implement Executable
type fluentInsertExecutable struct {
	fi *FluentInsert
}

func (f *fluentInsertExecutable) Execute() error {
	return f.fi.Execute()
}

// NewParallelExecutor creates a new parallel query executor with sensible defaults
func NewParallelExecutor() *ParallelExecutor {
	return &ParallelExecutor{
		queries:     make([]ExecQuery, 0),
		maxParallel: 10, // Default concurrency
	}
}

// Add appends queries to the executor. Accepts:
// - func() error (standard function)
// - *FluentSelect (will call .Fetch() automatically)
// - *FluentInsert (will call .Execute() automatically)
// - *DeferredExec (for store.Exec calls)
// - Executable interface (will call .Execute())
func (p *ParallelExecutor) Add(queries ...interface{}) *ParallelExecutor {
	for _, q := range queries {
		switch v := q.(type) {
		case func() error:
			// Direct function
			p.queries = append(p.queries, v)
		case *FluentSelect:
			// FluentSelect - wrap to call Fetch()
			p.queries = append(p.queries, func() error {
				return v.Fetch()
			})
		case *FluentInsert:
			// FluentInsert - wrap to call Execute()
			p.queries = append(p.queries, func() error {
				return v.Execute()
			})
		case *DeferredExec:
			// DeferredExec - wrap to call Execute()
			p.queries = append(p.queries, func() error {
				return v.Execute()
			})
		case Executable:
			// Generic Executable interface
			p.queries = append(p.queries, func() error {
				return v.Execute()
			})
		default:
			panic(fmt.Sprintf("unsupported query type: %T. Must be func() error, *FluentSelect, *FluentInsert, *DeferredExec, or Executable", q))
		}
	}
	return p
}

// MaxConcurrency sets the maximum number of concurrent queries
func (p *ParallelExecutor) MaxConcurrency(n int) *ParallelExecutor {
	if n <= 0 {
		panic("maxParallel must be greater than 0")
	}
	p.maxParallel = n
	return p
}

// Run executes all queries in parallel with the configured concurrency limit.
// Returns the first error encountered, if any. If an error occurs, remaining
// queries are cancelled via context.
func (p *ParallelExecutor) Run(ctx context.Context) error {
	if len(p.queries) == 0 {
		return nil
	}

	if p.maxParallel <= 0 {
		return fmt.Errorf("maxParallel must be set to a value greater than 0")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg        sync.WaitGroup
		errOnce   sync.Once
		firstErr  error
		semaphore = make(chan struct{}, p.maxParallel)
	)

	for _, query := range p.queries {
		// Early exit if context is done
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		go func(q ExecQuery) {
			defer wg.Done()

			// Acquire semaphore slot
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			// Execute query
			if err := q(); err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel() // Cancel remaining queries
				})
			}
		}(query)
	}

	wg.Wait()
	return firstErr
}

// RunAll executes all queries in parallel with the configured concurrency limit.
// Unlike Run(), this method does not stop on the first error. It continues
// executing all queries and returns a MultiError containing all errors that occurred.
// Returns nil if all queries succeed.
func (p *ParallelExecutor) RunAll(ctx context.Context) error {
	if len(p.queries) == 0 {
		return nil
	}

	if p.maxParallel <= 0 {
		return fmt.Errorf("maxParallel must be set to a value greater than 0")
	}

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		errors    []*QueryError
		semaphore = make(chan struct{}, p.maxParallel)
	)

	for i, query := range p.queries {
		// Check if context is already cancelled
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		go func(index int, q ExecQuery) {
			defer wg.Done()

			// Acquire semaphore slot
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			// Execute query
			if err := q(); err != nil {
				mu.Lock()
				errors = append(errors, &QueryError{
					Index: index,
					Err:   err,
				})
				mu.Unlock()
			}
		}(i, query)
	}

	wg.Wait()

	if len(errors) > 0 {
		return &MultiError{Errors: errors}
	}

	return nil
}

// QueryError contains information about a failed query execution
type QueryError struct {
	Index int   // Index of the query that failed
	Err   error // The error that occurred
}

// Error implements the error interface
func (qe *QueryError) Error() string {
	return fmt.Sprintf("query %d failed: %v", qe.Index, qe.Err)
}

// Unwrap returns the underlying error
func (qe *QueryError) Unwrap() error {
	return qe.Err
}

// MultiError contains multiple errors from failed queries
type MultiError struct {
	Errors []*QueryError
}

// Error implements the error interface
func (me *MultiError) Error() string {
	if len(me.Errors) == 0 {
		return "no errors"
	}
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}
	return fmt.Sprintf("%d queries failed: first error: %v", len(me.Errors), me.Errors[0].Err)
}

// Unwrap returns the slice of errors for use with errors.Is and errors.As
func (me *MultiError) Unwrap() []error {
	errs := make([]error, len(me.Errors))
	for i, qe := range me.Errors {
		errs[i] = qe
	}
	return errs
}
