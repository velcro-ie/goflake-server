package goflake

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	instance *uuidGenerator
	once     sync.Once
)

// GetNonLocalInterface returns the first interface that isn't a loopback.
func GetNonLocalInterface() (net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback == 0 {
			return iface, nil
		}
	}
	return net.Interface{}, errors.New("Unable to determine interface that isn't a loopback")

}

// putUInt Puts the lower numberOfBytes from longValue into the slice, starting at index pos.
func putUInt(byteArray []byte, longValue uint64, pos int, numberOfBytes int) {
	for i := 0; i < numberOfBytes; i++ {
		val := byte(longValue >> uint(i*8))
		byteArray[pos+numberOfBytes-i-1] = val
	}
}

// A uuidGenerator that allows for the generation for UUIDs
type uuidGenerator struct {
	sequence      uint64
	lastTimestamp timestamp
	addr          []byte
}

func goFlakeInstance(unique []byte) *uuidGenerator {
	if len(unique) > 6 {
		panic("Unique value must be byte array of 6 of less")
	}
	once.Do(func() {
		var random = rand.New(rand.NewSource(time.Now().UnixNano()))
		instance = &uuidGenerator{uint64(random.Int63()), timestamp{}, unique}
	})
	return instance
}

func GoFlakeInstanceUsingUnique(unique string) *uuidGenerator {
	return goFlakeInstance([]byte(unique))
}

// GoFlakeInstanceUsingMacAddress returns A singleton instance of the uuidGenerator
func GoFlakeInstanceUsingMacAddress() *uuidGenerator {
	localInterface, err := GetNonLocalInterface()
	if err != nil {
		panic("Failed to create uuidGenerator")
	}
	return goFlakeInstance(localInterface.HardwareAddr)
}

// GetBase64UUID returns A UUID encoded to Base64 URL safe
func (this *uuidGenerator) GetBase64UUID() string {

	mySequence := atomic.AddUint64(&this.sequence, 1) & 0xffffff
	flakeId := make([]byte, 15)
	myTimestamp := this.lastTimestamp.pushTimestamp(mySequence)
	putUInt(flakeId, myTimestamp, 0, 6)
	copy(flakeId[6:12], this.addr[:])
	putUInt(flakeId, mySequence, 12, 3)
	return base64.RawURLEncoding.EncodeToString(flakeId)

}

// The timestamp struct holds the last timestamp
type timestamp struct {
	lastTimeStamp uint64
	sync.Mutex
}

// pushTimestamp returns the last timestamp ensuring that the timestamp doesn't move backwards
func (this *timestamp) pushTimestamp(sequence uint64) uint64 {
	timestampHolder := uint64(time.Now().UnixNano()) / 1000000
	this.Lock()
	defer this.Unlock()
	this.lastTimeStamp = max(timestampHolder, this.lastTimeStamp)
	if sequence == 0 {
		this.lastTimeStamp++
	}
	return this.lastTimeStamp
}

// max returns the max value of two uint64 values
func max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
