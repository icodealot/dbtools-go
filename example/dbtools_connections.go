package example

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/databasetools"
)

// DBToolsConfig is the configuration for the example and is primarily used for
// setting the endpoint and payload.
//
//  - ConnectionId is the OCID of the DBTools connection
//  - ContentType should be either application/sql or application/json
//  - Payload is the SQL to execute (either SQL or JSON format)
//
// In addition to the above, the Go SDK for  provides a DefaultConfigProvider
// and should be configured with a [DEFAULT] section in the ~/.oci/config file.
type DBToolsConfig struct {
	ConnectionId string `json:"connectionId"`
	ContentType  string `json:"contentType"`
	Payload      string `json:"payload"`
}

// example.ExecuteDBToolsConnection() calls the Database Tools ORDS endpoint
// to execute some SQL queries against the database at the other end of the
// connection. This should work even if the database is sitting behind a private
// subnet as long as the connection was setup using a private endpoint.
//
// Queries are executed by signing requests using the details from the default
// profile in the ~/.oci/config file. (i.e. [DEFAULT])
func ExecuteDBToolsConnection(cfg DBToolsConfig) ([]byte, error) {

	endpoint := GetDatabaseToolsEndpoint(cfg)
	payload := bytes.NewBufferString(cfg.Payload)

	request, err := http.NewRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// UTC() time is important, otherwise your requests will likely fail
	request.Header.Set("date", time.Now().UTC().Format(http.TimeFormat))
	request.Header.Set("content-type", cfg.ContentType)

	signer := common.DefaultRequestSigner(common.DefaultConfigProvider())
	err = signer.Sign(request)
	if err != nil {
		return nil, fmt.Errorf("error signing request: %v", err)
	}

	return doHttpRequest(request)
}

// example.ValidateDBToolsCOnnection() calls the Database Tools API to validate
// a connection based on the OCID provided.
func ValidateDBToolsConnection(cfg DBToolsConfig) {
	log.Println("Validating a Database Tools Connection: " + cfg.ConnectionId)

	ocid := cfg.ConnectionId

	client, err := databasetools.NewDatabaseToolsClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		log.Fatalf("Error creating dbtools client: %v", err)
	}

	ctx := context.Background()

	connection, err := client.GetDatabaseToolsConnection(ctx,
		databasetools.GetDatabaseToolsConnectionRequest{
			DatabaseToolsConnectionId: common.String(ocid),
		})
	if err != nil {
		log.Fatalf("Error validating database tools connection: %v", err)
	}

	var connectionDetailsType databasetools.ValidateDatabaseToolsConnectionDetails

	// You need to specify the type of connection for validation so you can
	// do that with type assertion, for example:
	switch connection.DatabaseToolsConnection.(type) {
	case databasetools.DatabaseToolsConnectionMySql:
		connectionDetailsType = databasetools.ValidateDatabaseToolsConnectionMySqlDetails{}
	case databasetools.DatabaseToolsConnectionOracleDatabase:
		connectionDetailsType = databasetools.ValidateDatabaseToolsConnectionOracleDatabaseDetails{}
	default:
		log.Fatalf("unexpected connection type: %T", connection)
	}

	validate, err := client.ValidateDatabaseToolsConnection(ctx,
		databasetools.ValidateDatabaseToolsConnectionRequest{
			DatabaseToolsConnectionId:              common.String(ocid),
			ValidateDatabaseToolsConnectionDetails: connectionDetailsType,
		})
	if err != nil {
		log.Fatalf("Error validating database tools connection: %v", err)
	}

	log.Printf("Validation response: %+v\n", *validate.GetCode())
}

// Use the information from the ~/.oci/config file (for the region) combined with
// the OCID of the connection to get the ORDS endpoint. For example:
// https://sql.dbtools.us-phoenix-1.oci.oraclecloud.com/20201005/ords/ocid1.databasetoolsconnection...change-me/_/sql
func GetDatabaseToolsEndpoint(cfg DBToolsConfig) string {
	client, err := databasetools.NewDatabaseToolsClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		log.Fatalf("Error creating dbtools client: %v", err)
	}

	url := strings.Split(client.Endpoint(), "//")
	if len(url) <= 1 {
		log.Fatalf("error getting dbtools endpoint")
	}
	prot := url[0]
	host := url[1]

	return fmt.Sprintf("%s//sql.%s/%s/ords/%s/_/sql", prot, host, client.BasePath, cfg.ConnectionId)
}

// A helper function to do the actual HTTP request. Just to keep the example
// code clean and in case we want to add more logic in the future.
func doHttpRequest(request *http.Request) ([]byte, error) {
	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}
