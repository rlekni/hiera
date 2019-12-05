package hiera

import (
	"fmt"
	"io"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/dgoyaml/yaml"
	"github.com/lyraproj/hiera/hieraapi"
)

// RenderName is the name of the option value that describes how to render output
type RenderName string

const (
	// YAML render output in YAML
	YAML = RenderName(`yaml`)
	// JSON render output in JSON
	JSON = RenderName(`json`)
	// Binary render output as binary data
	Binary = RenderName(`binary`)
	// Text render output as plain text
	Text = RenderName(`s`)
)

// Render renders a value on a writer using a specified RenderName
func Render(s hieraapi.Session, renderAs RenderName, value dgo.Value, out io.Writer) {
	switch renderAs {
	case JSON:
		if value.Equals(vf.Nil) {
			util.WriteString(out, "null\n")
		} else {
			// Convert value to rich data format
			opts := streamer.DefaultOptions()
			opts.DedupLevel = streamer.NoDedup
			ser := streamer.New(s.AliasMap(), opts)
			ser.Stream(value, streamer.JSON(out))
			util.WriteByte(out, '\n')
		}

	case YAML:
		if value.Equals(vf.Nil) {
			util.WriteString(out, "\n")
		}
		// Convert value to rich data format
		ser := streamer.New(s.AliasMap(), streamer.DefaultOptions())
		dc := streamer.DataCollector()
		ser.Stream(value, dc)

		bs, err := yaml.Marshal(dc.Value())
		if err != nil {
			panic(err)
		}
		util.WriteString(out, string(bs))
	case Binary:
		bi := vf.New(typ.Binary, value).(dgo.Binary)
		_, err := out.Write(bi.GoBytes())
		if err != nil {
			panic(err)
		}
	case Text:
		util.Fprintln(out, value)
	default:
		panic(fmt.Errorf(`unknown rendering '%s'`, renderAs))
	}
}
