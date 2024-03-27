# dicom-service
Service that is able to accept and store an uploaded DICOM file, extract and return any DICOM header attribute based on a DICOM Tag as a query parameter, and finally convert the file into a PNG for browser-based viewing.

## How to run
```
docker-compose up --build
```

## Upload

Files are uploaded to `/images/dicom/`. We choose to accept files 10MB or less.

XRAY
```
curl -i -X POST -F "file=@images/raw/XRAY/DICOM/PA000001/ST000001/SE000001/IM000001" http://localhost:8080/v1/upload
```

MRI
```
curl -i -X POST -F "file=@images/raw/MRI/PA000001/ST000001/SE000001/IM000001" http://localhost:8080/v1/upload

```

## Extract

```
curl -i -X GET "http://localhost:8080/v1/extract?filename=/go/src/app/images/dicom/IM000001.dcm&tag=(0010,0010)"
```

Output:
```
Value for tag (0010,0010): [NAYYAR^HARSH]
```

## Convert

PNG for browser-based viewing and saved to `/images/dicom/`

```
curl -i -X GET "http://localhost:8080/v1/convert?filename=/go/src/app/images/dicom/IM000001.dcm"

```

Streaming is on by default for files over 1MB, option to turn on for smaller files

```
curl -i -X GET "http://localhost:8080/v1/convert?filename=/go/src/app/images/dicom/IM000001.dcm&streaming=true"
```

## Unit Test
Some test coverage to show how we would write unit tests for each handler function. Limited to extract handler due to time constraints.

```
go test ./... -v
```