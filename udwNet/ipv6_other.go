// +build !ios,!macAppStore

package udwNet

func NewSupportIpv6OnlyDialer(oldDialer Dialer) Dialer {
	return oldDialer
}
