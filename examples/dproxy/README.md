# dProxy curl examples

dProxy can use tokens for user authorization. These tokens can be supplied either in the HTTP Authorization header or as a part of an URL.
Below are examples for both ways.

When you execute docker-compose up, it will create `dproxy-localfiles` directory on your local filesystem.
Every file you put inside, will be accessible by dProxy.

dProxy can fetch resources from an HTTP server or from a local filesystem.
dProxy can also store uploaded files to a local filesystem.

## dProxy curl examples for token in the URL

#### Download files

To create a url_with_token for a filesystem resource which will last for 1 hour and have unlimited number of fetches, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/register -d '{"url": "<your-filename-here>","ttl": 1,"max_calls": 0}'```

- Replace `<your-filename-here>` with the the filename of the file you put in `dproxy-localfiles` directory

Once you have the url_with_token in curl response use it in the next curl request.

```curl <url_with_token> -o <downloaded-filename>```

- Replace the `<url_with_token>` part with the url you got from /api/register curl call.

- Also replace `<downloaded-filename>` with the desired filename.

dProxy will deliver your file contents and curl will save it with the specified filename.

#### Upload files

To create a url_with_token to store HTTP PUT/POST payload to a file which will last for 1 hour and have unlimited number of uploads, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/register -d '{"url": "<your-filename-here>","ttl": 1,"max_calls": 0}'```

- Replace `<your-filename-here>` with the the filename of the file you want to be created in the `dproxy-localfiles` directory. Please note that if file already exists, it will be overwritten with every HTTP PUT request

Then execute curl command.

##### PUT example (overwrites existing file)

Once you have the url in curl response use the following curl template.

```curl -X PUT -T <file-for-upload> <url_with_token>```

- Replace the `<url_token>` part with the url you got from /api/register curl call.

- Also replace `<file-for-upload>` with the desired filename.

dProxy will save contents of your payload to the filename specified in the url_token. Every time you issue same PUT command, dProxy will overwrite current file with new contents that you upload.

##### POST example (creates new file each time)

With the same url_with_token you can use the following curl template for POST requests.
The difference from PUT command is that, every time you issue POST command, a new file will be created with the filename+timestamp. So previous uploads will not be overwritten.

```curl -X POST -T <file-for-upload> <url_with_token>```

- Replace the `<url_with_token>` with the url response you got from /api/register call.

- Also replace `<file-for-upload>` with the desired filename.

dProxy will save contents of your payload to the filename specified in the token and will suffix it with the timestamp.


### Files on the network (HTTP server)

To create a url_with_token for a resource on the network which will last for 1 hour and have unlimited number of fetches, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/register -d '{"url": "http://<your-hostname-here>/<path-to-your-resource>","ttl": 1,"max_calls": 0}'```

- Replace `<your-hostname-here>` with hostname of your server.

- Replace `<path-to-your-resource>` with path to the resource on your server


Once you have the url_with_token in curl response use the following curl template.

```curl <url_with_token> -o <downloaded-filename>```
 
- Replace the `<url_with_token>` with the url you got from previous response.

- Also replace `<downloaded-filename>` with the desired filename.
  
dProxy will deliver your resource contents and curl will save it to the specified filename. 



## dProxy curl examples for token in the HTTP Authorization header

### Files on local filesystem

#### Download files

To create a token for a filesystem resource which will last for 1 hour and have unlimited number of fetches, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/token -d '{"url": "<your-filename-here>","ttl": 1,"max_calls": 0}'```

- Replace `<your-filename-here>` with the the filename of the file you put in `dproxy-localfiles` directory

Once you have the token in curl response use the following curl template.

```curl -H "Authorization: <your-token>" http://localhost:9090/fs -o <downloaded-filename>```

- Replace the `<your-token>` part below.

- Also replace `<downloaded-filename>` with the desired filename.
 
dProxy will deliver your file contents and curl will save it with the specified filename.

#### Upload files

To create a token to store HTTP PUT/POST payload to a file which will last for 1 hour and have unlimited number of uploads, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/token -d '{"url": "<your-filename-here>","ttl": 1,"max_calls": 0}'```

- Replace `<your-filename-here>` with the the filename of the file you want to be created in the `dproxy-localfiles` directory. Please note that if file already exists, it will be overwritten with every HTTP PUT request


##### PUT example (overwrites existing file)

Once you have the token in curl response use the following curl template.

```curl -X PUT -H "Authorization: <your-token>" -T <file-for-upload> http://localhost:9090/fs```

- Replace the `<your-token>` with the token you got from /api/token curl call.

- Also replace `<file-for-upload>` with the desired filename.

dProxy will save contents of your payload to the filename specified in the token. Every time you issue same PUT command, dProxy will overwrite current file with new contents that you upload.

##### POST example (creates new file each time)

With the same token you can use the following curl template for POST requests.
The difference from PUT command is that, every time you issue POST command, a new file will be created with the filename+timestamp. So previous uploads will not be overwritten.

```curl -X POST -H "Authorization: <your-token>" -T <file-for-upload> http://localhost:9090/fs```

- Replace the `<your-token>` with the token you got from /api/token.

- Also replace `<file-for-upload>` with the desired filename.

dProxy will save contents of your payload to the filename specified in the token and will suffix it with the timestamp.

### Files on the network (HTTP server)

To create a token for a resource on the network which will last for 1 hour and have unlimited number of fetches, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/token -d '{"url": "http://<your-hostname-here>/<path-to-your-resource>","ttl": 1,"max_calls": 0}'```

- Replace `<your-hostname-here>` with hostname of your server.

- Replace `<path-to-your-resource>` with path to the resource on your server

Once you have the token in curl response use the following curl template.

```curl -H "Authorization: <your-token>" http://localhost:9090/http -o <downloaded-filename>```
 
- Replace the `<your-token>` part below.

- Also replace `<downloaded-filename>` with the desired filename.
  
dProxy will deliver your resource contents and curl will save it to the specified filename. 

