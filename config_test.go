package yesql

import "testing"

func TestQuietIf(t *testing.T) {
	testCases := []struct {
		name string
		opts []func(*Config)
		want bool
	}{
		{
			name: "TrueDisablesLogging",
			opts: []func(*Config){QuietIf(true)},
			want: true,
		},
		{
			name: "FalseEnablesLogging",
			opts: []func(*Config){QuietIf(false)},
			want: false,
		},
		{
			name: "FalseOverridesOptQuiet",
			opts: []func(*Config){OptQuiet(), QuietIf(false)},
			want: false,
		},
		{
			name: "OptQuietOverridesFalse",
			opts: []func(*Config){QuietIf(false), OptQuiet()},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NewConfig(tc.opts...).quiet; got != tc.want {
				t.Errorf("quiet = %t; want %t", got, tc.want)
			}
		})
	}
}
