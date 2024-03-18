# IR File Upload Tool

This tool will zip and upload all files within the cwd to a speficied S3 bucket. The tool will ignore the executable and the zip file it creates when adding all of the files to the archive.

## Usage
Run the compiled executable.

## Change Credentials/Bucket
All of the AWS credentials are stored in `.env`.  
`example.env` contains the necessary template you can use to load new credentials into and compile the program again. just save it as `.env`

### Compiling the Tool
run `go build awsupload.go`to generate an executable.

