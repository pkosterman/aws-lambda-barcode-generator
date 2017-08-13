AWS Lambda Barcode-generator
===========

This is a sample project to demonstrate usage of Go for AWS Lambda. 

<a href="https://aws.amazon.com/lambda/">AWS Lambda</a> is a cloud computing service that lets you run code without provisioning or managing servers. AWS Lambda executes your code only when needed and scales automatically.

Currently AWS Lambda natively supports Java, Node.js, Python, and C#.

This project uses a Go Node.js wrapper to build a Go Lambda function that generates a and returns a barcode when triggered by AWS API Gateway. The barcode is returned as a base64 encoded PNG image string.

The Node.js wrapper keeps a Go process around to handle multiple invocations. The first time the function is run it will take a bit longer, but after that it greatly increases performance.

## Why Go?

I have been using Go for the past 3 years, it's our language of choice at <a href="https://passkit.com">PassKit</a> because of:

* <a href="https://hashnode.com/post/comparison-nodejs-php-c-go-python-and-ruby-cio352ydg000ym253frmfnt70">Speed</a>! Very fast & a perfect choice for CPU-intensive tasks.
* Quick & easy to master in a very short amount of time.
* Portability across platforms.
* Compiled binariers: plays nice with Docker.
* <a href="https://blog.golang.org/pipelines">Excellent concurreny primitives</a>. 	
* Well defined error handling patterns.
* Rich standard libraries.
* Standard code formatting / ease of maintenance.

## Inspiration
The Node.js wrapper used in this project is inspired by:
* <a href="https://github.com/jasonmoo/lambda_proc">lambda_proc</a>
* <a href="http://wiki.tcl.tk/44464">Amazon AWS Lambda &amp; Tcl</a>

## Full Demo & Instructions

Click <a href="https://blog.passkit.com/write-a-scalable-bar-code-generator-with-golang-aws-lambda">here</a> for a detailed article on how to set this up with AWS API Gateway and Route 53.

## Build
Clone this repo and cd into the project root then:

```bash
./build.sh
```

This will place a lambda.zip file into the build folder. You can update this zip file 
into AWS Lambda.

If you want to run `go test`, then you will also need to install lambda-test (Node command-line tool). I will
place a repo of this on github shortly - just need to rewrite some of the logic.