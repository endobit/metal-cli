package stack

import (
	"errors"

	"github.com/mattn/go-sqlite3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sqliteToGRPCErrorMap maps sqlite3 error codes to gRPC status codes. The
// comments are from the sqlite3 source code.
var sqliteToGRPCErrorMap = map[sqlite3.ErrNo]codes.Code{
	sqlite3.ErrError:      codes.Unknown,           // SQL error or missing database
	sqlite3.ErrInternal:   codes.Internal,          // Internal logic error in SQLite
	sqlite3.ErrPerm:       codes.PermissionDenied,  // Access permission denied
	sqlite3.ErrAbort:      codes.Aborted,           // Callback routine requested an abort
	sqlite3.ErrBusy:       codes.Unavailable,       // The database file is locked
	sqlite3.ErrLocked:     codes.ResourceExhausted, // A table in the database is locked
	sqlite3.ErrNomem:      codes.ResourceExhausted, // A malloc() failed
	sqlite3.ErrReadonly:   codes.PermissionDenied,  // Attempt to write a readonly database
	sqlite3.ErrInterrupt:  codes.Canceled,          // Operation terminated by sqlite3_interrupt()
	sqlite3.ErrIoErr:      codes.Internal,          // Some kind of disk I/O error occurred
	sqlite3.ErrCorrupt:    codes.DataLoss,          // The database disk image is malformed
	sqlite3.ErrNotFound:   codes.NotFound,          // Unknown opcode in sqlite3_file_control()
	sqlite3.ErrFull:       codes.ResourceExhausted, // Insertion failed because database is full
	sqlite3.ErrCantOpen:   codes.Unavailable,       // Unable to open the database file
	sqlite3.ErrProtocol:   codes.Internal,          // Database lock protocol error
	sqlite3.ErrEmpty:      codes.NotFound,          // Database is empty
	sqlite3.ErrSchema:     codes.Aborted,           // The database schema changed
	sqlite3.ErrTooBig:     codes.InvalidArgument,   // String or BLOB exceeds size limit
	sqlite3.ErrConstraint: codes.AlreadyExists,     // Abort due to constraint violation
	sqlite3.ErrMismatch:   codes.InvalidArgument,   // Data type mismatch
	sqlite3.ErrMisuse:     codes.Internal,          // Library used incorrectly
	sqlite3.ErrNoLFS:      codes.Unimplemented,     // Uses OS features not supported on host
	sqlite3.ErrAuth:       codes.PermissionDenied,  // Authorization denied
	sqlite3.ErrFormat:     codes.DataLoss,          // Auxiliary database format error
	sqlite3.ErrRange:      codes.OutOfRange,        // 2nd parameter to sqlite3_bind out of range
	sqlite3.ErrNotADB:     codes.InvalidArgument,   // File opened that is not a database file
	sqlite3.ErrNotice:     codes.OK,                // Notifications from sqlite3_log()
	sqlite3.ErrWarning:    codes.Unknown,           // Warnings from sqlite3_log()
}

func dberr(err error) error {
	var sqlerr sqlite3.Error

	if errors.As(err, &sqlerr) {
		return status.Errorf(sqliteToGRPCErrorMap[sqlerr.Code], "database: %v", sqlerr)
	}

	return status.Convert(err).Err()
}

var (
	errMissingAppliance   = status.Error(codes.InvalidArgument, "appliance not specified")
	errMissingAttribute   = status.Error(codes.InvalidArgument, "attribute not specified")
	errMissingCluster     = status.Error(codes.InvalidArgument, "cluster not specified")
	errMissingEnvironment = status.Error(codes.InvalidArgument, "environment not specified")
	errMissingHost        = status.Error(codes.InvalidArgument, "host not specified")
	errMissingInterface   = status.Error(codes.InvalidArgument, "inteface not specified")
	errMissingMake        = status.Error(codes.InvalidArgument, "make not specified")
	errMissingModel       = status.Error(codes.InvalidArgument, "model not specified")
	errMissingNetwork     = status.Error(codes.InvalidArgument, "network not specified")
	errMissingRack        = status.Error(codes.InvalidArgument, "rack not specified")
	errMissingZone        = status.Error(codes.InvalidArgument, "zone not specified")
	errNameOrGlob         = status.Error(codes.InvalidArgument, "must provide either name or glob")
)
