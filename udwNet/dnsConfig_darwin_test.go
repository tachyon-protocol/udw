package udwNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestParseNetworksetupListnetworkserviceorder(ot *testing.T) {
	content := `An asterisk (*) denotes that a network service is disabled.
(1) SAMSUNG_Android 5
(Hardware Port: Modem (usbmodem14212), Device: usbmodem14212)

(2) Bluetooth DUN
(Hardware Port: Bluetooth DUN, Device: Bluetooth-Modem)

(3) Thunderbolt Ethernet
(Hardware Port: Thunderbolt Ethernet, Device: en4)

(4) Wi-Fi
(Hardware Port: Wi-Fi, Device: en0)

(5) config6
(Hardware Port: PPTP, Device: )

(6) Vpn1
(Hardware Port: IPSec, Device: )
`
	out := parseNetworksetupListnetworkserviceorder(content)
	udwTest.Equal(len(out), 4)
	udwTest.Equal(out["en0"], "Wi-Fi")
	udwTest.Equal(out["en4"], "Thunderbolt Ethernet")
	udwTest.Equal(out["usbmodem14212"], "SAMSUNG_Android 5")
	udwTest.Equal(out["Bluetooth-Modem"], "Bluetooth DUN")

}
