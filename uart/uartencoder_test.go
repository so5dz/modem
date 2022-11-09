package uart

import (
	"reflect"
	"testing"
)

func TestUARTSymbolizer_Symbolize(t *testing.T) {
	type fields struct {
		Start  bool
		Data   int
		Parity Parity
		Stop   bool
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Symbol
	}{
		{
			name:   "no bytes",
			fields: fields{Start: true, Data: 1234, Parity: Even, Stop: true},
			args:   args{bytes: []byte{}},
			want:   []Symbol{},
		},
		{
			name:   "single word 5-N-1",
			fields: fields{Start: true, Data: 5, Parity: None, Stop: true},
			args:   args{bytes: []byte{0b10100}},
			want:   []Symbol{Start, Mark, Space, Mark, Space, Space, Stop},
		},
		{
			name:   "multiple words 4-N-1",
			fields: fields{Start: false, Data: 4, Parity: None, Stop: true},
			args:   args{bytes: []byte{0b1010, 0b0001}},
			want:   []Symbol{Start, Mark, Space, Mark, Space, Stop, Start, Space, Space, Space, Mark, Stop},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UARTEncoder{
				Data:   tt.fields.Data,
				Parity: tt.fields.Parity,
				Stop:   tt.fields.Stop,
			}
			if got := s.Symbolize(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UARTSymbolizer.Symbolize() = %v, want %v", got, tt.want)
			}
		})
	}
}
