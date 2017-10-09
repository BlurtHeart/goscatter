package test

import (
	"github.com/blurtheart/goscatter/testmock/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCompany_Meeting(t *testing.T) {
	ctl := gomock.NewController(t)
	mock_talker := mock_test.NewMockTalker(ctl)
	mock_talker.EXPECT().Talk(gomock.Eq("Wang")).Return("something random")

	company := NewCompany(mock_talker)
	t.Log(company.Meeting("Wang"))
}
