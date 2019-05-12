package service

import (
	"fmt"
	"simio-api/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	dnaSimianHorizontal         = []string{"CCCG", "AAAT", "GGGA", "TTTT"}
	dnaSimianVertical           = []string{"CCCC", "AAAC", "GGGC", "TTTC"}
	dnaSimianDiagonal           = []string{"CCTG", "ACAC", "GGCA", "TTTC"}
	dnaSimianDiagonal2          = []string{"TCTA", "GTAC", "GACA", "ATTC"}
	dnaHuman                    = []string{"CGAT", "GTCA", "TACG", "TCGA"}
	dnaSimianDiagonalAlta6x6    = []string{"CTGGGA", "CTATGC", "TATTGA", "AGAGAG", "CCCATA", "TCACTG"}
	dnaSimianDiagonalBaixa6x6   = []string{"CTGAGA", "CTATGC", "TATTGA", "AGAGAG", "CCCATA", "TCACTG"}
	dnaSimianDiagonalInversa6x6 = []string{"CTGGGA", "CATGGC", "TATTGA", "AGATTG", "CCCGTA", "TCACTG"}
	dna7x7                      = []string{"CTGAGAG", "CTATGCC", "TATTGTA", "AGAGGGA", "CCCCTAT", "TCACTGT", "TTAAGTA"}
	dna8x8                      = []string{"AGTCCCTA", "GCAGGAAT", "TTCCAAGG", "TCAATTGC", "GGTTCCAG", "CCTAGGCC", "TTGCGCAA", "AAACCGTG"}
	dna4x3                      = []string{"CTGZ", "CTAT", "TATT"}
	dna2x2                      = []string{"GT", "AT"}
	dna3x3                      = []string{"CAG", "CGA", "CCC"}
	dna1x1                      = []string{"C"}
	dna1x4                      = []string{"CCCC"}
	dnaDifCols                  = []string{"CTGZ", "CT", "TATTTGGAATTT"}
	dnaInvalidLastChar          = []string{"AGTCCCTA", "GCAGGAAT", "TTCCAAGG", "TCAATTGC", "GGTTCCAG", "CCTAGGCC", "TTGCGCAA", "AAACCGTZ"}
	dnaInvalidFirstChar         = []string{"ZGTCCCTA", "GCAGGAAT", "TTCCAAGG", "TCAATTGC", "GGTTCCAG", "CCTAGGCC", "TTGCGCAA", "AAACCGTA"}
	dnaEmpty                    = []string{}
)

//Mocking simio dao

type SimioDaoMock struct {
	mock.Mock
	database.DAO
}

func (sm *SimioDaoMock) Save(entity database.SimioEntity) error {
	args := sm.Called(entity)
	return args.Error(0)
}

func (sm *SimioDaoMock) GetData() map[string]database.SimioEntity {
	args := sm.Called()
	return args.Get(0).(map[string]database.SimioEntity)
}

//Start Tests

func TestProcessDNA(t *testing.T) {
	assert := assert.New(t)

	type Case struct {
		instance       []string
		expectedErr    error
		expectedResult bool
	}

	cases := []Case{
		//simians
		Case{instance: dnaSimianHorizontal, expectedResult: true, expectedErr: nil},
		Case{instance: dnaSimianVertical, expectedResult: true, expectedErr: nil},
		Case{instance: dnaSimianDiagonal, expectedResult: true, expectedErr: nil},

		//Valid but not simian
		Case{instance: dna3x3, expectedResult: false, expectedErr: nil},
		Case{instance: dnaEmpty, expectedResult: false, expectedErr: nil},
		Case{instance: dna1x1, expectedResult: false, expectedErr: nil},

		//Invalids matrix
		Case{instance: dnaDifCols, expectedResult: false, expectedErr: nil},
		Case{instance: dnaInvalidFirstChar, expectedResult: false, expectedErr: fmt.Errorf("")},
		Case{instance: dnaInvalidLastChar, expectedResult: false, expectedErr: fmt.Errorf("")},
	}

	for _, currentCase := range cases {
		simioDaoMock := new(SimioDaoMock)
		simioDaoMock.On("Save", mock.Anything).Return(nil)
		simioService := NewSimioService(4, simioDaoMock)

		res, err := simioService.ProcessDNA(currentCase.instance)

		if err != nil {
			assert.False(res)
			simioDaoMock.AssertNotCalled(t, "Save")
		} else {
			simioDaoMock.AssertNumberOfCalls(t, "Save", 1)
		}
	}

}

