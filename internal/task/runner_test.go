/*
Copyright (c) 2023 Gemba Advantage

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

package task_test

import (
	"errors"
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	m := &MockedTask{}
	m.On("Run", mock.Anything).Return(nil)
	m.On("Skip", mock.Anything).Return(false)

	err := task.Execute(&context.Context{}, []task.Runner{m})

	require.NoError(t, err)
	m.AssertExpectations(t)
}

func TestExecuteError(t *testing.T) {
	m := &MockedTask{}
	m.On("Run", mock.Anything).Return(errors.New("unexpected error"))
	m.On("Skip", mock.Anything).Return(false)

	err := task.Execute(&context.Context{}, []task.Runner{m})

	require.EqualError(t, err, "unexpected error")
	m.AssertExpectations(t)
}

func TestExecute_Skips(t *testing.T) {
	m := &MockedTask{}
	m.On("Skip", mock.Anything).Return(true)

	err := task.Execute(&context.Context{}, []task.Runner{m})

	require.NoError(t, err)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "Run")
}

type MockedTask struct {
	mock.Mock
}

func (m *MockedTask) Run(ctx *context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockedTask) Skip(ctx *context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}
