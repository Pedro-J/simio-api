package service

import (
	"crypto/sha1"
	"fmt"
	"simio-api/database"
)

type Stats struct {
	CountMutantDNA int     `json:"count_mutant_dna"`
	CountHumanDNA  int     `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

type SimioService interface {
	ProcessDNA(dna []string) (bool, error)
	GetSimiansProportion() Stats
}

type SimioServiceImpl struct {
	sequenceSize int
	simioDAO     database.DAO
}

func (ss *SimioServiceImpl) ProcessDNA(DNA []string) (bool, error) {
	err := ss.validateDNA(DNA)

	if err != nil {
		return false, err
	}

	isSimian := ss.isSimian(DNA)
	ss.simioDAO.Save(ss.mapToSimioEntity(DNA, isSimian))

	return isSimian, nil
}

func (ss *SimioServiceImpl) GetSimiansProportion() Stats {
	data := ss.simioDAO.GetData()

	simians, humans := 0, 0

	for _, entity := range data {
		if entity.IsSimian {
			simians++
		} else {
			humans++
		}
	}

	var ratio float64
	if humans != 0 {
		ratio = float64(simians) / float64(humans)
	} else {
		ratio = float64(0)
	}

	return Stats{
		Ratio:          ratio,
		CountHumanDNA:  humans,
		CountMutantDNA: simians,
	}
}

func (ss *SimioServiceImpl) mapToSimioEntity(dna []string, isSimian bool) database.SimioEntity {
	stringDNA := ss.getStringDNA(dna)
	return database.SimioEntity{
		DNA:      stringDNA,
		ID:       ss.generateId(stringDNA),
		IsSimian: isSimian,
	}
}

func (ss *SimioServiceImpl) generateId(dna string) string {
	hashFunc := sha1.New()
	hashFunc.Write([]byte(dna))
	hashBytes := hashFunc.Sum(nil)
	return fmt.Sprintf("%x", hashBytes)
}

func (ss *SimioServiceImpl) getStringDNA(arrayDNA []string) string {

	var stringDNA string
	for i := 0; i < len(arrayDNA); i++ {
		stringDNA = stringDNA + "|" + arrayDNA[i]
	}

	return stringDNA[1:]
}

func (ss *SimioServiceImpl) isSimian(dna []string) bool {

	if ss.checkHorizontals(dna) || ss.checkVerticals(dna) || ss.checkDiagonals(dna) {
		return true
	}

	return false
}

func (ss *SimioServiceImpl) checkVerticals(dna []string) bool {
	n := len(dna)

	var lastBase byte
	var sequenceCount int
	for col := 0; col < n; col++ {
		for row := 0; row < n; row++ {
			if ss.verifySequence(&lastBase, dna[row][col], &sequenceCount) {
				return true
			}
		}
		lastBase = 0
		sequenceCount = 0
	}
	return false
}

func (ss *SimioServiceImpl) checkHorizontals(dna []string) bool {
	n := len(dna)

	var lastBase byte
	var sequenceCount int
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			if ss.verifySequence(&lastBase, dna[row][col], &sequenceCount) {
				return true
			}
		}
		lastBase = 0
		sequenceCount = 0
	}
	return false
}

func (ss *SimioServiceImpl) checkDiagonals(dna []string) bool {
	numDiagonals := (2 * len(dna)) - 1
	numTargetDiaginals := numDiagonals - (ss.sequenceSize-1)*2
	n := len(dna)
	lastIndex := n - 1

	var lastBase byte
	var sequenceCount int

	for left, right := 0, 0; left < numTargetDiaginals; left++ {
		startRow := (ss.sequenceSize - 1) + left

		if startRow < n {
			for row, col := startRow, 0; row >= 0; row, col = row-1, col+1 {
				if ss.verifySequence(&lastBase, dna[row][col], &sequenceCount) {
					return true
				}
			}
		} else {
			startCol := n - (ss.sequenceSize + right)
			for row, col := lastIndex, startCol; col < n; row, col = row-1, col+1 {
				if ss.verifySequence(&lastBase, dna[row][col], &sequenceCount) {
					return true
				}
			}
			right++
		}
		lastBase = 0
		sequenceCount = 0
	}

	for left, right := 0, 0; left < numTargetDiaginals; left++ {
		startRow := (n - ss.sequenceSize) - left

		if startRow >= 0 {
			for row, col := startRow, 0; row < n; row, col = row+1, col+1 {
				if ss.verifySequence(&lastBase, dna[row][col], &sequenceCount) {
					return true
				}
			}
		} else {
			startCol := (n - ss.sequenceSize) - right

			for row, col := 0, startCol; col < n; row, col = row+1, col+1 {
				if ss.verifySequence(&lastBase, dna[row][col], &sequenceCount) {
					return true
				}
			}
			right++
		}

		lastBase = 0
		sequenceCount = 0
	}
	return false
}

func (ss *SimioServiceImpl) verifySequence(lastBase *byte, currentBase byte, sequenceCount *int) bool {

	if isCharacterNotValid(currentBase) {
		lastBase = nil
		sequenceCount = nil
		return false
	}

	if *lastBase != 0 && currentBase == *lastBase {
		*sequenceCount++
	} else {
		*sequenceCount = 1
		*lastBase = currentBase
	}

	if *sequenceCount == ss.sequenceSize {
		return true
	}

	return false
}

func isCharacterNotValid(currentBase byte) bool {
	const baseA, baseT, baseC, baseG = byte('A'), byte('T'), byte('C'), byte('G')

	if currentBase != baseA && currentBase != baseT && currentBase != baseC && currentBase != baseG {
		return true
	}

	return false
}

func (ss *SimioServiceImpl) validateDNA(DNA []string) error {
	const baseA, baseT, baseC, baseG = byte('A'), byte('T'), byte('C'), byte('G')
	size := len(DNA)

	if len(DNA) == 0 {
		return fmt.Errorf("Invalid DNA size. the matrix is empty")
	}

	for row := 0; row < size; row++ {
		if len(DNA[row]) != size {
			return fmt.Errorf("Invalid DNA size. It has to be NxN")
		}

		for col := 0; col < size; col++ {
			if isCharacterNotValid(DNA[row][col]) {
				return fmt.Errorf("Matrix has invalid character ( %c )", DNA[row][col])
			}
		}
	}
	return nil
}

func BuildSimioService() SimioService {
	return NewSimioService(4, database.BuildSimioDAO())
}

func NewSimioService(sequenceSize int, dao database.DAO) SimioService {
	return &SimioServiceImpl{
		sequenceSize: sequenceSize,
		simioDAO:     dao,
	}
}
