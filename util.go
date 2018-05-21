package viddl

import (
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func CreateLoadingBar(p *mpb.Progress) *mpb.Bar {
	return p.AddBar(100,
		mpb.PrependDecorators(
			// Display our static name with one space on the right
			decor.StaticName("video", len("video")+1, decor.DidentRight),
			// DwidthSync bit enables same column width synchronization
			decor.Percentage(0, decor.DwidthSync),
		),
		mpb.AppendDecorators(
			// Replace our ETA decorator with "done!", on bar completion event
			decor.OnComplete(decor.ETA(3, 0), "done!", 0, 0),
		),
	)
}
