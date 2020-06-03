package main





import (
	"fmt"
	"time"
	"encoding/binary"
	"math/rand" 
)



func main() {
	fmt.Printf("A test harness for byte manipulation\n")


		rand.Seed(time.Now().UnixNano())

	// date is transmitted as signed 64 bit signed integer it is tv_sec of struct timeval


	tvsec := time.Now().Unix()
	tvsec32 := int32(tvsec);
	tvsecu32 := uint32(tvsec32)

	fmt.Printf("time: %d - %d - %d\n", tvsec, tvsec32, tvsecu32)


	// encode 
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf,tvsecu32)
	fmt.Printf("%v\n", buf)


	// decode
	tvsecu32 = binary.BigEndian.Uint32(buf)
	tvsec32 = int32(tvsecu32)
	tvsec = int64(tvsec32)

	
	fmt.Printf("time: %d - %d - %d\n", tvsec, tvsec32, tvsecu32)

	t := time.Now()

	fmt.Println( t.Format(time.RFC3339))
	

// 	func main() {
//     buf := make([]byte, 10)
//     ts := uint32(time.Now().Unix())
//     binary.BigEndian.PutUint16(buf[0:], 0xa20c) // sensorID
//     binary.BigEndian.PutUint16(buf[2:], 0x04af) // locationID
//     binary.BigEndian.PutUint32(buf[4:], ts)     // timestamp
//     binary.BigEndian.PutUint16(buf[8:], 479)    // temp 
//     fmt.Printf("% x\n", buf)
// }

	b := []byte{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18}
	
	min := 4
	max := len(b) - 4
	split := rand.Intn(max - min + 1) + min


	
	ba := b[0:split]
	bb := b[split:]


	fmt.Printf("%v\n%v\n%v\n",b,ba,bb)
	
}
