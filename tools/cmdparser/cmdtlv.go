package cmdparser

type Tlv struct {
    LeafType Type
    Id string
    Value string
}

type SerBuff struct {
    Data *Tlv
    Next *SerBuff
}
