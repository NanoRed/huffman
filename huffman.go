package huffman

import (
	"bytes"
	"sort"
)

type Node struct {
	Left   *Node
	Right  *Node
	Weight int
	Value  interface{}
}

func (n *Node) traverse(depth int, code []byte, offset int, encoder *Encoder, decoder *Decoder) {
	i := depth / 8
	o := depth % 8
	if n.Left.Left == nil {
		ncode := make([]byte, len(code), len(code)+1)
		copy(ncode, code)
		if o == 7 {
			ncode = append(ncode, 1)
		} else {
			ncode[i] |= 1 << (o + 1)
		}
		(*encoder)[n.Left.Value] = ncode
		(*decoder)[string(ncode)] = n.Left.Value
	} else {
		defer func() {
			if o == 7 {
				code = append(code, 0)
			}
			n.Left.traverse(depth+1, code, o, encoder, decoder)
		}()
	}
	if n.Right.Left == nil {
		ncode := make([]byte, len(code), len(code)+1)
		copy(ncode, code)
		if o == 7 {
			ncode[i] |= 1 << o
			ncode = append(ncode, 1)
		} else {
			ncode[i] |= 3 << o
		}
		(*encoder)[n.Right.Value] = ncode
		(*decoder)[string(ncode)] = n.Right.Value
	} else {
		defer func() {
			if o == 7 {
				code = append(code, 0)
			}
			code[i] |= 1 << o
			n.Right.traverse(depth+1, code, o, encoder, decoder)
		}()
	}
}

func (n *Node) NewCoder() (Encoder, Decoder) {
	encoder := make(Encoder)
	decoder := make(Decoder)
	if n.Left == nil {
		encoder[n.Value] = []byte{2}
		decoder[string(encoder[n.Value])] = n.Value
	} else {
		n.traverse(0, []byte{0}, 0, &encoder, &decoder)
	}
	return encoder, decoder
}

type Leaves []*Node

func (l Leaves) Len() int           { return len(l) }
func (l Leaves) Less(i, j int) bool { return l[i].Weight < l[j].Weight }
func (l Leaves) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func (l Leaves) Build() *Node {
	sort.Stable(l)
	return l.BuildSorted()
}

func (l Leaves) BuildSorted() *Node {
	if len(l) == 0 {
		return nil
	}

	for len(l) > 1 {
		parent := &Node{
			Weight: l[0].Weight + l[1].Weight,
		}
		if l[1].Weight >= l[0].Weight {
			parent.Left, parent.Right = l[0], l[1]
		} else {
			parent.Right, parent.Left = l[0], l[1]
		}
		l[1] = parent
		l = l[1:]
	}

	return l[0]
}

type Encoder map[interface{}][]byte

// make sure all the values are contained in the encoder
func (e Encoder) EncodeStringSlice(values []string) []byte {
	b := &bytes.Buffer{}
	var remain byte
	var offset byte = 1
	for _, val := range values {
		i := 0
		for ; i < len(e[val])-1; i++ {
			var ofs byte = 1
			for k := 0; k < 8; k++ {
				if e[val][i]&ofs == ofs {
					remain |= offset
				}
				if offset == 128 {
					b.WriteByte(remain)
					remain = 0
					offset = 1
				} else {
					offset = offset << 1
				}
				ofs = ofs << 1
			}
		}
		end := 0
		var ofs byte = 128
		for k := 0; k < 8; k++ {
			if e[val][i]&ofs != ofs {
				ofs = ofs >> 1
				continue
			}
			end = 7 - k
			break
		}
		ofs = 1
		for k := 0; k < end; k++ {
			if e[val][i]&ofs == ofs {
				remain |= offset
			}
			if offset == 128 {
				b.WriteByte(remain)
				remain = 0
				offset = 1
			} else {
				offset = offset << 1
			}
			ofs = ofs << 1
		}
	}
	remain |= offset
	b.WriteByte(remain)
	return b.Bytes()
}

type Decoder map[string]interface{}

func (d Decoder) DecodeToStringSlice(code []byte) (result []string) {
	search := []byte{0}
	var index int = 0
	var offset byte = 1
	i := 0
	for ; i < len(code)-1; i++ {
		var ofs byte = 1
		for k := 0; k < 8; k++ {
			if code[i]&ofs == ofs {
				search[index] |= offset
			}
			if offset == 128 {
				index++
				if len(search) == index {
					search = append(search, 0)
				}
				offset = 1
			} else {
				offset = offset << 1
			}
			search[index] |= offset
			v, ok := d[string(search[:index+1])]
			search[index] ^= offset
			if ok {
				result = append(result, v.(string))
				for j := 0; j <= index; j++ {
					search[j] = 0
				}
				index = 0
				offset = 1
			}
			ofs = ofs << 1
		}
	}
	end := 0
	var ofs byte = 128
	for k := 0; k < 8; k++ {
		if code[i]&ofs != ofs {
			ofs = ofs >> 1
			continue
		}
		end = 7 - k
		break
	}
	ofs = 1
	for k := 0; k < end; k++ {
		if code[i]&ofs == ofs {
			search[index] |= offset
		}
		if offset == 128 {
			index++
			if len(search) == index {
				search = append(search, 0)
			}
			offset = 1
		} else {
			offset = offset << 1
		}
		search[index] |= offset
		v, ok := d[string(search[:index+1])]
		search[index] ^= offset
		if ok {
			result = append(result, v.(string))
			for j := 0; j <= index; j++ {
				search[j] = 0
			}
			index = 0
			offset = 1
		}
		ofs = ofs << 1
	}
	return
}
