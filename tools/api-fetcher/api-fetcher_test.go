package apifetcher_test

import (
	"testing"
	"time"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage/postgres"
	apifetcher "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/tools/api-fetcher"
)

func NewFetcher_Test()

func TestNewFetcher(t *testing.T) {
	type args struct {
		s      *postgres.Storage
		period time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apifetcher.NewFetcher(tt.args.s, tt.args.period)
		})
	}
}
