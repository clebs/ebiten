// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphicsutil

import (
	"github.com/hajimehoshi/ebiten/internal/graphics"
	"github.com/hajimehoshi/ebiten/internal/opengl"
)

var (
	theVerticesBackend = &verticesBackend{}
)

type verticesBackend struct {
	backend []float32
	head    int
}

func (v *verticesBackend) sliceForOneQuad() []float32 {
	const num = 256
	size := 4 * graphics.VertexSizeInBytes() / opengl.Float.SizeInBytes()
	if v.backend == nil {
		v.backend = make([]float32, size*num)
	}
	s := v.backend[v.head : v.head+size]
	v.head += size
	if v.head+size > len(v.backend) {
		v.backend = nil
		v.head = 0
	}
	return s
}

func isPowerOf2(x int) bool {
	return (x & (x - 1)) == 0
}

func QuadVertices(width, height int, sx0, sy0, sx1, sy1 int, a, b, c, d, tx, ty float32) []float32 {
	if !isPowerOf2(width) {
		panic("not reached")
	}
	if !isPowerOf2(height) {
		panic("not reached")
	}

	if sx0 >= sx1 || sy0 >= sy1 {
		return nil
	}
	if sx1 <= 0 || sy1 <= 0 {
		return nil
	}

	wf := float32(width)
	hf := float32(height)
	u0, v0, u1, v1 := float32(sx0)/wf, float32(sy0)/hf, float32(sx1)/wf, float32(sy1)/hf
	return quadVerticesImpl(float32(sx1-sx0), float32(sy1-sy0), u0, v0, u1, v1, a, b, c, d, tx, ty)
}

func quadVerticesImpl(x, y, u0, v0, u1, v1, a, b, c, d, tx, ty float32) []float32 {
	// Specifying a range explicitly here is redundant but this helps optimization
	// to eliminate boundry checks.
	vs := theVerticesBackend.sliceForOneQuad()[0:24]

	ax, by, cx, dy := a*x, b*y, c*x, d*y

	// Vertex coordinates
	vs[0] = tx
	vs[1] = ty

	// Texture coordinates: first 2 values indicates the actual coodinate, and
	// the second indicates diagonally opposite coodinates.
	// The second is needed to calculate source rectangle size in shader programs.
	vs[2] = u0
	vs[3] = v0
	vs[4] = u1
	vs[5] = v1

	// and the same for the other three coordinates
	vs[6] = ax + tx
	vs[7] = cx + ty
	vs[8] = u1
	vs[9] = v0
	vs[10] = u0
	vs[11] = v1

	vs[12] = by + tx
	vs[13] = dy + ty
	vs[14] = u0
	vs[15] = v1
	vs[16] = u1
	vs[17] = v0

	vs[18] = ax + by + tx
	vs[19] = cx + dy + ty
	vs[20] = u1
	vs[21] = v1
	vs[22] = u0
	vs[23] = v0

	return vs
}
