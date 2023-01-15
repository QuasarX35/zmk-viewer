package img

import (
	"fmt"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/mrmarble/zmk-viewer/pkg/keyboard"
	"github.com/mrmarble/zmk-viewer/pkg/keymap"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	keySize  = 60.0
	spacer   = 5.0
	margin   = keySize / 8
	radius   = 5.0
	fontSize = 10.0
)

// ParseKeymap returns struct from a .keymap file.
func ParseKeymap(file string) (*keymap.File, bool) {
	if file == "" {
		return nil, false
	}
	log.Info().Msg("Parsing keymap file.")
	r, err := os.Open(file)
	if err != nil {
		log.Err(err).Send()
		return nil, false
	}

	ast, err := keymap.Parse(r)
	defer r.Close()

	if err != nil {
		log.Err(err).Send()
		return nil, false
	}
	return ast, true
}

// CreateContext from the calculated keyboard size.
func CreateContext(layout *keyboard.Layout) *gg.Context {
	mx := maxX(layout.Layout) + 1
	my := maxY(layout.Layout) + 1

	imageW := int((mx*keySize)+(mx*spacer)) + spacer
	imageH := int(math.Ceil((my*keySize)+(my+1.)*spacer) + (fontSize + spacer*2))

	log.Debug().Int("Image Width", imageW).Int("Image Height", imageH).Float64("Max X", mx).Float64("Max Y", my).Send()

	return gg.NewContext(imageW, imageH)
}

// drawLaout of the keyboard. Blank keys.
func DrawLayout(ctx *gg.Context, transparent bool, layout keyboard.Layout) error {
	if !transparent {
		ctx.SetHexColor("#eeeeee")
		ctx.Clear()
	}

	for _, key := range layout.Layout {
		x, y := getKeyCoords(key)
		w := keySize
		h := keySize
		if key.H != nil {
			h = *key.H * keySize
		}
		if key.W != nil {
			w = *key.W * keySize
		}
		ctx.DrawRoundedRectangle(x, y, w, h, radius)
		ctx.SetRGB(0., 0., 0.)
		ctx.StrokePreserve()
		ctx.SetHexColor("#888888")
		ctx.Fill()
		ctx.DrawRoundedRectangle(x+(margin*2)/2, y+2, w-(margin*2), h-(margin*2), radius)
		ctx.SetHexColor("#a7a7a7")
		ctx.Fill()
	}
	return nil
}

// DrawKeymap of the keyboard. Legend on top of the keys.
func DrawKeymap(ctx *gg.Context, layout keyboard.Layout, layer *keymap.Layer) error {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 10})
	ctx.SetFontFace(face)

	ctx.SetRGB(0., 0., 0.)
	ctx.DrawString(layer.Name, spacer, fontSize+spacer)

	for i, key := range layout.Layout {
		x, y := getKeyCoords(key)
		drawBehavior(ctx, layer.Bindings[i], x+margin+3, y+margin*2.5)
	}
	return nil
}

func getKeyCoords(key keyboard.Key) (float64, float64) {
	x := key.X*keySize + spacer*key.X + spacer
	y := key.Y*keySize + spacer*key.Y + (fontSize + spacer*2)

	return x, y
}

func drawBehavior(ctx *gg.Context, key *keymap.Behavior, x float64, y float64) {
	log.Debug().Str("Action", key.Action).Interface("Params", key.Params).Send()
	ctx.SetRGB(0., 0., 0.)
	for i, v := range key.Params {
		str := ""
		if v.KeyCode == nil {
			str += fmt.Sprintf("%v", *v.Number)
		} else {
			str += *v.KeyCode
		}

		_, dh := ctx.MeasureString(str)
		ctx.DrawString(str, x, y-dh/2.+float64(i)*10.)

	}
}

func maxX(l []keyboard.Key) float64 {
	curr := 0.
	for _, v := range l {
		if v.X > curr {
			curr = v.X
		}
	}
	return curr
}

func maxY(l []keyboard.Key) float64 {
	curr := 0.
	for _, v := range l {
		if v.Y > curr {
			curr = v.Y
		}
	}
	return curr
}
