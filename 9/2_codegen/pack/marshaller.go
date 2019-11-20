package main

import "encoding/binary"
import "bytes"

func (in *User) Unpack(data []byte) error {
	r := bytes.NewReader(data)

	// ID
	var IDRaw uint32
	binary.Read(r, binary.LittleEndian, &IDRaw)
	in.ID = int(IDRaw)

	// Name
	var NameLenRaw uint32
	binary.Read(r, binary.LittleEndian, &NameLenRaw)
	NameRaw := make([]byte, NameLenRaw)
	binary.Read(r, binary.LittleEndian, &NameRaw)
	in.Name = string(NameRaw)

	// Flags
	var FlagsRaw uint32
	binary.Read(r, binary.LittleEndian, &FlagsRaw)
	in.Flags = int(FlagsRaw)
	return nil
}

