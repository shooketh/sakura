package generator_test

import (
	"testing"

	"github.com/shooketh/sakura/common/errors"
	"github.com/shooketh/sakura/module/generator"
	"github.com/stretchr/testify/assert"
)

func TestInvalidMachineID(t *testing.T) {
	_, err := generator.New(31, 31, 0)
	assert.Nil(t, err, "unexpected error")

	_, err = generator.New(32, 31, 0)
	assert.Equal(t, err, errors.InvalidDatacenterID, "invalid error for overranged datacenterID")

	_, err = generator.New(31, 32, 0)
	assert.Equal(t, err, errors.InvalidWorkerID, "invalid error for overranged workerID")
}

func TestGenerateIDs(t *testing.T) {
	g, err := generator.New(0, 0, 0)
	assert.Nil(t, err, "failed to create new generator")

	var ids []int64

	for i := 0; i < 10000; i++ {
		id := g.Generate()

		for _, otherID := range ids {
			assert.True(t, otherID != id, "id duplicated!!")
		}

		l := len(ids)
		assert.True(t, l >= 0 && (ids == nil || id > ids[l-1]), "generated smaller id!!")

		ids = append(ids, id)
	}

	t.Logf("%d ids are tested", len(ids))
}

func BenchmarkGenerate(b *testing.B) {
	g, _ := generator.New(0, 0, 0)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = g.Generate()
	}
}
