package stream

import (
	"errors"
	"io"
	"iter"

	"google.golang.org/grpc"
)

// Reader is a gRPC streaming response reader. It's only purpose is to provide
// an iterator for streaming responses. See NewReader for more details.
type Reader[T any] struct {
	grpc.ServerStreamingClient[T]
	err error
}

// NewReader returns a new Reader. It takes as parameters the return values from
// a gRPC streaming client method. For example:
//
//	r := NewReader(client.ReadZones(ctx, req))
//	for zone, err := range r.Responses() {
//	...
//	}
func NewReader[T any](client grpc.ServerStreamingClient[T], err error) *Reader[T] {
	return &Reader[T]{
		ServerStreamingClient: client,
		err:                   err,
	}
}

// Responses returns an iterator for processing gRPC streaming responses from the Reader.
func (r *Reader[T]) Responses() iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		// constructors are allowed to set and not return an error
		if r.err != nil {
			if !yield(nil, r.err) {
				return
			}
		}

		for {
			resp, err := r.ServerStreamingClient.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if !yield(resp, err) {
				return
			}
		}
	}
}
