package udwSys

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMemory(t *testing.T) {

}

func TestNetworkConnection(t *testing.T) {
	count := networkConnection(`
	39`)
	udwTest.Equal(count, 39)
}

func TestIKEUserCount(t *testing.T) {
	c := ikeUserCount(`uptime: 45 minutes, since Jul 17 16:03:02 2015
worker threads: 32 total, 27 idle, working: 4/0/1/0
job queues: 0/0/0/0
jobs scheduled: 471
IKE_SAs: 79 total, 0 half-open
mallinfo: sbrk 7364608, mmap 0, used 2095152, free 5269456`)
	udwTest.Equal(c, 79)
}
