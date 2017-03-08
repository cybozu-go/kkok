package maildir

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func TestScan(t *testing.T) {
	t.Parallel()

	files, err := ioutil.ReadDir("testdata/new")
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ndir := filepath.Join(dir, "new")
	if err := os.Mkdir(ndir, 0755); err != nil {
		t.Fatal(err)
	}

	for _, fi := range files {
		err := copyFile(filepath.Join("testdata/new", fi.Name()),
			filepath.Join(ndir, fi.Name()))
		if err != nil {
			t.Fatal(err)
		}
	}

	alerts := scan(dir)
	if len(alerts) != len(files) {
		t.Errorf(`len(alerts) != len(files) (%d, %d)`, len(alerts), len(files))
	}
}
