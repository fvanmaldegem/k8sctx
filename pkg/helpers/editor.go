package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func OpenAndRead() *api.Config {
	f, err := createTempFile()
	if err != nil {
		panic(err)
	}
	defer remove(f)

	// close the file and open in the default editor
	f.Close()
	openFileInDefaultEditor(f)

	c, err := clientcmd.LoadFromFile(f.Name())
	if err != nil {
		panic(err)
	}
	return c
}

func remove(f *os.File) {
	err := os.Remove(f.Name())
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func createTempFile() (*os.File, error) {
	f, err := os.CreateTemp("", "k8sctx-*.yml")
	if err != nil {
		return nil, err
	}

	return f, nil
}

func openFileInDefaultEditor(f *os.File) {
	cmd := exec.Command(getEditor(), f.Name())
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", f.Name(), "-e")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic("oops")
	}
}

func getEditor() string {
	editor, found := os.LookupEnv("EDITOR")
	if found {
		return editor
	}

	editor, found = os.LookupEnv("VISUAL")
	if found {
		return editor
	}

	if runtime.GOOS == "windows" {
		return "notepad"
	}

	if runtime.GOOS == "linux" {
		return "vi"
	}

	return ""
}
