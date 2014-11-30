package main

import (
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	cmds map[string]string
)

func checkError(e error) {
	if e != nil {
		log.Fatalf("%s", e)
	}
}

func execCmd(f string, arg ...string) {
	// Check to see if program is in $PATH
	_, e := exec.LookPath(f)
	checkError(e)

	cmd := exec.Command(f, arg...)

	out, _ := cmd.Output()
	s := string(out)
	if len(s) > 1 {
		log.Println("Output for", f, s)
	} else {
		log.Println("Executed", f, arg)
	}
}

func findCommand(n string) {

	for k, v := range cmds {
		m, e := filepath.Glob(k)
		checkError(e)
		l := len(m)

		for i := 0; i < l; i++ {
			if m[i] == n {
				execCmd(v, m[i])
			}
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	watcher, e := fsnotify.NewWatcher()
	checkError(e)
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file: ", ev.Name)
					findCommand(ev.Name)
				}
			case er := <-watcher.Errors:
				checkError(er)
			}
		}
	}()

	m := os.Args[1:]
	l := len(m)
	cmds = make(map[string]string)

	for i := 0; i < l; i++ {
		if i%2 == 0 {
			if cmds[m[i]] == "" {
				cmds[m[i]] = m[i+1]
			}
		}
	}

	for k, _ := range cmds {
		e := watcher.Add(k)
		checkError(e)
	}

	<-done
}
