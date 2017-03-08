package maildir

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func testSourceDirMissing(t *testing.T) {
	t.Parallel()

	_, err := ctor(map[string]interface{}{})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testSourceDirWrongType(t *testing.T) {
	t.Parallel()

	_, err := ctor(map[string]interface{}{
		"dir": 1,
	})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testSourceDirRelative(t *testing.T) {
	t.Parallel()

	_, err := ctor(map[string]interface{}{
		"dir": "relative/path",
	})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testSourceDirNotDirectory(t *testing.T) {
	t.Parallel()

	f, err := ioutil.TempFile("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = ctor(map[string]interface{}{
		"dir": f.Name(),
	})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testSourceDirNotExist(t *testing.T) {
	t.Parallel()

	_, err := ctor(map[string]interface{}{
		"dir": "/not/existing",
	})
	if err != nil {
		t.Error(err)
	}
}

func testSourceIntervalDefault(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	src, err := ctor(map[string]interface{}{
		"dir": d,
	})
	if err != nil {
		t.Fatal(err)
	}
	if src.(*source).interval != defaultInterval*time.Second {
		t.Error(`src.interval != defaultInterval * time.Second`)
	}
}

func testSourceIntervalWrongType(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	_, err = ctor(map[string]interface{}{
		"dir":      d,
		"interval": true,
	})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testSourceIntervalWrongValue(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	_, err = ctor(map[string]interface{}{
		"dir":      d,
		"interval": 0,
	})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testSourceIntervalCustom(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	src, err := ctor(map[string]interface{}{
		"dir":      d,
		"interval": float64(99),
	})
	if err != nil {
		t.Fatal(err)
	}
	if src.(*source).interval != 99*time.Second {
		t.Error(`src.(*source).interval != 99 * time.Second`)
	}
}

func TestSource(t *testing.T) {
	t.Run("dir/missing", testSourceDirMissing)
	t.Run("dir/wrongtype", testSourceDirWrongType)
	t.Run("dir/relative", testSourceDirRelative)
	t.Run("dir/notexist", testSourceDirNotExist)
	t.Run("dir/notdirectory", testSourceDirNotDirectory)
	t.Run("interval/default", testSourceIntervalDefault)
	t.Run("interval/wrongtype", testSourceIntervalWrongType)
	t.Run("interval/wrongvalue", testSourceIntervalWrongValue)
	t.Run("interval/custom", testSourceIntervalCustom)
}
