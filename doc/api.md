# Bigsense API

## Status

**Request:**

```bash
http localhost:8080/api/status
```

**Response:**

```json
HTTP/1.0 200 OK
Date: Fri, 13 Oct 2017 10:28:32 GMT
Server: WSGIServer/0.2 CPython/3.4.2
content-length: 16
content-type: application/json; charset=UTF-8

{
    "status": "OK"
}
```

## Users

### Crete User

**Request:** 

```bash
http POST localhost:8080/api/users Content-Type:application/json email="john.doe@email.com" password="123"
```

**Response:**

```
HTTP/1.0 201 Created
Date: Fri, 13 Oct 2017 10:15:39 GMT
Server: WSGIServer/0.2 CPython/3.4.2
content-length: 0
content-type: application/json; charset=UTF-8
location: a713f645-b7af-4676-9a32-b1eec49c381b
```

### Get User

First, get user's `jwt` token. Use the token to compose `Authorization` header entry value as `Bearer {jwt}`.

**Request:** 

```bash
http GET localhost:8080/api/users Content-Type:application/json Authorization:'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjp7ImVtYWlsIjoiam9obi5kb2VAZW1haWwuY29tIn0sImlhdCI6MTUxMjQ5MzE1OSwibmJmIjoxNTEyNDkzMTU5LCJleHAiOjE1MTI1Nzk1NTl9.6R8eJYD0GfoJnlgCp5VwmZgWvLgq4HIaRyFVLueLNYM'
```

**Response:**

```json
HTTP/1.0 200 OK
Date: Tue, 05 Dec 2017 17:26:41 GMT
Server: WSGIServer/0.2 CPython/3.6.3
content-length: 338
content-type: application/json; charset=UTF-8

{
    "_id": {
        "$oid": "5a201840e6c760446d7fd25e"
    },
    "balance": 25,
    "email": "john.doe@email.com",
    "subscriptions": [
        {
            "$oid": "5a2438a7e6c7603b639197cc"
        },
        {
            "$oid": "5a243dd4e6c76005e7e3c943"
        }
    ]
}
```

## JWT

### Get JWT token

JWT tokens are used to compose `Authorization` header entry value as `Bearer {jwt}` for the requests that need an authorization.

**Request:** 

```bash
http POST localhost:8080/api/token Content-Type:application/json email="john.doe@email.com" password="123"
```

**Response:**

```json
HTTP/1.0 200 OK
Date: Fri, 13 Oct 2017 10:15:59 GMT
Server: WSGIServer/0.2 CPython/3.4.2
content-length: 270
content-type: application/json; charset=UTF-8

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWJqZWN0Ijoiam9obi5kb2VAZW1haWwuY29tIiwiaXNzdWVkX2F0IjoiMjAxNy0xMC0xMyAxMDoxNTo1OS44MDE0NDAiLCJleHBpcmVzX2F0IjoiMjAxNy0xMC0xMyAyMDoxNTo1OS44MDE0NjQiLCJpc3N1ZXIiOiJOT0tJQSJ9.xMZOWifI4JfF4xJTlpJFOKMO1ak8b7O1GmnCXMfcjFw"
}
```

## Streams

### Create stream

First, get creator's `jwt` token. Use the token to compose `Authorization` header entry. The `name` and `url` values of the JSON have to be unique across the bigsense deployment. 

**Request:** 

```bash
http POST localhost:8080/api/streams Content-Type:application/json Authorization:'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjp7ImVtYWlsIjoiYmxvY2tzZW5zZUBub2tpYS5jb20ifSwiaWF0IjoxNTE0NDUwMTAzLCJuYmYiOjE1MTQ0NTAxMDMsImV4cCI6MTUxNDUzNjUwM30.ovyaOAk2EX8nwcKE3qZiBgJ0SdwS4GlKhfxT05C-8Xg' name='HR office temperature' type='floating point number' description='reads the current HR office room temperature' url='http://www.space.com/offices/HR/temperature' price='2' long='2.349014' lat='48.864716'
```

**Response:**

```json
HTTP/1.0 201 Created
Date: Tue, 05 Dec 2017 18:47:10 GMT
Server: WSGIServer/0.2 CPython/3.6.3
content-length: 2
content-type: application/json; charset=UTF-8

{}
```

### Get stream

In order to compose URL for the request, you need to know the stream's id. To compose the `Authorization` header entry, first get stream owner's `jwt` token.

