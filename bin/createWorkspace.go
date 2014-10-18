package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
		fmt.Println("use with a projectname, eg. github.com/ulrichSchreiner/gl")
		return
	}
	pt := os.Args[1]
	dir, wks := filepath.Split(pt)
	wd, err := os.Getwd()
	check("current workingdir", err)

	wksdir := filepath.Join(wd, wks)

	check("mkdir bin", os.MkdirAll(fmt.Sprintf("%s/bin", wksdir), 0755))
	check("mkdir pkg", os.MkdirAll(fmt.Sprintf("%s/pkg", wksdir), 0755))
	check("mkdir src", os.MkdirAll(fmt.Sprintf("%s/src/%s", wksdir, dir), 0755))
	check("write .env", ioutil.WriteFile(fmt.Sprintf("%s/.env", wksdir), []byte(autoenv), 0755))

	goget := exec.Command("go", "get", os.Args[1])
	goget.Dir = fmt.Sprintf("%s/src/%s", wksdir, dir)
	goget.Env = []string{fmt.Sprintf("GOPATH=%s", wksdir), fmt.Sprintf("PATH=%s", os.ExpandEnv("$PATH"))}
	out, err := goget.CombinedOutput()
	check(string(out), err)
}
