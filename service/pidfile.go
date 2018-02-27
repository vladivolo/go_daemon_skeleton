package service

import (
	"fmt"
	"os"
	"syscall"
)

type Pidfile struct {
	f    *os.File
	path string
}

func (p *Pidfile) Open(path string) error {
	if path == "" {
		return nil
	}

	pf, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			p.CloseAndRemove()
		}
	}()

	err = syscall.Flock(int(pf.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return err
	}

	err = pf.Truncate(0)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(pf, "%d\n", os.Getpid())
	if err != nil {
		return err
	}

	p.CloseAndRemove() // maybe there's an old one

	p.f = pf
	p.path = path

	return nil
}

func (p *Pidfile) Close() (err error) {
	if p.f != nil {
		err = p.f.Close()
		p.f = nil
	}
	return
}

func (p *Pidfile) CloseAndRemove() (err error) {
	if p.f != nil {
		p.Close()
		err = os.Remove(p.path)
	}
	return
}

func PidfileOpen(path string) (*Pidfile, error) {
	pidfile := &Pidfile{
		f:    nil,
		path: path,
	}

	return pidfile, pidfile.Open(path)
}
