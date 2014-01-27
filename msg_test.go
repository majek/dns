// Copyright 2011 Miek Gieben. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dns

import (
	"testing"
	"encoding/hex"
)

var testMessages = []string{
	// google
	"064e81800001000b0004000506676f6f676c6503636f6d0000010001c00c00010001000000050004adc22986c00c00010001000000050004adc22987c00c00010001000000050004adc22988c00c00010001000000050004adc22989c00c00010001000000050004adc2298ec00c00010001000000050004adc22980c00c00010001000000050004adc22981c00c00010001000000050004adc22982c00c00010001000000050004adc22983c00c00010001000000050004adc22984c00c00010001000000050004adc22985c00c00020001000000050006036e7331c00cc00c00020001000000050006036e7332c00cc00c00020001000000050006036e7333c00cc00c00020001000000050006036e7334c00cc0d800010001000000050004d8ef200ac0ea00010001000000050004d8ef220ac0fc00010001000000050004d8ef240ac10e00010001000000050004d8ef260a0000290500000000050000",
	// amazon
	"6de1818000010004000a000806616d617a6f6e03636f6d0000010001c00c000100010000000500044815c2d4c00c000100010000000500044815d7e8c00c00010001000000050004b02062a6c00c00010001000000050004cdfbf236c00c000200010000000500140570646e733408756c747261646e73036f726700c00c000200010000000500150570646e733508756c747261646e7304696e666f00c00c000200010000000500160570646e733608756c747261646e7302636f02756b00c00c00020001000000050014036e7331037033310664796e656374036e657400c00c00020001000000050006036e7332c0cfc00c00020001000000050006036e7333c0cfc00c00020001000000050006036e7334c0cfc00c000200010000000500110570646e733108756c747261646e73c0dac00c000200010000000500080570646e7332c127c00c000200010000000500080570646e7333c06ec0cb00010001000000050004d04e461fc0eb00010001000000050004cc0dfa1fc0fd00010001000000050004d04e471fc10f00010001000000050004cc0dfb1fc12100010001000000050004cc4a6c01c121001c000100000005001020010502f3ff00000000000000000001c13e00010001000000050004cc4a6d01c13e001c0001000000050010261000a1101400000000000000000001",
	// yahoo
	"fc2d81800001000300070008057961686f6f03636f6d0000010001c00c00010001000000050004628afd6dc00c00010001000000050004628bb718c00c00010001000000050004cebe242dc00c00020001000000050006036e7336c00cc00c00020001000000050006036e7338c00cc00c00020001000000050006036e7331c00cc00c00020001000000050006036e7332c00cc00c00020001000000050006036e7333c00cc00c00020001000000050006036e7334c00cc00c00020001000000050006036e7335c00cc07b0001000100000005000444b48310c08d00010001000000050004448eff10c09f00010001000000050004cb54dd35c0b100010001000000050004628a0b9dc0c30001000100000005000477a0f77cc05700010001000000050004ca2bdfaac06900010001000000050004caa568160000290500000000050000",
	// microsoft.com
	"f4368180000100020005000b096d6963726f736f667403636f6d0000010001c00c0001000100000005000440040b25c00c0001000100000005000441373ac9c00c0002000100000005000e036e7331046d736674036e657400c00c00020001000000050006036e7332c04fc00c00020001000000050006036e7333c04fc00c00020001000000050006036e7334c04fc00c00020001000000050006036e7335c04fc04b000100010000000500044137253ec04b001c00010000000500102a010111200500000000000000010001c0650001000100000005000440043badc065001c00010000000500102a010111200600060000000000010001c07700010001000000050004d5c7b435c077001c00010000000500102a010111202000000000000000010001c08900010001000000050004cf2e4bfec089001c00010000000500102404f800200300000000000000010001c09b000100010000000500044137e28cc09b001c00010000000500102a010111200f000100000000000100010000290500000000050000",
	// MX google.com
	"724b8180000100050004000b06676f6f676c6503636f6d00000f0001c00c000f000100000005000c000a056173706d78016cc00cc00c000f0001000000050009001404616c7431c02ac00c000f0001000000050009001e04616c7432c02ac00c000f0001000000050009002804616c7433c02ac00c000f0001000000050009003204616c7434c02ac00c00020001000000050006036e7332c00cc00c00020001000000050006036e7333c00cc00c00020001000000050006036e7334c00cc00c00020001000000050006036e7331c00cc02a00010001000000050004adc2421bc02a001c00010000000500102a00145040080c01000000000000001bc04200010001000000050004adc2461bc05700010001000000050004adc2451bc06c000100010000000500044a7d8f1bc081000100010000000500044a7d191bc0ca00010001000000050004d8ef200ac09400010001000000050004d8ef220ac0a600010001000000050004d8ef240ac0b800010001000000050004d8ef260a0000290500000000050000",
	// reddit.com
	"12b98180000100080000000c0672656464697403636f6d0000020001c00c0002000100000005000f046175733204616b616d036e657400c00c000200010000000500070475736534c02dc00c000200010000000500070475737733c02dc00c000200010000000500070475737735c02dc00c00020001000000050008056173696131c02dc00c00020001000000050008056173696139c02dc00c00020001000000050008056e73312d31c02dc00c0002000100000005000a076e73312d313935c02dc02800010001000000050004c30a242ec04300010001000000050004451f1d39c05600010001000000050004451f3bc7c0690001000100000005000460073240c07c000100010000000500046007fb81c090000100010000000500047c283484c090001c00010000000500102a0226f0006700000000000000000064c0a400010001000000050004c16c5b01c0a4001c000100000005001026001401000200000000000000000001c0b800010001000000050004c16c5bc3c0b8001c0001000000050010260014010002000000000000000000c30000290500000000050000",
}

func TestMsgLenTest2(t *testing.T) {
	for i, hexData := range testMessages {
		input, err := hex.DecodeString(hexData)
		if err != nil {
			panic("")
		}
		m := new(Msg)
		if m.Unpack(input) != nil {
			t.Errorf("Unpack failure on input %d: %v", i, hexData)
			continue;
		}
		m.Compress = true
		lenComp := m.Len()
		b, err := m.Pack()
		if err != nil {
			panic("")
		}
		pacComp := len(b)
		m.Compress = false
		lenUnComp := m.Len()
		b, err = m.Pack()
		if err != nil {
			panic("")
		}
		pacUnComp := len(b)
		if pacComp != lenComp {
			t.Errorf("msg.Len(compressed)=%d actual=%d for test %d: %v",
				lenComp, pacComp, i, hexData)
		}

		if pacUnComp != lenUnComp {
			t.Errorf("msg.Len(uncompressed)=%d actual=%d for test %d: %v",
				lenUnComp, pacUnComp, i, hexData)
		}

		if pacUnComp != m.packLen() {
			t.Errorf("msg.packLength()=%d actual=%d for test %d: %v",
				m.packLen(), pacUnComp, i, hexData)
		}
	}
}