**Request:** 

```bash
http GET localhost:8080/api/streams/5a201d93e6c760446d7fd260 Authorization:'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjp7ImVtYWlsIjoiam9obi5kb2VAZW1haWwuY29tIn0sImlhdCI6MTUxMjQ5NzY1MiwibmJmIjoxNTEyNDk3NjUyLCJleHAiOjE1MTI1ODQwNTJ9.mZKy_NfPecymNLK9vx_wkDcgrAFmmhfHqT-ImlB98EI'
```

**Response:**

```json
HTTP/1.0 200 OK
Date: Tue, 05 Dec 2017 18:59:16 GMT
Server: WSGIServer/0.2 CPython/3.6.3
content-length: 242
content-type: application/json; charset=UTF-8

{
    "data": {
        "_id": {
            "$oid": "5a201d93e6c760446d7fd260"
        },
        "description": "reads the number of car wheel rotations per minute",
        "name": "car wheel rotations",
        "owner": {
            "$oid": "5a201840e6c760446d7fd25e"
        },
        "price": 2,
        "type": "integer counter"
    }
}
```
### Purchase stream

To compose the `Authorization` header entry, first get stream owner's `jwt` token. JSON's `id` refers to stream's id.

**Request:** 

```bash
http POST localhost:8080/api/streams/purch Authorization:'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjp7ImVtYWlsIjoiam9obi5kb2VAZW1haWwuY29tIn0sImlhdCI6MTUxMjUwMTI3NSwibmJmIjoxNTEyNTAxMjc1LCJleHAiOjE1MTI1ODc2NzV9.kyUm7KjP-2qwFdM7GHCJqTXyv-jPUSUrWaQYacjFdsw' id='5a201d93e6c760446d7fd260' hours:='2'
```

**Response:**

```json
HTTP/1.0 201 Created
Date: Tue, 05 Dec 2017 19:19:00 GMT
Server: WSGIServer/0.2 CPython/3.6.3
content-length: 2
content-type: application/json; charset=UTF-8

{}
```

### Search stream by geolocation

To compose the `Authorization` header entry, first get stream owner's `jwt` token. Query parameters are for latitude and longitude coords used for geofencing.

**Request:** 

```bash
http GET localhost:8080/api/streams/search Authorization:'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjp7ImVtYWlsIjoiam9obi5kb2VAZW1haWwuY29tIn0sImlhdCI6MTUxNTY1OTQ2MCwibmJmIjoxNTE1NjU5NDYwLCJleHAiOjE1MTU3NDU4NjB9.RoZ9QvsDsv82mQnMPsi2ryyaAD56jkZQtqbvzdVCSL4' type=='geo' x0==20 y0==20 x1==30 y1==20 x2==20 y2==40 x3==30 y3==40
```

**Response:**

```json
HTTP/1.0 200 OK
Date: Thu, 11 Jan 2018 08:42:26 GMT
Server: WSGIServer/0.2 CPython/3.6.3
content-length: 287
content-type: application/json; charset=UTF-8

{
    "data": [
        {
            "_id": {
                "$oid": "5a572113e6c76023f4dfd79a"
            },
            "description": "a continous stream of data",
            "longlat": {
                "coordinates": [
                    25.123,
                    35.123
                ],
                "type": "Point"
            },
            "name": "Stream 05",
            "owner": {
                "$oid": "5a57205ee6c76024dd3b8bfb"
            },
            "price": 2,
            "type": "data type",
            "url": "/stream/02"
        }
    ]
}
```

## Buy

### Buy tokens from the superuser
To compose the `Authorization` header entry, first get buyer's `jwt` token.

**Request:** 

```bash
http POST localhost:8080/api/transfer/buy Authorization:'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjp7ImVtYWlsIjoiam9obi5kb2VAZW1haWwuY29tIn0sImlhdCI6MTUxMjUwMzA1NywibmJmIjoxNTEyNTAzMDU3LCJleHAiOjE1MTI1ODk0NTd9.otPIPZ75LKZWyayt75Bt6AXpuhTts4mOxAn8bhpsk8A' tokens:='2'
```

**Response:**

```json
HTTP/1.0 200 OK
Date: Tue, 05 Dec 2017 19:45:33 GMT
Server: WSGIServer/0.2 CPython/3.6.3
content-length: 2
content-type: application/json; charset=UTF-8

{}
```
