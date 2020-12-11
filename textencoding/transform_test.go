package textencoding

import (
	"reflect"
	"testing"
)

func TestTransform(t *testing.T) {
	type args struct {
		s    []byte
		from string
		to   string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"error1",
			args{
				[]byte("中文"),
				"unknown",
				"gbk",
			},
			nil,
			true,
		},
		{
			"error1",
			args{
				[]byte("中文"),
				"utf8",
				"unknown",
			},
			nil,
			true,
		},
		{
			"utf8 -> gbk",
			args{
				[]byte("中文"),
				"UTF8",
				"gbk",
			},
			[]byte{0xD6, 0xD0, 0xCE, 0xC4},
			false,
		},
		{
			"gbk -> utf8",
			args{
				[]byte{0xD6, 0xD0, 0xCE, 0xC4},
				"gbk",
				"UTF8",
			},
			[]byte("中文"),
			false,
		},
		{
			"gbk -> Big5",
			args{
				[]byte{0xD6, 0xD0, 0xCE, 0xC4},
				"gbk",
				"Big5",
			},
			[]byte{0xA4, 0xA4, 0xA4, 0xE5},
			false,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got, err := Transform(tt.args.s, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		s  []byte
		to string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"",
			args{
				[]byte("中文"),
				"gbk",
			},
			[]byte{0xD6, 0xD0, 0xCE, 0xC4},
			false,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.s, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		s    []byte
		from string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"",
			args{
				[]byte{0xD6, 0xD0, 0xCE, 0xC4},
				"gbk",
			},
			[]byte("中文"),
			false,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.s, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
