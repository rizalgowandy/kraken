package storage

import (
	"fmt"

	"code.uber.internal/infra/kraken/torlib"
)

// InfoHashMissMatchError implements error and contains expected and actual torlib.InfoHash
// TODO (@evelynl): this seems to be a fairly common error
type InfoHashMissMatchError struct {
	expected torlib.InfoHash
	actual   torlib.InfoHash
}

func (ie InfoHashMissMatchError) Error() string {
	return fmt.Sprintf("InfoHash missmatch: expected %s, actual %s", ie.expected.HexString(), ie.actual.HexString())
}

// IsInfoHashMissMatchError returns true if error type is InfoHashMissMatchError
func IsInfoHashMissMatchError(err error) bool {
	switch err.(type) {
	case InfoHashMissMatchError:
		return true
	}
	return false
}

// ConflictedPieceWriteError implements error and contains torrent name and piece index
type ConflictedPieceWriteError struct {
	torrent string
	piece   int
}

func (ce ConflictedPieceWriteError) Error() string {
	return fmt.Sprintf("Another thread is writing to the same piece %d for torrent %s", ce.piece, ce.torrent)
}

// IsConflictedPieceWriteError returns true if error type is ConflictedPieceWriteError
func IsConflictedPieceWriteError(err error) bool {
	switch err.(type) {
	case ConflictedPieceWriteError:
		return true
	}
	return false
}