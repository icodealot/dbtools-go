# Database Tools Examples Using Go
Examples for the OCI Go SDK made possible through Database Tools connections and private endpoints

## Executing SQL via DBTOOLS Connection

Executing SQL using a Database Tools Connection is fairly straightforward. To make this work you will need a few pieces of information:

1. The connection ID (OCID)
2. The ORDS endpoint in your region
3. The SQL to execute

Note: The examples here use the ```DefaultConfigProvider``` to read the necessary bits from a ```[DEFAULT]``` profile for signing requests but you can use also use Session Token authenticaion (see the [OCI Go SDK README](https://github.com/oracle/oci-go-sdk/blob/master/README.md#configuring) for more information).

### The Connection ID

The Database Tools connection ID will look similar to:

```go
connection := "ocid1.databasetoolsconnection...change-me"
```

You can find this in the OCI console of your tenancy under Developer Services -> Connections -> Inspect the connection -> Connection information -> OCID.

### The Endpoint

A Database Tools endpoint will look something like the following:

```go
endpoint :=  "https://sql.dbtools.us-phoenix-1.oci.oraclecloud.com/20201005/ords/ocid1.databasetoolsconnection...change-me/_/sql"
```

What you see above is a dummy URL to get you started. Make sure you adjust the region and include the OCID of your Database Tools connection in the URL.

### Calling ORDS

Ignoring error checks for brevity and given the above, we just need to craft an HTTP POST request, sign it, and send it to the endpoint.

```go
	sql := bytes.NewBufferString("select sysdate from dual;")

	request, _ := http.NewRequest(http.MethodPost, endpoint, sql)

	// UTC() time is important, otherwise your requests will likely fail
	request.Header.Set("date", time.Now().UTC().Format(http.TimeFormat))
	request.Header.Set("content-type", "application/sql")

	// The DefaultConfigProvider assumes you have an .oci/config file setup
	// with a profile called [DEFAULT]
	signer := common.DefaultRequestSigner(common.DefaultConfigProvider())
	signer.Sign(request)
```

Once you have a signed request you can use Go's ```http.Client``` to send the request to the OCI endpoint. At this point, you are just doing normal web stuff.

```go

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	json, _ := io.ReadAll(response.Body)

	fmt.Println(string(json))

    // ... do stuff with the response from ORDS
```