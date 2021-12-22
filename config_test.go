package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNewAria2(t *testing.T) {
	type args struct {
		id     string
		links  string
		secret string
	}
	tests := []struct {
		name string
		args args
		want *Aria2
	}{
		{name: "xx", args: args{id: "id", links: "links", secret: "secret"}, want: &Aria2{Id: "id", Params: []interface{}{"token:secret", []string{"links"}}, Jsonrpc: "2.0", Method: "aria2.addUri"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAria2(tt.args.id, tt.args.links, tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				b, _ := json.Marshal(got)
				fmt.Println(string(b))

				t.Errorf("NewAria2() = %v, want %v", got, tt.want)
			}
		})
	}
}
