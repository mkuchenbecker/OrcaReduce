package exceptions

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mkuchenbecker/orcareduce/orcareduce/exceptions"
	"github.com/mkuchenbecker/orcareduce/orcareduce/exceptions/mock"
	"github.com/stretchr/testify/assert"
)

var errGeneric = fmt.Errorf("error")

func TestHandler_RunAsync(t *testing.T) {
	t.Parallel()

	t.Run("SyncFunc Synchronizes", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockLogger := mock.NewMockLogger(mockCtrl)
		handler := NewHandler(mockLogger)

		i := 0
		runFunc := func() error {
			time.Sleep(10 * time.Millisecond)
			i = 1
			return nil
		}

		sync := handler.RunAsync(runFunc)
		assert.Equal(t, 0, i)
		sync()
		assert.Equal(t, 1, i)
	})

	t.Run("Panic Handled and Logged", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logs := mock.NewMockLogger(mockCtrl)
		logs.EXPECT().Errorf("encountered a panic: %+v", gomock.Not(nil)).Times(1)
		handler := NewHandler(logs)

		runFunc := func() error {
			panic("")
		}
		sync := handler.RunAsync(runFunc)
		sync()
	})
}

func TestHandler_HandlePanic(t *testing.T) {
	t.Parallel()

	t.Run("Function Does Not Panic", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockLogger := mock.NewMockLogger(mockCtrl)
		handler := NewHandler(mockLogger)

		noOpFunction := func() error { return nil }
		sync := handler.RunAsync(noOpFunction)
		assert.NoError(t, sync())
	})

	t.Run("Panic Logs and Returns Error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logs := mock.NewMockLogger(mockCtrl)
		logs.EXPECT().Errorf("encountered a panic: %+v", gomock.Not(nil)).Times(1)
		handler := NewHandler(logs)

		runFunc := func() error {
			panic("")
		}
		sync := handler.RunAsync(runFunc)
		err := sync()
		assert.Error(t, err)
	})

	t.Run("Panic Populates Return Error and Logs Error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logs := mock.NewMockLogger(mockCtrl)
		logs.EXPECT().Error(errGeneric, "encountered an error").Times(1)
		handler := NewHandler(logs)

		runFunc := func() error {
			return errGeneric
		}
		sync := handler.RunAsync(runFunc)
		assert.Equal(t, errGeneric, sync())
	})

	t.Run("Panic Logs Error Nil Error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logs := mock.NewMockLogger(mockCtrl)
		logs.EXPECT().Errorf("encountered a panic: %+v", gomock.Not(nil)).Times(1)
		handler := NewHandler(logs)

		runFunc := func() error {
			panic("")
		}
		done := make(chan bool)
		go func() {
			defer func() {
				done <- true
			}()
			defer handler.HandlePanic(nil)
			_ = runFunc()
		}()
		<-done
	})
}

func TestHandler_HandleError(t *testing.T) {
	t.Parallel()

	t.Run("Precondition Error Logged as Info", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockLogger := mock.NewMockLogger(mockCtrl)
		mockLogger.EXPECT().Infof("preconditions not met: error").Times(1)
		handler := NewHandler(mockLogger)

		sourceError := exceptions.PreconditionError(errGeneric.Error())
		err := handler.HandleError(sourceError)
		assert.Equal(t, sourceError, err)
	})

}
