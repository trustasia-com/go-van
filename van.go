// Package van provides ...
package van

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/trustasia-com/go-van/pkg/registry"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// NewService create and returns a new service
func NewService(opts ...Option) Service {
	options := options{
		// context can not be null
		context: context.Background(),
		signal:  true,
	}
	// process options
	for _, o := range opts {
		o(&options)
	}
	// service id
	id, err := uuid.NewUUID()
	if err == nil {
		options.id = id.String()
	}
	return Service{options: options}
}

// Service for micro services
type Service struct {
	options options
}

// Run run the micro service
func (s *Service) Run() error {
	g, ctx := errgroup.WithContext(s.options.context)
	// start server
	if err := s.start(g); err != nil {
		return err
	}
	// os signal
	ch := make(chan os.Signal, 1)
	if s.options.signal {
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
		if e := s.stop(g); e != nil {
			err = e
		}
		return err
	})
	// block unsless error
	if err := g.Wait(); err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// start the service
func (s *Service) start(g *errgroup.Group) (err error) {

	for _, srv := range s.options.servers {
		srv := srv
		g.Go(func() error { return srv.Start() })
	}
	// register service
	if s.options.registry != nil {
		srv := s.regService()
		err = s.options.registry.Register(s.options.context, srv)
	}
	return
}

// stop the service
func (s *Service) stop(g *errgroup.Group) (err error) {
	// deregister service
	if s.options.registry != nil {
		srv := s.regService()
		err = s.options.registry.Deregister(s.options.context, srv)
	}
	for _, srv := range s.options.servers {
		srv := srv
		g.Go(func() error { return srv.Stop() })
	}
	return
}

// regService discovery registry service
func (s *Service) regService() *registry.Instance {
	if len(s.options.endpoints) == 0 {
		for _, srv := range s.options.servers {
			if e, err := srv.Endpoint(); err == nil {
				s.options.endpoints = append(s.options.endpoints, e)
			}
		}
	}
	return &registry.Instance{
		ID:        s.options.id,
		Name:      s.options.name,
		Version:   s.options.version,
		Metadata:  s.options.metadata,
		Endpoints: s.options.endpoints,
	}
}
