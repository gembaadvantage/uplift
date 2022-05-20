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

package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLog_BreakingFooter(t *testing.T) {
	inc := ParseLog(`
commit 95bdec4c8fe888ae2fd4e6cecea99f5b7ae2a045
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Wed May 18 20:44:10 2022 +0100

    docs: document about new breaking change

commit 0a437181e47e79ac80b683f411677ce94859633a
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 21:13:13 2022 +0100

    fix: annoying bug has now been fixed

commit f51d067556e8cc0eadcabeb5a1f3d27577bc0a84
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 08:33:59 2022 +0100

    refactor: changed a really important part of the app

	BREAKING CHANGE: the cli has been completely refactored with no backwards compatibility

commit a7095058f2b42a87d772a084f427c0e645440308
Author: paul.t <paul.t@gembaadvantage.com>
Date:   Mon May 16 12:12:34 2022 +0100

    docs(config): document new configuration option`)

	assert.Equal(t, MajorIncrement, inc)
}

func TestParseLog_BreakingBang(t *testing.T) {
	inc := ParseLog(`
commit 95bdec4c8fe888ae2fd4e6cecea99f5b7ae2a045
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Wed May 18 20:44:10 2022 +0100

    feat: a new snazzy feature has been added

commit 0a437181e47e79ac80b683f411677ce94859633a
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 21:13:13 2022 +0100

    fix: annoying bug has now been fixed

commit f51d067556e8cc0eadcabeb5a1f3d27577bc0a84
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 08:33:59 2022 +0100

    feat!: changed a really important part of the app`)

	assert.Equal(t, MajorIncrement, inc)
}

func TestParseLog_Minor(t *testing.T) {
	inc := ParseLog(`
commit 95bdec4c8fe888ae2fd4e6cecea99f5b7ae2a045
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Wed May 18 20:44:10 2022 +0100

    ci: change to the existing workflow

commit 0a437181e47e79ac80b683f411677ce94859633a
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 21:13:13 2022 +0100

    fix: annoying bug has now been fixed

commit f51d067556e8cc0eadcabeb5a1f3d27577bc0a84
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 08:33:59 2022 +0100

    feat: shiny new feature has been added`)

	assert.Equal(t, MinorIncrement, inc)
}

func TestParseLog_Patch(t *testing.T) {
	inc := ParseLog(`
commit 95bdec4c8fe888ae2fd4e6cecea99f5b7ae2a045
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Wed May 18 20:44:10 2022 +0100

    ci: change to the existing workflow

commit 0a437181e47e79ac80b683f411677ce94859633a
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 21:13:13 2022 +0100

    docs: updated documented to talk about fix

commit f51d067556e8cc0eadcabeb5a1f3d27577bc0a84
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 08:33:59 2022 +0100

    fix: small bug fixed`)

	assert.Equal(t, PatchIncrement, inc)
}

func TestParseLog_NoIncrement(t *testing.T) {
	inc := ParseLog(`
commit 0a437181e47e79ac80b683f411677ce94859633a
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 21:13:13 2022 +0100

    docs(ci): documented additional CI support

commit f51d067556e8cc0eadcabeb5a1f3d27577bc0a84
Author: Paul T <paul.t@gembaadvantage.com>
Date:   Tue May 17 08:33:59 2022 +0100

    ci: sped up the existing build

commit a7095058f2b42a87d772a084f427c0e645440308
Author: paul.t <paul.t@gembaadvantage.com>
Date:   Mon May 16 12:12:34 2022 +0100

    docs(config): documented new configuration option`)

	assert.Equal(t, NoIncrement, inc)
}
