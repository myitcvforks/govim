package main

import (
	"encoding/json"
	"fmt"

	"github.com/govim/govim/cmd/govim/internal/types"
)

const (
	exprAutocmdCurrBufInfo = `{"Num": eval(expand('<abuf>')), "Name": fnamemodify(bufname(eval(expand('<abuf>'))),':p'), "Contents": join(getbufline(eval(expand('<abuf>')), 0, "$"), "\n")."\n", "Loaded": bufloaded(eval(expand('<abuf>')))}`
)

// currentBufferInfo is a helper function to unmarshal autocmd current
// buffer details from expr
func (v *vimstate) currentBufferInfo(expr json.RawMessage) *types.Buffer {
	var buf struct {
		Num      int
		Name     string
		Contents string
		Loaded   int
	}
	v.Parse(expr, &buf)
	return types.NewBuffer(buf.Num, buf.Name, []byte(buf.Contents), buf.Loaded == 1)
}

func (v *vimstate) cursorPos() (b *types.Buffer, p types.Point, err error) {
	var pos struct {
		BufNum int `json:"bufnum"`
		Line   int `json:"line"`
		Col    int `json:"col"`
	}
	expr := v.ChannelExpr(`{"bufnum": bufnr(""), "line": line("."), "col": col(".")}`)
	v.Parse(expr, &pos)
	b, ok := v.buffers[pos.BufNum]
	if !ok {
		err = fmt.Errorf("failed to resolve buffer %v", pos.BufNum)
		return
	}
	p, err = types.PointFromVim(b, pos.Line, pos.Col)
	return
}
