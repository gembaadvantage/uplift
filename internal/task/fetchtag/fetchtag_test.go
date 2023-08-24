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

package fetchtag

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "fetching all tags", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.True(t, Task{}.Skip(&context.Context{
		FetchTags: false,
	}))
}

func TestRun(t *testing.T) {
	log := `(tag: 0.2.1) fix: this is a fix
(tag: 0.2.0) feat: this was another feature
(tag: 0.1.0) feat: this was a feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))
	require.Empty(t, gittest.Tags(t))

	err := Task{}.Run(&context.Context{})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0", "0.2.1"}, gittest.Tags(t))
}
