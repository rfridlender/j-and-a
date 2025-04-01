package repositories

import (
	"log"
	"reflect"
	"sort"
	"time"

	"j-and-a/internal/models"
)

type lessFunc func(d1, d2 models.ModelData) bool

type multiSorter struct {
	modelDatas []models.ModelData
	lesses     []lessFunc
}

func (ms *multiSorter) Sort(modelDatas []models.ModelData) {
	ms.modelDatas = modelDatas
	sort.Sort(ms)
}

func OrderedBy(lesses ...lessFunc) *multiSorter {
	return &multiSorter{
		lesses: lesses,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.modelDatas)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.modelDatas[i], ms.modelDatas[j] = ms.modelDatas[j], ms.modelDatas[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.modelDatas[i], ms.modelDatas[j]
	var k int
	for k = 0; k < len(ms.lesses)-1; k++ {
		less := ms.lesses[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.lesses[k](p, q)
}

func isDeleted(d1, d2 models.ModelData) bool {
	v1 := reflect.ValueOf(d1)
	v2 := reflect.ValueOf(d2)
	if v1.Kind() != reflect.Pointer || v2.Kind() != reflect.Pointer {
		log.Fatalln("model data must be pointer to struct")
	}
	v1 = v1.Elem()
	v2 = v2.Elem()
	if v1.Kind() != reflect.Struct || v2.Kind() != reflect.Struct {
		log.Fatalln("model data must be pointer to struct")
	}
	v1 = v1.FieldByName("DeletedAt")
	v2 = v2.FieldByName("DeletedAt")
	return v1.Len() < v2.Len()
}

func updatedAt(d1, d2 models.ModelData) bool {
	v1 := reflect.ValueOf(d1)
	v2 := reflect.ValueOf(d2)
	if v1.Kind() != reflect.Pointer || v2.Kind() != reflect.Pointer {
		log.Fatalln("model data must be pointer to struct")
	}
	v1 = v1.Elem()
	v2 = v2.Elem()
	if v1.Kind() != reflect.Struct || v2.Kind() != reflect.Struct {
		log.Fatalln("model data must be pointer to struct")
	}
	updatedAt1 := v1.FieldByName("DeletedAt")
	if updatedAt1.Len() == 0 {
		updatedAt1 = v1.FieldByName("CreatedAt")
	}
	updatedAt2 := v2.FieldByName("DeletedAt")
	if updatedAt2.Len() == 0 {
		updatedAt2 = v2.FieldByName("CreatedAt")
	}
	t1, err := time.Parse(time.RFC3339, updatedAt1.String())
	if err != nil {
		log.Fatalln(err)
	}
	t2, err := time.Parse(time.RFC3339, updatedAt2.String())
	if err != nil {
		log.Fatalln(err)
	}
	return t1.After(t2)
}
