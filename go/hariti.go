package main

import (
	"./vcs"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type Bundle struct {
	Id   string
	Path string
	Url  string
}

type Vcs struct {
	Install func(*Bundle) error
	Update  func(*Bundle) error
}

var logger *log.Logger = func() *log.Logger {
	logflags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

    if file, err := os.OpenFile(path.Join(path.Dir(os.Args[0]), "../logs/hariti.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err == nil {
		return log.New(file, "", logflags)
	} else {
		return log.New(os.Stderr, "", logflags)
	}
}()

func parseLine(line []byte) (*Vcs, *Bundle, error) {
	items := strings.SplitN(string(line), "\t", 4)
	if len(items) != 4 {
		return nil, nil, errors.New("Too few arguments.")
	}
	id, vcsName, url, path := items[0], items[1], items[2], items[3]

    logger.Printf("Given type is `%s'\n", vcsName)

	var vcs *Vcs
	switch vcsName {
	case "git":
        git := git.NewGitWithLogger(logger)
		vcs = &Vcs{
			Install: func(bundle *Bundle) error {
				return git.Install(bundle.Url, bundle.Path)
			},
			Update: func(bundle *Bundle) error {
				return git.Update(bundle.Url, bundle.Path)
			},
		}
	default:
		return nil, nil, fmt.Errorf("Unknown version control system name: ", vcsName)
	}

	bundle := &Bundle{
		Id:   id,
		Url:  url,
		Path: path,
	}

	return vcs, bundle, nil
}

func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func main() {
	var wg sync.WaitGroup

	in := bufio.NewReader(os.Stdin)
	for {
		bytes, _, err := in.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Panic(err)
		}

		vcs, bundle, err := parseLine(bytes)
		if err != nil {
			logger.Panic(err)
		}

		wg.Add(1)
		go func(vcs *Vcs, bundle *Bundle) {
			defer func() {
				fmt.Printf("%s\t<FINISH>\n", bundle.Id)
				wg.Done()
			}()

			fmt.Printf("%s\t<START>\n", bundle.Id)
			if isDirectory(bundle.Path) {
				err := vcs.Update(bundle)
				if err != nil {
					fmt.Printf("%s\t<ERROR>\t%v\n", bundle.Id, strings.Replace(err.Error(), "\n", "\\n", -1))
				}
			} else {
				err := vcs.Install(bundle)
				if err != nil {
					fmt.Printf("%s\t<ERROR>\t%v\n", bundle.Id, strings.Replace(err.Error(), "\n", "\\n", -1))
				}
			}
		}(vcs, bundle)
	}
	wg.Wait()
}
