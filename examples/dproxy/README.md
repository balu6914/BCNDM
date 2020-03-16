# dProxy curl examples

Assuming that you used docker-compose-dproxy.yml to configure dproxy, you can use commands below to test it.

When you execute docker-compose up, it will create `dproxy-localfiles` directory on your local filesystem.
Every file you put inside, will be accessible by dProxy.
In the example below, a file in that directory is fetched by dProxy

### Files on local filesystem

#### Download files

To create a token for a filesystem resource which will last for 1 hour and have unlimited number of fetches, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/register -d '{"url": "<your-filename-here>","ttl": 1,"max_calls": 0}'```

- Replace `<your-filename-here>` with the the filename of the file you put in `dproxy-localfiles` directory

Then execute curl command.


Once you have the token in curl response use the following curl template.

```curl -H "Authorization: <your-token>" http://localhost:9090/fs -o <downloaded-filename>```

- Replace the `<your-token>` part below.

- Also replace `<downloaded-filename>` with the desired filename.
 
Then execute curl command.
 
dProxy will deliver your file contents and curl will save it with the specified filename.

#### Upload files

To create a token to store HTTP PUT payload to a file which will last for 1 hour and have unlimited number of uploads, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/register -d '{"url": "<your-filename-here>","ttl": 1,"max_calls": 0}'```

- Replace `<your-filename-here>` with the the filename of the file you want to be created in the `dproxy-localfiles` directory. Please note that if file already exists, it will be overwritten with every HTTP PUT request

Then execute curl command.


Once you have the token in curl response use the following curl template.

```curl -X PUT -H "Authorization: <your-token>" -T <file-for-upload> http://localhost:9090/fs```

- Replace the `<your-token>` part below.

- Also replace `<file-for-upload>` with the desired filename.


dProxy will save contents of your payload to the filename specified in the token.


### Files on the network

To create a token for a resource on the network which will last for 1 hour and have unlimited number of fetches, use the following curl request template.

```curl -X POST -H "Content-Type: application/json" http://localhost:9090/api/register -d '{"url": "http://<your-hostname-here>/<path-to-your-resource>","ttl": 1,"max_calls": 0}'```

- Replace `<your-hostname-here>` with hostname of your server.

- Replace `<path-to-your-resource>` with path to the resource on your server

Then execute curl command.

Once you have the token in curl response use the following curl template.

```curl -H "Authorization: <your-token>" http://localhost:9090/http -o <downloaded-filename>```
 
- Replace the `<your-token>` part below.

- Also replace `<downloaded-filename>` with the desired filename.

Then execute curl command. 
  
dProxy will deliver your resource contents and curl will save it to the specified filename. 
