// Package stack implements the stack command line.
package stack

import (
	"database/sql"
	"log/slog"

	"github.com/endobit/stack/internal/generated/go/db"
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
	"github.com/endobit/stack/internal/logging"
)

//go:generate gotip tool "github.com/dmarkham/enumer" -type AttrScope -linecomment -text

type AttrScope int

const (
	GlobalScope      AttrScope = iota // global
	ModelScope                        // model
	ZoneScope                         // zone
	ClusterScope                      // cluster
	RackScope                         // rack
	ApplianceScope                    // appliance
	EnvironmentScope                  // environment
	HostScope                         // host
	SwitchScope                       // switch
	BMCScope                          // bmc
)

// Service implements the stackd grpc service.
type Service struct {
	pb.UnimplementedStackServiceServer
	logger *slog.Logger
	db     *db.Queries
}

// WithDB is an option setting function for NewService. It sets the db to db.
func WithDB(database *sql.DB) func(*Service) {
	return func(s *Service) {
		s.db = db.New(logging.DB{
			DB:     database,
			Logger: s.logger,
			Level:  slog.LevelDebug,
		})
	}
}

// WithLogger is an option setting function for NewService. It sets the logger
// to l.
func WithLogger(l *slog.Logger) func(*Service) {
	return func(s *Service) {
		s.logger = l
	}
}

// NewService returns a new stack service.
func NewService(opts ...func(*Service)) *Service {
	svc := Service{
		logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(&svc)
	}

	return &svc
}

// Logger returns the logger.
func (s Service) Logger() *slog.Logger {
	return s.logger
}
