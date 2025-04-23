package app

import (
	"context"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	gaugeTimer *gauge.Gauge
	disType    *segmentdisplay.SegmentDisplay
	txtInfo    *text.Text
	// txtTimer   *text.Text
	txtTimer *segmentdisplay.SegmentDisplay

	updateWidgetTimer chan []int
	updateTxtInfo     chan string
	updateTxtTimer    chan string
	updateTxtType     chan string
}

func (w *widgets) update(timer []int, txtType, txtInfo, txtTimer string,
	redrawCh chan<- bool) {

	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}

	if txtType != "" {
		w.updateTxtType <- txtType
	}

	if txtTimer != "" {
		w.updateTxtTimer <- txtTimer
	}

	if len(timer) > 0 {
		w.updateWidgetTimer <- timer
	}

	redrawCh <- true
}

func newWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error

	w.updateWidgetTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer = make(chan string)

	w.gaugeTimer, err = newGauge(ctx, w.updateWidgetTimer, errorCh)
	if err != nil {
		return nil, err
	}

	w.disType, err = newSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		return nil, err
	}

	w.txtInfo, err = newText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}

	w.txtTimer, err = newSegmentDisplayText(ctx, w.updateTxtTimer, errorCh)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newSegmentDisplayText(ctx context.Context, updateText <-chan string,
	errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}

	// Goroutine to update Segment Display
	go func() {
		for {
			select {
			case t := <-updateText:
				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(t),
				})
			case <-ctx.Done():
				return
			}
		}
	}()

	return sd, nil
}

func newText(ctx context.Context, updateText <-chan string,
	errorCh chan<- error) (*text.Text, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}

	// Goroutines to update text
	go func() {
		for {
			select {
			case t := <-updateText:
				txt.Reset()
				errorCh <- txt.Write(t)
			case <-ctx.Done():
				return
			}
		}
	}()

	return txt, nil
}

func newGauge(ctx context.Context, gaugeUpdater <-chan []int,
	errorCh chan<- error) (*gauge.Gauge, error) {

	gauge, err := gauge.New(
		gauge.Height(2),
		gauge.Border(linestyle.Round),
		gauge.Color(cell.ColorAqua),
	)

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case g := <-gaugeUpdater:
				if g[0] <= g[1] {
					errorCh <- gauge.Percent(int(g[0] * 100 / g[1]))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return gauge, nil
}

// func newDonut(ctx context.Context, donUpdater <-chan []int,
// 	errorCh chan<- error) (*donut.Donut, error) {

// 	don, err := donut.New(
// 		donut.Clockwise(),
// 		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	go func() {
// 		for {
// 			select {
// 			case d := <-donUpdater:
// 				if d[0] <= d[1] {
// 					errorCh <- don.Absolute(d[0], d[1])
// 				}
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()

// 	return don, nil
// }

func newSegmentDisplay(ctx context.Context, updateText <-chan string,
	errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {

	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}

	// Goroutine to update SegmentDisplay
	go func() {
		for {
			select {
			case t := <-updateText:
				if t == "" {
					t = " "
				}

				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(t),
				})
			case <-ctx.Done():
				return
			}
		}
	}()

	return sd, nil
}
