package main

import (
	ll "github.com/sirupsen/logrus"
)

type Servers struct {
	listen listener
	errors chan error
}

// Wait waits until the end of the execution of the server.
func (s *Servers) Wait() error {
	err := <-s.errors
	s.Close()
	return err
}

// Close closes all listening connections
func (s *Servers) Close() {
	s.listen.Close()
}

// Start will start the server asynchronously. See `Wait` to wait until the execution ends.
func Start() (*Servers, error) {
	srv := Servers{
		errors: make(chan error),
	}

	// listen
	ll.Infoln("Starting DHCPv4 server")

	l4, err := listen4()
	if err != nil {
		goto cleanup
	}
	srv.listen = l4
	go func() {
		srv.errors <- l4.Serve()
	}()

	return &srv, nil

cleanup:
	srv.Close()
	return nil, err
}
