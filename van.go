// Package van provides ...
package van

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/deepzz0/go-van/registry"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// NewService create and returns a new service
func NewService(opts ...Option) Service {
	opt := defaultOptions()
	// process options
	for _, o := range opts {
		o(&opt)
	}
	// service id
	id, err := uuid.NewUUID()
	if err == nil {
		opt.id = id.String()
	}

	return Service{opts: opt}
}

// Service for micro services
type Service struct {
	opts options
}

// Run run the micro service
func (s *Service) Run() error {
	g, ctx := errgroup.WithContext(s.opts.ctx)
	// start server
	if err := s.start(ctx, g); err != nil {
		return err
	}
	// os signal
	ch := make(chan os.Signal, 1)
	if s.opts.signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT,
			syscall.SIGINT)
	}
	// block until got a signal or context done
	g.Go(func() (err error) {
		select {
		case <-ch:
		case <-ctx.Done():
			err = ctx.Err()
		}
		if e := s.stop(ctx, g); e != nil {
			err = e
		}
		return err
	})
	// block unsless error
	if err := g.Wait(); err != nil &&
		!errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// start the service
func (s *Service) start(ctx context.Context,
	g *errgroup.Group) (err error) {

	for _, svr := range s.opts.servers {
		svr := svr
		g.Go(func() error { return svr.Start() })
	}
	// register service
	if s.opts.registry != nil {
		srv := s.regService()
		err = s.opts.registry.Register(ctx, srv)
	}
	return
}

// stop the service
func (s *Service) stop(ctx context.Context,
	g *errgroup.Group) (err error) {
	for _, svr := range s.opts.servers {
		svr := svr
		g.Go(func() error { return svr.Stop() })
	}
	// deregister service
	if s.opts.registry != nil {
		srv := s.regService()
		err = s.opts.registry.Deregister(ctx, srv)
	}
	return
}

// regService discovery registry service
func (s *Service) regService() *registry.Service {
	if len(s.opts.endpoints) == 0 {
		for _, srv := range s.opts.servers {
			if e, err := srv.Endpoint(); err == nil {
				s.opts.endpoints = append(s.opts.endpoints, e)
			}
		}
	}
	return &registry.Service{
		ID:        s.opts.id,
		Name:      s.opts.name,
		Version:   s.opts.version,
		Metadata:  s.opts.metadata,
		Endpoints: s.opts.endpoints,
	}
}
