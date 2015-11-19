package vcs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type Bundle struct {
	Id   string
	Path string
	Url  string
}

type Cmd interface {
	Install(*Bundle) error
	Update(*Bundle) error
}

var drivers = map[string]Cmd{}

func Register(name string, vcscmd Cmd) {
	drivers[name] = vcscmd
}

func Command(name string) Cmd {
	vcscmd, _ := drivers[name]
	return vcscmd
}

func Run(name string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Got an error: %s\n---\n%s\n", err, stderr.String())
		return fmt.Errorf("%s\n%v", stderr.String(), err)
	}
	return nil
}
