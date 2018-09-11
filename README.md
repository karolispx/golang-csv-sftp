# Generate a CSV file and upload it to the server through SFTP.

It took me a while to get these 2 parts to work on my project so I decided to share my solution in case someone else is looking for these 2 features.

## Code put together using these 2 resources:
* Generating CSV file: http://www.golangprograms.com/sample-program-to-create-csv-and-write-data.html
* SFTP upload: https://github.com/pkg/sftp/issues/154

You will need to get the following packages to make it work:
* `go get -u github.com/pkg/sftp`
* `go get -u golang.org/x/crypto/ssh`
