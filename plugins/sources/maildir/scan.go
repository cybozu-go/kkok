package maildir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/log"
)

func parseFile(p string) (*kkok.Alert, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parse(f)
}

// scan scans a Maildir directory dir and generates alerts.
func scan(dir string) []*kkok.Alert {
	files, err := ioutil.ReadDir(filepath.Join(dir, "new"))
	if err != nil {
		return nil
	}

	var alerts []*kkok.Alert
	for _, f := range files {
		if !f.Mode().IsRegular() {
			continue
		}

		fname := f.Name()
		p := filepath.Join(dir, "new", fname)
		a, err := parseFile(p)
		if err != nil {
			log.Error("failed to parse a mail", map[string]interface{}{
				"source":    "maildir",
				log.FnError: err.Error(),
				"dir":       dir,
				"filename":  fname,
			})
		} else {
			alerts = append(alerts, a)
		}
		err = os.Remove(p)
		if err != nil {
			log.Critical("failed to remove a mail", map[string]interface{}{
				"source":    "maildir",
				log.FnError: err.Error(),
				"dir":       dir,
				"filename":  fname,
			})
			return nil
		}
	}

	if len(alerts) > 0 {
		log.Info("new alerts", map[string]interface{}{
			"source": "maildir",
			"count":  len(alerts),
			"dir":    dir,
		})
	}

	return alerts
}