func TestGetSimiansProportion(t *testing.T) {
	assert := assert.New(t)

	type Case struct {
		data           map[string]database.SimioEntity
		expectedResult Stats
	}

	dataMoreHumans := make(map[string]database.SimioEntity)
	dataMoreHumans["1"] = database.SimioEntity{IsSimian: true}
	dataMoreHumans["2"] = database.SimioEntity{IsSimian: false}
	dataMoreHumans["3"] = database.SimioEntity{IsSimian: false}

	dataMoreSimians := make(map[string]database.SimioEntity)
	dataMoreSimians["1"] = database.SimioEntity{IsSimian: true}
	dataMoreSimians["2"] = database.SimioEntity{IsSimian: true}
	dataMoreSimians["3"] = database.SimioEntity{IsSimian: true}
	dataMoreSimians["4"] = database.SimioEntity{IsSimian: false}

	dataNoHumans := make(map[string]database.SimioEntity)
	dataNoHumans["1"] = database.SimioEntity{IsSimian: true}

	dataNoSimians := make(map[string]database.SimioEntity)
	dataNoSimians["1"] = database.SimioEntity{IsSimian: false}

	dataEqualHumanAndSimian := make(map[string]database.SimioEntity)
	dataEqualHumanAndSimian["1"] = database.SimioEntity{IsSimian: true}
	dataEqualHumanAndSimian["2"] = database.SimioEntity{IsSimian: true}
	dataEqualHumanAndSimian["3"] = database.SimioEntity{IsSimian: false}
	dataEqualHumanAndSimian["4"] = database.SimioEntity{IsSimian: false}

	cases := []Case{
		Case{data: dataMoreHumans, expectedResult: Stats{CountMutantDNA: 1, CountHumanDNA: 2, Ratio: float64(1) / float64(2)}},
		Case{data: dataMoreSimians, expectedResult: Stats{CountMutantDNA: 3, CountHumanDNA: 1, Ratio: float64(3) / float64(1)}},
		Case{data: dataNoHumans, expectedResult: Stats{CountMutantDNA: 1, CountHumanDNA: 0, Ratio: float64(0)}},
		Case{data: dataNoSimians, expectedResult: Stats{CountMutantDNA: 0, CountHumanDNA: 1, Ratio: float64(0)}},
		Case{data: dataEqualHumanAndSimian, expectedResult: Stats{CountMutantDNA: 2, CountHumanDNA: 2, Ratio: float64(2) / float64(2)}},
	}

	for _, currentCase := range cases {
		simioDaoMock := new(SimioDaoMock)
		simioDaoMock.On("GetData", mock.Anything).Return(currentCase.data)
		simioService := NewSimioService(4, simioDaoMock)

		statsResult := simioService.GetSimiansProportion()

		assert.Equal(currentCase.expectedResult.Ratio, statsResult.Ratio)
		assert.Equal(currentCase.expectedResult.CountHumanDNA, statsResult.CountHumanDNA)
		assert.Equal(currentCase.expectedResult.CountMutantDNA, statsResult.CountMutantDNA)
	}
}

func TestMapToSimioEntity(t *testing.T) {
	assert := assert.New(t)

	simioDaoMock := new(SimioDaoMock)
	simioService := NewSimioService(4, simioDaoMock).(*SimioServiceImpl)

	type Case struct {
		dna            []string
		isSimian       bool
		expectedResult database.SimioEntity
	}

	cases := []Case{
		Case{
			dna: dnaHuman, isSimian: false, expectedResult: database.SimioEntity{
				DNA:      simioService.getStringDNA(dnaHuman),
				ID:       simioService.generateId(simioService.getStringDNA(dnaHuman)),
				IsSimian: false,
			},
		},

		Case{
			dna: dnaSimianDiagonal, isSimian: true, expectedResult: database.SimioEntity{
				DNA:      simioService.getStringDNA(dnaSimianDiagonal),
				ID:       simioService.generateId(simioService.getStringDNA(dnaSimianDiagonal)),
				IsSimian: true,
			},
		},
	}

	for _, currentCase := range cases {

		entityResult := simioService.mapToSimioEntity(currentCase.dna, currentCase.isSimian)

		assert.Equal(currentCase.expectedResult.DNA, entityResult.DNA)
		assert.Equal(currentCase.expectedResult.ID, entityResult.ID)
		assert.Equal(currentCase.expectedResult.IsSimian, entityResult.IsSimian)
	}

}

func TestIsSimian(t *testing.T) {
	assert := assert.New(t)

	simioDaoMock := new(SimioDaoMock)
	simioService := NewSimioService(4, simioDaoMock).(*SimioServiceImpl)

	type Case struct {
		dna            []string
		expectedResult bool
	}

	cases := []Case{
		//Simians
		Case{dna: dnaSimianHorizontal, expectedResult: true},
		Case{dna: dnaSimianVertical, expectedResult: true},
		Case{dna: dnaSimianDiagonal, expectedResult: true},
		Case{dna: dnaSimianDiagonal2, expectedResult: true},
		Case{dna: dnaSimianDiagonalAlta6x6, expectedResult: true},
		Case{dna: dnaSimianDiagonalBaixa6x6, expectedResult: true},
		Case{dna: dnaSimianDiagonalInversa6x6, expectedResult: true},

		//Not simians
		Case{dna: dna3x3, expectedResult: false},
		Case{dna: dnaEmpty, expectedResult: false},
		Case{dna: dna1x1, expectedResult: false},
	}

	for _, currentCase := range cases {
		resultIsSimian := simioService.isSimian(currentCase.dna)

		assert.Equal(currentCase.expectedResult, resultIsSimian)
	}

}
