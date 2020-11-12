package bin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type SlabUninitialized struct {
}

type SlabInnerNode struct {
	//    u32('prefixLen'),
	//    u128('key'),
	//    seq(u32(), 2, 'children'),
	PrefixLen uint32
	Key       Uint128
	Children  [2]uint32
}

type SlabLeafNode struct {
	OwnerSlot     uint8
	FeeTier       uint8
	Padding       [2]byte `json:"-"`
	Key           Uint128
	Owner         PublicKey
	Quantity      Uint64
	ClientOrderId Uint64
}
type SlabFreeNode struct {
	Next uint32
}

type SlabLastFreeNode struct {
}

type PublicKey [32]byte

var SlabFactoryImplDef = NewVariantDefinition(Uint32TypeIDEncoding, []VariantType{
	{"uninitialized", (*SlabUninitialized)(nil)},
	{"inner_node", (*SlabInnerNode)(nil)},
	{"leaf_node", (*SlabLeafNode)(nil)},
	{"free_node", (*SlabFreeNode)(nil)},
	{"last_free_node", (*SlabLastFreeNode)(nil)},
})

type Slab struct {
	BaseVariant
}

func (s *Slab) UnmarshalBinary(decoder *Decoder) error {
	return s.BaseVariant.UnmarshalBinaryVariant(decoder, SlabFactoryImplDef)
}
func (s *Slab) MarshalBinary(encoder *Encoder) error {
	err := encoder.writeUint16(uint16(s.TypeID))
	if err != nil {
		return err
	}
	return encoder.Encode(s.Impl)
}

type Orderbook struct {
	// ORDERBOOK_LAYOUT
	SerumPadding [5]byte `json:"-"`
	AccountFlags uint64
	// SLAB_LAYOUT
	// SLAB_HEADER_LAYOUT
	BumpIndex    uint32  `bin:"sizeof=Nodes"`
	ZeroPaddingA [4]byte `json:"-"`
	FreeListLen  uint32
	ZeroPaddingB [4]byte `json:"-"`
	FreeListHead uint32
	Root         uint32
	LeafCount    uint32
	ZeroPaddingC [4]byte `json:"-"`

	// SLAB_NODE_LAYOUT
	Nodes []*Slab `bin: ""`
}

func TestDecoder_DecodeSol(t *testing.T) {

	hexData := `736572756d2100000000000000650000000000000044000000000000004c0000000000000011000000000000000100000035000000010babffffffffff4105000000000000400000003f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000007060000e92ba8ffffffffff6a050000000000005b4388a3431832af5742b863e200b8e733ce451f27006e2723d9198b363355d18813000000000000c6a1933c53bc44160200000003060000ad2ba8ffffffffff6b050000000000005b4388a3431832af5742b863e200b8e733ce451f27006e2723d9198b363355d18813000000000000f49ba2465bbc441603000000340000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000380000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000140600003334a8ffffffffff66050000000000005ae01b52d00a090c6dc6fce8e37a225815cff2223a99c6dfdad5aae56d3db6705246000000000000c837c70ea5737ca70200000013060000e234a8ffffffffff67050000000000005ae01b52d00a090c6dc6fce8e37a225815cff2223a99c6dfdad5aae56d3db670e367000000000000b1223a1fe53f3899030000006200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030000000900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030000003200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000003d000000965ca8ffffffffff600500`
	cnt, err := hex.DecodeString(hexData)
	require.NoError(t, err)

	decoder := NewDecoder(cnt)
	var ob *Orderbook
	err = decoder.Decode(&ob)
	require.NoError(t, err)

	//json, err := json.MarshalIndent(ob, "", "   ")
	//require.NoError(t, err)
	//fmt.Println(string(json))

	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf)
	err = encoder.Encode(ob)
	require.NoError(t, err)

	obHex := hex.EncodeToString(buf.Bytes())
	require.Equal(t, hexData, obHex)

}

func TestDecoder_Slabs(t *testing.T) {

	//zlog, _ := zap.NewDevelopment()
	//EnableDebugLogging(zlog)

	rawSlabs := []string{
		"0100000035000000010babffffffffff4105000000000000400000003f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0200000014060000b2cea5ffffffffff23070000000000005ae01b52d00a090c6dc6fce8e37a225815cff2223a99c6dfdad5aae56d3db670e62c000000000000140b0fadcf8fcebf",
		"030000003400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	}

	for _, s := range rawSlabs {
		cnt, err := hex.DecodeString(s)
		require.NoError(t, err)

		decoder := NewDecoder(cnt)
		var slab *Slab
		err = decoder.Decode(&slab)
		require.NoError(t, err)

		json, err := json.MarshalIndent(slab, "", "   ")
		require.NoError(t, err)
		fmt.Println(string(json))

		//require.Equal(t, 0, decoder.remaining())

	}
}
