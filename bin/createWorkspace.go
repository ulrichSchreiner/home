//usr/bin/env go run $0 $@ ; exit

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const autoenv = `
fn=$(readlink -f $1)
dn=$(dirname $fn)
export GOPATH=$dn
`

func check(msg string, err error) {
	if err != nil {
		log.Fatalf(msg+" :%s", err)
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("use with a git repo, eg. git@github.com:ulrichSchreiner/gl.git")
		return
	}
	u, e := url.Parse(os.Args[1])
	check("parse given url", e)

	if u.Scheme == "" {
		// ssh git repo
		u, e = url.Parse("ssh://" + strings.Replace(os.Args[1], ":", "/", 1))
	}
	var cmd string
	if u.User != nil {
		cmd = u.User.Username()
	}
	// if there is an extension (.git), remove it, we don't need it
	ext := filepath.Ext(u.Path)
	if ext != "" {
		u.Path = u.Path[0 : len(u.Path)-len(ext)]
		// perhaps an anonymous clone via https? use extension as command
		if cmd == "" {
			cmd = ext[1:]
		}
	}

	dir, wks := filepath.Split(u.Path)
	wd, err := os.Getwd()
	check("current workingdir", err)

	wksdir := filepath.Join(wd, wks)

	check("mkdir bin", os.MkdirAll(fmt.Sprintf("%s/bin", wksdir), 0755))
	check("mkdir pkg", os.MkdirAll(fmt.Sprintf("%s/pkg", wksdir), 0755))
	check("mkdir src", os.MkdirAll(fmt.Sprintf("%s/src/%s/%s", wksdir, u.Host, dir), 0755))
	check("write .env", ioutil.WriteFile(fmt.Sprintf("%s/.env", wksdir), []byte(autoenv), 0755))

	if cmd == "" {
		log.Fatal("unknown clone command")
	}
	goget := exec.Command(cmd, "clone", os.Args[1])
	goget.Dir = fmt.Sprintf("%s/src/%s/%s", wksdir, u.Host, dir)
	goget.Env = []string{fmt.Sprintf("GOPATH=%s", wksdir), fmt.Sprintf("PATH=%s", os.ExpandEnv("$PATH"))}
	out, err := goget.CombinedOutput()
	check(string(out), err)
}
