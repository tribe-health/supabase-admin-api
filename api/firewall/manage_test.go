package firewall

import (
	"net"
	"reflect"
	"testing"
)

func TestFirewall_Process(t *testing.T) {
	type fields struct {
		PrivilegedPorts          []int
		CustomerPorts            []int
		PrivilegedPortsWhitelist []net.IPNet
	}
	type args struct {
		userNetworks []net.IPNet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "normalOperation",
			fields: fields{PrivilegedPorts: []int{22}, CustomerPorts: []int{5432, 6543}, PrivilegedPortsWhitelist: getNetworks([]string{"54.45.65.123/32"}, t)},
			args:   args{getNetworks([]string{"20.34.0.0/16"}, t)},
			want: `flush table inet supabase_managed

define SUPABASE_INTERNAL = {
       10.0.0.0/8,
       172.16.0.0/12,
       192.168.0.0/16
}

define USER_WHITELIST = {
       20.34.0.0/16
}

define TEMP_PRIV = {
       54.45.65.123/32
}

define CUSTOMER_PORTS = {
       5432, 6543
}

define PRIVILEGED_PORTS = {
       22
}

table inet supabase_managed {
    chain inbound {
        type filter hook input priority 0; policy drop;
        ct state vmap { established : accept, related : accept, invalid : drop }
        iifname lo accept
        tcp dport $CUSTOMER_PORTS jump customer_traffic
        tcp dport $PRIVILEGED_PORTS jump filter_privileged_ports
    }

    chain customer_traffic {
        ip saddr $USER_WHITELIST accept
        ip saddr $SUPABASE_INTERNAL accept
    }

    chain filter_privileged_ports {
        ip saddr $SUPABASE_INTERNAL accept
        ip saddr $TEMP_PRIV accept
    }

    chain forward {
        type filter hook forward priority 0; policy drop;
    }
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Manager{
				PrivilegedPorts:          tt.fields.PrivilegedPorts,
				CustomerPorts:            tt.fields.CustomerPorts,
				PrivilegedPortsWhitelist: tt.fields.PrivilegedPortsWhitelist,
			}
			got, err := f.Process(tt.args.userNetworks)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getNetworks(nets []string, t *testing.T) []net.IPNet {
	out := make([]net.IPNet, 0)
	for _, subnet := range nets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			t.Fatal("failed to parse network", err)
		}
		out = append(out, *ipNet)
	}
	return out
}
