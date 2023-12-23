package mb8600

import (
	"reflect"
	"testing"
)

func TestNewDownstreamChannelsFromResponse(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		args    args
		want    []*DownstreamChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDownstreamChannelsFromResponse(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDownstreamChannelsFromResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDownstreamChannelsFromResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDownstreamChannelFromLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *DownstreamChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDownstreamChannelFromLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDownstreamChannelFromLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDownstreamChannelFromLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUpstreamChannelsFromResponse(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		args    args
		want    []*UpstreamChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUpstreamChannelsFromResponse(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUpstreamChannelsFromResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpstreamChannelsFromResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUpstreamChannelFromLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *UpstreamChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUpstreamChannelFromLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUpstreamChannelFromLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpstreamChannelFromLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
