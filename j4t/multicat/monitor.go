package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
)

type Formatter struct {
}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	writer := tabwriter.NewWriter(entry.Buffer, 0, 0, 1, ' ', 0)
	_, _ = writer.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	_, _ = writer.Write([]byte("\t["))
	_, _ = writer.Write([]byte(entry.Level.String()))
	_, _ = writer.Write([]byte("]\t"))
	var names = make([]string, 0, len(entry.Data))
	for name := range entry.Data {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		_, _ = writer.Write([]byte("["))
		_, _ = writer.Write([]byte(name))
		_, _ = writer.Write([]byte{':'})
		_, _ = writer.Write([]byte(fmt.Sprint(entry.Data[name])))
		_, _ = writer.Write([]byte("]\t"))
	}
	_, _ = writer.Write([]byte(entry.Message))
	_ = writer.Flush()
	entry.Buffer.WriteByte('\n')
	return entry.Buffer.Bytes(), nil
}

type Item struct {
	Name      string
	CalSpeed  bool
	SpeedOnly bool
	Val       uint64
	preVal    uint64
}

type Monitor struct {
	Items []*Item
	t     *time.Ticker
}

func (m *Monitor) Monit(d time.Duration) {
	log.SetFormatter(&Formatter{})
	log.SetOutput(os.Stdout)
	t := time.NewTicker(d)
	m.t = t
	start := time.Now()
	for range t.C {
		loge := log.NewEntry(log.StandardLogger())
		end := time.Now()
		for _, item := range m.Items {
			val := atomic.LoadUint64(&item.Val)
			if !item.SpeedOnly {
				if strings.HasSuffix(item.Name, "mem") {
					loge = loge.WithField(item.Name, Humanize1024Size(int64(val)))
				} else {
					loge = loge.WithField(item.Name, val)
				}
			}
			if item.CalSpeed {
				if strings.HasSuffix(item.Name, "mem") {
					loge = loge.WithField(item.Name+"_speed", fmt.Sprintf("%s/s", Humanize1024Size(int64(float64(val-item.preVal)/end.Sub(start).Seconds()))))
				} else {
					loge = loge.WithField(item.Name+"_speed", fmt.Sprintf("%.2f/s", float64(val-item.preVal)/end.Sub(start).Seconds()))
				}
			}
			item.preVal = val
		}
		loge.Info("monit")
		start = end
	}
}

func (m *Monitor) Stop() {
	m.t.Stop()
}
func (i *Item) Incr() {
	atomic.AddUint64(&i.Val, 1)
}

func (i *Item) IncrN(n uint64) {
	atomic.AddUint64(&i.Val, n)
}

const (
	kib = 1024
	mib = 1024 * 1024
	gib = 1024 * 1024 * 1024
	tib = 1024 * 1024 * 1024 * 1024
)

// Humanize1024Size humanizes size based on powers of 1024
func Humanize1024Size(size int64) string {
	if size < kib {
		return fmt.Sprintf("%d Bytes", size)
	}
	if size < mib {
		return fmt.Sprintf("%.2f KiB", float64(size)/kib)
	}
	if size < gib {
		return fmt.Sprintf("%.2f MiB", float64(size)/mib)
	}
	if size < tib {
		return fmt.Sprintf("%.2f GiB", float64(size)/gib)
	}
	return fmt.Sprintf("%.2f TiB", float64(size)/tib)
}
