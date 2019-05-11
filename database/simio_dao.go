package database

import "log"

type SimioEntity struct {
	ID       string
	DNA      string
	IsSimian bool
}

type DAO interface {
	Save(entity SimioEntity) error
	GetData() map[string]SimioEntity
}

type SimioDAO struct {
	Data map[string]SimioEntity
}

func (sDB *SimioDAO) Save(entity SimioEntity) error {
	_, hasEntity := sDB.Data[entity.ID]

	if !hasEntity {
		if !checkFileExist(entity.ID) {

			err := saveEntityOnFile(entity.ID, entity)

			if err != nil {
				return err
			}

			sDB.Data[entity.ID] = entity
		}
	} else {
		log.Printf("The DNA %s has been already saved", entity.ID)
	}
	return nil
}

func (sDB *SimioDAO) GetData() map[string]SimioEntity {
	return sDB.Data
}

func BuildSimioDAO() DAO {
	data, err := LoadAll()

	if err != nil {
		return &SimioDAO{
			Data: make(map[string]SimioEntity),
		}
	}

	return &SimioDAO{
		Data: data,
	}
}
