package ixgo

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/ixgo/load"
)

type goModDriver struct {
	listDrvier *load.ListDriver
}

func newGoModDriver() *goModDriver {
	return &goModDriver{listDrvier: new(load.ListDriver)}
}

func (g *goModDriver) Lookup(root string, path string) (dir string, found bool) {
	dir, found = g.listDrvier.Lookup(root, path)
	if found {
		return
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); os.IsNotExist(err) {
		execCommand(root, "go", "mod", "init", filepath.Base(root))
	}
	execCommand(root, "go", "get", path)

	ret := execCommand(root, "go", "mod", "download", "-json", path)

	var modDownload struct {
		Dir string
	}
	json.Unmarshal(ret, &modDownload)

	if modDownload.Dir != "" {
		found = true
		dir = modDownload.Dir
	}

	return
}

func execCommand(dir, mainCmd string, subcmd ...string) []byte {
	cmd := exec.Command(mainCmd, subcmd...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	ret, err := cmd.Output()
	if err != nil {
		panic(string(ret))
	}
	return ret
}
