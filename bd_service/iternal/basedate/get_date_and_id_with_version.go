package basedate

import myerrors "dkl.ru/pact/bd_service/iternal/my_errors"

type FileWithVersion struct {
	FileData  File   `json:"file"`
	CreatedAt string `json:"created_at"`
}

func (d *Database) GetDateIdWithWersion(version int) ([]FileWithVersion, error) {
	return nil, myerrors.NotRealizeable("GetDateIdWithWersion")
}
