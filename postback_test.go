package tracking

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostbackProvider_Request(t *testing.T) {
	t.Parallel()

	p := NewPostbackProvider("example.com", true)
	ip := net.IPv4(5, 6, 7, 8)

	s, err := http.NewRequest(http.MethodGet, "https://example.com/postback?click_id=111111111111111111111111", nil)
	require.NoError(t, err)

	f, err := http.NewRequest(http.MethodGet, "https://example.com/postback?action_id=2&click_id=111111111111111111111111&comment=10&custom_field1=16&custom_field10=25&custom_field11=26&custom_field12=27&custom_field13=28&custom_field14=29&custom_field15=30&custom_field2=17&custom_field3=18&custom_field4=19&custom_field5=20&custom_field6=21&custom_field7=22&custom_field8=23&custom_field9=24&device_type=13&fbclid=12&goal=3&ip=5.6.7.8&referrer=9&secure=11&status=1&sum=4&user_id=14", nil)
	require.NoError(t, err)

	for _, tc := range []struct {
		name string
		pb   *Postback
		exp  *http.Request
		err  error
	}{
		{
			name: "error",
			pb:   &Postback{ClickID: "1"},
			exp:  nil,
			err:  ErrInvalidClickID,
		},
		{
			name: "simple",
			pb:   &Postback{ClickID: "111111111111111111111111"},
			exp:  s,
			err:  nil,
		},
		{
			name: "full",
			pb: &Postback{
				ClickID:    "111111111111111111111111",
				ActionID:   "2",
				Goal:       "3",
				Sum:        float64Ptr(4),
				IP:         &ip,
				Status:     PostbackStatusConfirmed,
				Referrer:   "9",
				Comment:    "10",
				Secure:     "11",
				FbClID:     "12",
				DeviceType: "13",
				UserID:     "14",
				CustomFields: [15]string{
					"16", "17", "18", "19", "20",
					"21", "22", "23", "24", "25",
					"26", "27", "28", "29", "30",
				},
			},
			exp: f,
			err: nil,
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, err := p.Request(context.Background(), tc.pb)
			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.exp, r)
		})
	}
}

func float64Ptr(x float64) *float64 {
	return &x
}
