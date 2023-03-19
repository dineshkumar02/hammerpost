package mysql

import (
	"os"

	"hammerpost/model"
)

// Update my.cnf file
func UpdateMysqlParameter(mycnfPath string, param []model.Param) error {
	// append my.cnf with the given parameter

	f, e := os.OpenFile(mycnfPath, os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		return e
	}

	defer f.Close()

	for _, p := range param {
		if _, e = f.WriteString(p.Name + "=" + p.Value); e != nil {
			return e
		}
	}
	return nil
}
