package output

import (
	"bytes"
	"os"

	"miller/clitypes"
	"miller/types"
)

type RecordWriterDKVP struct {
	ofs string
	ops string
	ors string
}

func NewRecordWriterDKVP(writerOptions *clitypes.TWriterOptions) *RecordWriterDKVP {
	return &RecordWriterDKVP{
		ofs: writerOptions.OFS,
		ops: writerOptions.OPS,
		ors: writerOptions.ORS,
	}
}

func (this *RecordWriterDKVP) Write(
	outrec *types.Mlrmap,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(*pe.Key)
		buffer.WriteString(this.ops)
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(this.ofs)
		}
	}
	buffer.WriteString(this.ors)
	os.Stdout.WriteString(buffer.String())
}
