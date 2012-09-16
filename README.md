# Nostra

Nostra is a JSON response system that runs on your servers to retrieve
status-related information.

## Query Format

Open up a TCP session with Nostra on its given port and query as such:

    {
		    "version" : "1.0",
				"params" : [
				    "hostname",
						"time"
				]
		}

## Responses

Nostra will respond with data of the form:

    {
        "code" : 0,
				"data" : {
				    "hostname" : "athena",
						"time" : "2012-09-16T17:20:54-04:00"
				}
		}

