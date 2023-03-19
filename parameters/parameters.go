package parameters

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"hammerpost/model"
	"hammerpost/util"

	"github.com/schwarmco/go-cartesian-product"
)

func ReadParameterFile(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var parameters map[string]interface{}
	err = json.Unmarshal(bytes, &parameters)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

func PrepareParameters(parameters map[string]interface{}) []model.Parameters {
	var parametersList []model.Parameters

	// Prepare parameters
	for key, value := range parameters {
		var parameter model.Parameters
		parameter.Name = key
		parameter.Values = util.InterFaceToStrArray(value)
		parametersList = append(parametersList, parameter)
	}

	return parametersList
}

func GeneratePossibleParameters(params []model.Parameters) (result []interface{}) {

	var allParams [][]interface{}

	for _, param := range params {
		var paramValues []interface{}
		var p model.Param
		p.Name = param.Name
		for _, value := range param.Values {
			p.Value = value
			paramValues = append(paramValues, p)
		}
		allParams = append(allParams, paramValues)
	}

	// Generate cartesian product of all the parameters
	c := cartesian.Iter(allParams...)

	for product := range c {
		result = append(result, product)
	}

	// fmt.Println(result)
	return result
}

func GetParametersAsString(params []model.Param) string {
	var result string
	for _, param := range params {
		result += param.Name + ":" + param.Value + "\n"
	}

	return result
}
