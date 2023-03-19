package templates

import (
	"fmt"
	"net/url"
	"os"

	"hammerpost/global"

	"github.com/flosch/pongo2/v6"
)

func processDBUriGetCreds(dbUri string) (host string, port string, dbuser string, dbpassword string, err error) {
	// Parse the uri
	uri, err := url.Parse(dbUri)
	if err != nil {
		return "", "", "", "", err
	}

	// Get the user and password
	dbuser = uri.User.Username()
	dbpassword, _ = uri.User.Password()
	host = uri.Hostname()
	port = uri.Port()

	return host, port, dbuser, dbpassword, nil
}

func createPGSchemaTemplate(uri string, users int, warehouses int) error {
	// Replace the placeholders with the actual values

	tpl, err := pongo2.FromFile("hammer-templates/pg/schema.tpl")
	if err != nil {
		return err
	}

	// Create the new schema template file
	w, err := os.Create(fmt.Sprintf("schema_%s.tcl", global.BenchmarkID))
	if err != nil {
		return err
	}

	dbhost, dbport, dbuser, dbpassword, err := processDBUriGetCreds(uri)
	if err != nil {
		return err
	}

	// Create the schema file
	err = tpl.ExecuteWriter(pongo2.Context{
		"db_host":     dbhost,
		"db_port":     dbport,
		"db_user":     dbuser,
		"db_password": dbpassword,
		"users":       users, "warehouses": warehouses,
	}, w)
	return err
}

func createPGRunTemplate(uri string, users int, trx int, duration int, rampup int, allwarehouses bool, warehouses int) error {
	// Replace the placeholders with the actual values

	tpl, err := pongo2.FromFile("hammer-templates/pg/run.tpl")
	if err != nil {
		return err
	}

	// Create the new run template file
	w, err := os.Create(fmt.Sprintf("run_%s.tcl", global.BenchmarkID))
	if err != nil {
		return err
	}

	dbhost, dbport, dbuser, dbpassword, err := processDBUriGetCreds(uri)
	if err != nil {
		return err
	}

	// Create the run file
	err = tpl.ExecuteWriter(pongo2.Context{"db_host": dbhost,
		"db_port":     dbport,
		"db_user":     dbuser,
		"db_password": dbpassword, "users": users, "total_transactions": trx, "test_duration": duration, "rampup_duration": rampup,
		"all_warehouses": allwarehouses, "warehouses": warehouses,
	}, w)
	return err
}

func createMySQLSchemaTemplate(uri string, users int, warehouses int) error {
	// Replace the placeholders with the actual values

	tpl, err := pongo2.FromFile("hammer-templates/mysql/schema.tpl")
	if err != nil {
		return err
	}

	// Create the new schema template file
	w, err := os.Create(fmt.Sprintf("schema_%s.tcl", global.BenchmarkID))
	if err != nil {
		return err
	}

	dbhost, dbport, dbuser, dbpassword, err := processDBUriGetCreds(uri)
	if err != nil {
		return err
	}

	// Create the schema file
	err = tpl.ExecuteWriter(pongo2.Context{
		"db_host":     dbhost,
		"db_port":     dbport,
		"db_user":     dbuser,
		"db_password": dbpassword,
		"users":       users, "warehouses": warehouses,
	}, w)
	return err
}

func createMySQLRunTemplate(uri string, users int, trx int, duration int, rampup int, allwarehouses bool, warehouses int) error {
	// Replace the placeholders with the actual values

	tpl, err := pongo2.FromFile("hammer-templates/mysql/run.tpl")
	if err != nil {
		return err
	}

	// Create the new run template file
	w, err := os.Create(fmt.Sprintf("run_%s.tcl", global.BenchmarkID))
	if err != nil {
		return err
	}

	dbhost, dbport, dbuser, dbpassword, err := processDBUriGetCreds(uri)
	if err != nil {
		return err
	}

	// Create the run file
	err = tpl.ExecuteWriter(pongo2.Context{"db_host": dbhost,
		"db_port":     dbport,
		"db_user":     dbuser,
		"db_password": dbpassword,
		"users":       users, "total_transactions": trx, "test_duration": duration, "rampup_duration": rampup,
		"warehouses": warehouses}, w)
	return err
}

func CreateTemplateFiles(dbType string, uri string, users int, warehouses int, trx int, duration int, rampup int, allwarehouses bool) error {

	if dbType == "postgres" {
		err := createPGSchemaTemplate(uri, users, warehouses)
		if err != nil {
			return err
		}
		return createPGRunTemplate(uri, users, trx, duration, rampup, allwarehouses, warehouses)
	} else if dbType == "mysql" {
		err := createMySQLSchemaTemplate(uri, users, warehouses)
		if err != nil {
			return err
		}

		return createMySQLRunTemplate(uri, users, trx, duration, rampup, allwarehouses, warehouses)
	}
	return fmt.Errorf("unsupported dbtype: %s", dbType)

}
