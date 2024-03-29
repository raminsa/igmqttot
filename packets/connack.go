/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Allan Stockdill-Mander
 */

package packets

import (
	"bytes"
	"fmt"
	"io"
)

// ConnackPacket is an internal representation of the fields of the
// Connack MQTT packet
type ConnackPacket struct {
	FixedHeader
	SessionPresent bool
	ReturnCode     byte
	AckFlags       byte
	Payload        []byte
}

func (ca *ConnackPacket) String() string {
	return fmt.Sprintf("%s sessionpresent: %t returncode: %d", ca.FixedHeader, ca.SessionPresent, ca.ReturnCode)
}

func (ca *ConnackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.WriteByte(boolToByte(ca.SessionPresent))
	body.WriteByte(ca.ReturnCode)
	ca.FixedHeader.RemainingLength = 2
	packet := ca.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (ca *ConnackPacket) Unpack(b io.Reader) error {
	flags, err := decodeByte(b)
	if err != nil {
		return err
	}
	ca.AckFlags = flags
	ca.SessionPresent = 1&flags > 0
	ca.ReturnCode, err = decodeByte(b)
	if err != nil {
		return err
	}
	if ca.RemainingLength > 2 {
		ca.Payload = make([]byte, ca.RemainingLength-2)
		_, err = io.ReadFull(b, ca.Payload)
		if err != nil {
			return err
		}
	} else {
		ca.Payload = []byte{}
	}
	return err
}

// Details returns a Details struct containing the Qos and
// MessageID of this ControlPacket
func (ca *ConnackPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}
