package chaos

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/mkuchenbecker/orcareduce/orcareduce/mock"
	"github.com/stretchr/testify/assert"
)

func TestRandomStaticErrors_100Percent(t *testing.T) {
	errGeneric := fmt.Errorf("error")
	errorInjector := NewRandomStaticErrorInjector(errGeneric, 1)
	for i := 0; i < 100; i++ {
		assert.Equal(t, errGeneric, errorInjector.Error())
	}
}

func TestRandomStaticErrors_ZeroPercent(t *testing.T) {
	errGeneric := fmt.Errorf("error")
	errorInjector := NewRandomStaticErrorInjector(errGeneric, 0)
	for i := 0; i < 100; i++ {
		assert.Equal(t, nil, errorInjector.Error())
	}
}

func TestMetaErrors(t *testing.T) {
	errGeneric := fmt.Errorf("error")
	t.Parallel()
	rand.Seed(1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockError0 := mock.NewMockErrorInjector(mockCtrl)
	mockError0.EXPECT().Error().Return(nil).Times(1)
	mockError1 := mock.NewMockErrorInjector(mockCtrl)
	mockError1.EXPECT().Error().Return(errGeneric).Times(1)

	meta := NewMetaErrorInjector([]orcareduce.ErrorInjector{mockError0, mockError1})
	assert.Equal(t, errGeneric, meta.Error())
}
