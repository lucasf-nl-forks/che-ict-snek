package utils

import (
	"io"

	"github.com/cheggaaa/pb"
)

// ProgressTracker implements the required interface
type ProgressTracker struct {
	bar *pb.ProgressBar
}

// NewProgressTracker creates a new progress tracker instance
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		bar: pb.StartNew(100),
	}
}

// TrackProgress implements the required interface method
func (pt *ProgressTracker) TrackProgress(
	src string,
	currentSize, totalSize int64,
	stream io.ReadCloser,
) (body io.ReadCloser) {
	// Create a proxy reader that tracks progress
	proxyReader := &progressProxyReader{
		Reader:     stream,
		TotalSize:  totalSize,
		CurrentPos: currentSize,
		Bar:        pt.bar,
	}

	// Set the total size for the progress bar
	pt.bar.Total = int64(totalSize)

	// Return our proxy reader wrapped in a ReadCloser
	return &proxyReadCloser{proxyReader}
}

// Helper types for implementing the progress tracking
type progressProxyReader struct {
	io.Reader
	TotalSize  int64
	CurrentPos int64
	Bar        *pb.ProgressBar
}

type proxyReadCloser struct {
	*progressProxyReader
}

func (pr *progressProxyReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	if n > 0 {
		pr.CurrentPos += int64(n)
		pr.Bar.Set64(pr.CurrentPos)
	}
	return
}

func (pr *proxyReadCloser) Close() error {
	if closer, ok := pr.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
