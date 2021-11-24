# hipeople-coding-challenge
REST API written in GoLang, using only standard libraries. It allows users to upload images and access them later.
## Instructions
1. Have go installed
2. Clone the repository
3. Run the server by:
    1. Open a terminal session in the project root folder
    2. Start the server locally at port 8080 by running `go run .`
4. Sending requests can be done using the terminal (cURL) or using Postman
5. To execute the unit tests, you can run `go test ./...` from the project root folder
## Functionality

### Upload
To upload an image you have to specify the path to it in the request.

If successful, you will receive the image ID in the response.
```
curl --location --request POST 'localhost:8080/images/' \
--form 'imageFile=@"$IMAGEPATH"'
```
### List
To simply get all image IDs, use this endpoint.

If successful, you will receive them in the response as a coma separated list.
```
curl --location --request GET 'localhost:8080/images/'
```
### Get by ID
To retrieve an image, you need to know its ID. A valid ID is of integer type.

If successful, you will receive the image in the response and you will be able to download it from there.
```
curl --location --request GET 'localhost:8080/images/566274931'
```