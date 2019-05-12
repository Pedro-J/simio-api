package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func cleanFiles() {
	currentDir, _ := os.Getwd()
	os.RemoveAll(currentDir + "/database/")
}

func TestSave(t *testing.T) {
	assert := assert.New(t)

	dataMoreHumans := make(map[string]SimioEntity)
	dataMoreHumans["111"] = SimioEntity{ID: "111", DNA: "ACCG|DGCT", IsSimian: true}
	dataMoreHumans["222"] = SimioEntity{ID: "222", DNA: "AGCG|GGCT", IsSimian: false}
	dataMoreHumans["333"] = SimioEntity{ID: "333", DNA: "AACG|DTTT", IsSimian: false}

	simioDAO := SimioDAO{
		Data: dataMoreHumans,
	}

	newEntity := SimioEntity{ID: "4454", DNA: "AACG|DTTT", IsSimian: false}

	err := simioDAO.Save(newEntity)
	assert.Nil(err)
	simioDAO.Save(newEntity)

	defer cleanFiles()
}

func TestBuildSimioDAO(t *testing.T) {
	assert := assert.New(t)

	simioDAO := NewSimioDAO("dasdadasd")

	assert.Empty(simioDAO.GetData())

	entities := []SimioEntity{
		SimioEntity{ID: "111", DNA: "ACCG|DGCT", IsSimian: true},
		SimioEntity{ID: "222", DNA: "AGCG|GGCT", IsSimian: false},
		SimioEntity{ID: "333", DNA: "AACG|DTTT", IsSimian: false},
	}

	for _, entity := range entities {
		saveEntityOnFile(entity.ID, entity)
	}

	simioDAO = BuildSimioDAO()

	assert.NotNil(simioDAO)
	assert.Equal(len(entities), len(simioDAO.GetData()))

	defer cleanFiles()
}
