// Package alog provides a simple asynchronous logger that will write to provided io.Writers without blocking calling
// goroutines.
package alog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Alog is a type that defines a logger. It can be used to write log messages synchronously (via the Write method)
// or asynchronously via the channel returned by the MessageChannel accessor.
type Alog struct {
	dest               io.Writer
	m                  *sync.Mutex
	msgCh              chan string
	errorCh            chan error
	shutdownCh         chan struct{}
	shutdownCompleteCh chan struct{}
}

// New creates a new Alog object that writes to the provided io.Writer.
// If nil is provided the output will be directed to os.Stdout.
func New(w io.Writer) *Alog {
		
	if w == nil {
		w = os.Stdout
	}
	return &Alog{
		dest: w,
		msgCh: make(chan string),
		errorCh: make(chan error),
	}
}

// Start begins the message loop for the asynchronous logger. It should be initiated as a goroutine to prevent
// the caller from being blocked.
func (al Alog) Start() {

}

func (al Alog) formatMessage(msg string) string {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return fmt.Sprintf("[%v] - %v", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func (al Alog) write(msg string, wg *sync.WaitGroup) {

	_,error := al.dest.Write([]byte(al.formatMessage(msg)))

	if error != nil {
		al.errorCh <- error
	}
	

	
}

func (al Alog) shutdown() {
}

// MessageChannel returns a channel that accepts messages that should be written to the log.
func (al Alog) MessageChannel() chan<- string {
	// return send only channel msgCh

	return al.msgCh // send only channel

}

// ErrorChannel returns a channel that will be populated when an error is raised during a write operation.
// This channel should always be monitored in some way to prevent deadlock goroutines from being generated
// when errors occur.
func (al Alog) ErrorChannel() <-chan error {
	return al.errorCh
}

// Stop shuts down the logger. It will wait for all pending messages to be written and then return.
// The logger will no longer function after this method has been called.
func (al Alog) Stop() {
}

// Write synchronously sends the message to the log output
func (al Alog) Write(msg string) (int, error) {

	result,error := al.dest.Write([]byte(al.formatMessage(msg)))
	if error != nil {
		return 0,error
	}
	return result, nil;
}
