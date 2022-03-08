/*
Copyright (c) 2022 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package bump

import (
	"io/ioutil"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/tidwall/sjson"
)

func jsonBump(ctx *context.Context, path string, bumps []config.JSONBump) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	str := string(data)

	for _, bump := range bumps {
		// Strip any 'v' prefix if this must be a semantic version
		v := ctx.NextVersion.Raw
		if bump.SemVer {
			v = strictSemVer(v)
		}

		str, err = sjson.Set(str, bump.Path, v)
		if err != nil {
			return false, err
		}
	}

	return true, ioutil.WriteFile(path, []byte(str), 0644)
}
