# Error Codes
If the payload POSTed to the Inbox does not contain well-formed JSON, then the HTTP error code returned is "400 Bad Request".

If the payload POSTed to the Inbox contains well-formed JSON, but it is not valid Notify JSON-LD, then the error code is "422 Unprocessable Entity".