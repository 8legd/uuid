/**
 * Date: 31/01/14
 * Time: 3:34 PM
 */
package uuid

import (
	"fmt"
	"net"
)

// UUIDStruct is used for RFC4122 Version 1 UUIDs
type UUIDStruct struct {
	timeLow                   uint32
	timeMid                   uint16
	timeHiAndVersion          uint16
	sequenceHiAndVariant      byte
	sequenceLow               byte
	node                      []byte
	size                      int
}

func (o UUIDStruct) Size() int {
	return o.size
}

func (o UUIDStruct) Version() int {
	return int(o.timeHiAndVersion>>12)
}

func (o UUIDStruct) Variant() byte {
	return variant(o.sequenceHiAndVariant)
}

// Sets the four most significant bits (bits 12 through 15) of the
// timeHiAndVersion field to the 4-bit version number.
func (o *UUIDStruct) setVersion(pVersion int) {
	o.timeHiAndVersion &= 0x0FFF
	o.timeHiAndVersion |= (uint16(pVersion)<<12)
}

func (o *UUIDStruct) setVariant(pVariant byte) {
	setVariant(&o.sequenceHiAndVariant, pVariant)
}

func (o *UUIDStruct) Unmarshal(pData []byte) {
	o.timeLow = uint32(pData[3]) | uint32(pData[2])<<8 | uint32(pData[1])<<16 | uint32(pData[0])<<24
	o.timeMid = uint16(pData[5]) | uint16(pData[4])<<8
	o.timeHiAndVersion = uint16(pData[7]) | uint16(pData[6])<<8
	o.sequenceHiAndVariant = pData[8]
	o.sequenceLow = pData[9]
	o.node = pData[10:o.Size()]
}

func (o *UUIDStruct) Bytes() (data []byte) {
	data = []byte{
		byte(o.timeLow>>24), byte(o.timeLow>>16), byte(o.timeLow>>8), byte(o.timeLow),
		byte(o.timeMid>>8), byte(o.timeMid),
		byte(o.timeHiAndVersion>>8), byte(o.timeHiAndVersion),
	}
	data = append(data, o.sequenceHiAndVariant)
	data = append(data, o.sequenceLow)
	data = append(data, o.node...)
	return
}

// Marshals the UUID bytes into a slice
func (r UUIDStruct) MarshalBinary() ([]byte, error) {
	return r.Bytes(), nil
}

// Un-marshals the data bytes into the UUID struct.
// Implements the BinaryUn-marshaller interface
func (o *UUIDStruct) UnmarshalBinary(pData []byte) error {
	return UnmarshalBinary(o, pData)
}

func (o UUIDStruct) String() string {
	b := o.Bytes()
	return fmt.Sprintf(format, b[0:4], b[4:6], b[6:8], b[8:9], b[9:10], b[10:o.Size()])
}

// Set the three most significant bits (bits 0, 1 and 2) of the
// sequenceHiAndVariant to variant mask 0x80.
func (o *UUIDStruct) setRFC4122Variant() {
	o.sequenceHiAndVariant &= variantSet
	o.sequenceHiAndVariant |= ReservedRFC4122
}

// Unmarshals data into struct for V1 UUIDs
func newV1(pNow Timestamp, pVersion uint16, pVariant byte, pNode net.HardwareAddr) UUID {
	o := new(UUIDStruct)
	o.timeLow = uint32(pNow & 0xFFFFFFFF)
	o.timeMid = uint16((pNow>>32) & 0xFFFF)
	o.timeHiAndVersion = uint16((pNow>>48) & 0x0FFF)
	o.timeHiAndVersion |= uint16(pVersion<<12)
	o.sequenceLow = byte(state.sequence & 0xFF)
	o.sequenceHiAndVariant = byte(( state.sequence & 0x3F00)>>8)
	o.sequenceHiAndVariant |= pVariant
	o.node = pNode
	return o
}
