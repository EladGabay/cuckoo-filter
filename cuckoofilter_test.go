/*
 * Copyright (C) linvon
 * Date  2021/2/18 10:29
 */

package cuckoo

import (
	"crypto/rand"
	"fmt"
	"io"
	"reflect"
	"testing"
)

const size = 100000

var testBucketSize = []uint{2, 4, 8}
var testFingerprintSize = []uint{2, 4, 5, 6, 7, 8, 9, 10, 12, 13, 16, 17, 32}
var testTableType = []uint{TableTypeSingle, TableTypePacked}

func TestFilter(t *testing.T) {
	var insertNum uint = 50000
	var hash [32]byte

	for _, b := range testBucketSize {
		for _, f := range testFingerprintSize {
			for _, table := range testTableType {
				if f == 2 && table == TableTypePacked {
					continue
				}
				if table == TableTypePacked {
					b = 4
				}
				cf := NewFilter(b, f, 8190, table)
				//fmt.Println(cf.Info())
				a := make([][]byte, 0)
				for i := uint(0); i < insertNum; i++ {
					_, _ = io.ReadFull(rand.Reader, hash[:])
					if cf.AddUnique(hash[:]) {
						tmp := make([]byte, 32)
						copy(tmp, hash[:])
						a = append(a, tmp)
					}
				}

				count := cf.Size()
				if count != uint(len(a)) {
					t.Errorf("Expected count = %d, instead count = %d, b %v f %v", uint(len(a)), count, b, f)
				}

				for _, v := range a {
					if !cf.Contain(v) {
						t.Errorf("Expected contain, instead not contain")
					}
				}

				for _, v := range a {
					cf.Delete(v)
				}

				count = cf.Size()
				if count != 0 {
					t.Errorf("Expected count = 0, instead count == %d", count)
				}

				bytes := cf.Encode()
				ncf, err := Decode(bytes)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !reflect.DeepEqual(cf, ncf) {
					t.Errorf("Expected %v, got %v", cf, ncf)
				}

				cf.Info()
				cf.BitsPerItem()
				cf.SizeInBytes()
				cf.LoadFactor()
				fmt.Printf("Filter bucketSize %v fingerprintSize %v tableType %v falsePositive Rate %v \n", b, f, table, cf.FalsePositiveRate())
			}
		}
	}

}

func BenchmarkFilterSingle_Reset(b *testing.B) {
	filter := NewFilter(4, 8, size, TableTypeSingle)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Reset()
	}
}

func BenchmarkFilterSingle_Insert(b *testing.B) {
	filter := NewFilter(4, 8, size, TableTypeSingle)

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		_, _ = io.ReadFull(rand.Reader, hash[:])
		filter.Add(hash[:])
	}
}

func BenchmarkFilterSingle_Lookup(b *testing.B) {
	filter := NewFilter(4, 8, size, TableTypeSingle)

	var hash [32]byte
	for i := 0; i < size; i++ {
		_, _ = io.ReadFull(rand.Reader, hash[:])
		filter.Add(hash[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = io.ReadFull(rand.Reader, hash[:])
		filter.Contain(hash[:])
	}
}

func BenchmarkFilterPacked_Reset(b *testing.B) {
	filter := NewFilter(4, 9, size, TableTypePacked)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Reset()
	}
}

func BenchmarkFilterPacked_Insert(b *testing.B) {
	filter := NewFilter(4, 9, size, TableTypePacked)

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		_, _ = io.ReadFull(rand.Reader, hash[:])
		filter.Add(hash[:])
	}
}

func BenchmarkFilterPacked_Lookup(b *testing.B) {
	filter := NewFilter(4, 9, size, TableTypePacked)

	var hash [32]byte
	for i := 0; i < size; i++ {
		_, _ = io.ReadFull(rand.Reader, hash[:])
		filter.Add(hash[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = io.ReadFull(rand.Reader, hash[:])
		filter.Contain(hash[:])
	}
}
