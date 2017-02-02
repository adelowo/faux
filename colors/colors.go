package colors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/influx6/faux/utils"
)

//==============================================================================

var hsl = regexp.MustCompile("hsl\\((\\d+),\\s*([\\d.]+)%,\\s*([\\d.]+)%\\)")

// IsHSL returns true/false if the giving string is a rgba format data.
func IsHSL(c string) bool {
	return hsl.MatchString(c)
}

// colorReg defines a regexp for matching rgb/rgba header content.
var colorReg = regexp.MustCompile("[rgb|rgba]\\(([\\d\\.,\\s]+)\\)")

// IsRGBFormat returns true/false if the giving string is a rgb/rgba format data.
func IsRGBFormat(c string) bool {
	return colorReg.MatchString(c)
}

// rgbHeader defines a regexp for matching rgb/rgba header content.
var rgbHeader = regexp.MustCompile("rgb\\(([\\d\\.,\\s]+)\\)")

// IsRGB returns true/false if the giving string is a rgb format data.
func IsRGB(c string) bool {
	return rgbHeader.MatchString(c)
}

// rgbaHeader defines a regexp for matching rgb/rgba header content.
var rgbaHeader = regexp.MustCompile("rgba\\(([\\d\\.,\\s]+)\\)")

// IsRGBA returns true/false if the giving string is a rgba format data.
func IsRGBA(c string) bool {
	return rgbaHeader.MatchString(c)
}

// ParseHSL pulls out the rgb/rgba information from a hsl color format  from the
// provided string.
func ParseHSL(rgbData string) (float64, float64, float64) {
	subs := hsl.FindStringSubmatch(rgbData)

	h := utils.ParseFloat(subs[1]) / 360
	s := utils.ParseFloat(subs[2]) / 100
	l := utils.ParseFloat(subs[3]) / 100

	return h, s, l
}

// HSL2RGB converts color values in hsl to rgb.
func HSL2RGB(h, s, l float64) (int, int, int) {
	if s == 0 {
		return int(l), int(l), int(l)
	}

	var q float64

	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}

	p := 2 * (l - q)

	r := Hue(p, q, h+1/3)
	g := Hue(p, q, h)
	b := Hue(p, q, h-1/3)

	return int(r), int(g), int(b)
}

// RGB returns the value of the rgb values as a string.
func RGB(r, g, b, alpha int) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r*255, g*255, b*255, float64(alpha)/100)
}

// Hue takes the provided values and returns a hue value.
func Hue(p, q, t float64) float64 {
	if t < 0 {
		t++
	}

	if t > 1 {
		t--
	}

	if t < 1/6 {
		return p + (q-p)*6*t
	}

	if t < 1/2 {
		return q
	}
	if t < 2/3 {
		return p + (q-p)*(2/3-t)*6
	}

	return p
}

// ParseRGB pulls out the rgb/rgba information from a rgba(9,9,9,9) type
// formatted string.
func ParseRGB(rgbData string) (int, int, int, float64) {
	subs := colorReg.FindStringSubmatch(rgbData)

	if len(subs) < 2 {
		return 0, 0, 0, 0
	}

	rc := strings.Split(subs[1], ",")

	var r, g, b int
	var alpha float64

	r = utils.ParseInt(rc[0])
	g = utils.ParseInt(rc[1])
	b = utils.ParseInt(rc[2])

	if len(rc) > 3 {
		alpha = utils.ParseFloat(rc[3])
	} else {
		alpha = 1
	}

	return r, g, b, alpha
}

// HexToRGB turns a hexademicmal color into rgba format.
// Returns the read, green and blue values as int.
func HexToRGB(hex string) (red, green, blue int) {
	if strings.HasPrefix(hex, "#") {
		hex = strings.TrimPrefix(hex, "#")
	}

	// We are dealing with a 3 string hex.
	if len(hex) < 6 {
		parts := strings.Split(hex, "")
		red = utils.ParseIntBase16(doubleString(parts[0]))
		green = utils.ParseIntBase16(doubleString(parts[1]))
		blue = utils.ParseIntBase16(doubleString(parts[2]))
		return
	}

	red = utils.ParseIntBase16(hex[0:2])
	green = utils.ParseIntBase16(hex[2:4])
	blue = utils.ParseIntBase16(hex[4:6])

	return
}

// HexToRGBA turns a hexademicmal color into rgba format.
// Alpha values ranges from 0-100
func HexToRGBA(hex string, alpha int) string {
	r, g, b := HexToRGB(hex)
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r, g, b, float64(alpha)/100)
}

// doubleString doubles the giving string.
func doubleString(c string) string {
	return fmt.Sprintf("%s%s", c, c)
}
