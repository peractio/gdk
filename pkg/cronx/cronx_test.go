package cronx

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestEvery(t *testing.T) {
	type args struct {
		duration time.Duration
		job      cron.Job
		mock     func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Uninitialized",
			args: args{
				duration: 0,
				job:      nil,
				mock: func() {
					Default()
					commandController.Commander = nil
				},
			},
		},
		{
			name: "Success",
			args: args{
				duration: 5 * time.Minute,
				job:      Func(func() {}),
				mock: func() {
					Default()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mock()
			Every(tt.args.duration, tt.args.job)
		})
	}
}

func TestFunc_Run(t *testing.T) {
	tests := []struct {
		name string
		r    Func
	}{
		{
			name: "Success",
			r:    Func(func() {}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Run()
		})
	}
}

func TestGetEntries(t *testing.T) {
	tests := []struct {
		name string
		mock func()
		want bool
	}{
		{
			name: "Uninitialized",
			mock: func() {
				Default()
				commandController.Commander = nil
			},
		},
		{
			name: "Success",
			mock: func() {
				Default()
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got := GetEntries()
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		id   cron.EntryID
		mock func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Uninitialized",
			args: args{
				id: 1,
				mock: func() {
					Default()
					commandController.Commander = nil
				},
			},
		},
		{
			name: "Success",
			args: args{
				id: 1,
				mock: func() {
					Default()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mock()
			Remove(tt.args.id)
		})
	}
}

func TestSchedule(t *testing.T) {
	type args struct {
		spec string
		job  cron.Job
		mock func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Uninitialized",
			args: args{
				spec: "@every 5m",
				job:  Func(func() {}),
				mock: func() {
					Default()
					commandController.Commander = nil
				},
			},
			wantErr: true,
		},
		{
			name: "Broken spec",
			args: args{
				spec: "this is clearly not a spec",
				job:  Func(func() {}),
				mock: func() {
					Default()
				},
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				spec: "@every 5m",
				job:  Func(func() {}),
				mock: func() {
					Default()
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mock()
			if err := Schedule(tt.args.spec, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Schedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Default()
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				config: Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.args.config)
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "Uninitialized",
			mock: func() {
				Default()
				commandController.Commander = nil
			},
		},
		{
			name: "Success",
			mock: func() {
				Default()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			Stop()
		})
	}
}