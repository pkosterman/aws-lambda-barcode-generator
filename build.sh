#!/bin/bash

# set -e
# set -x

# clean the build directory
echo "Cleaning build folder"
test -d "./build" && rm -rf build
mkdir ./build

echo "Setting Production Mode"
mv config/config.go config/config.tmp
sed -e "/const Develop = /s/.*/const Develop = false/g" config/config.tmp > config/config.go

echo "Compiling for Linux"
GOOS=linux go build -o ./build/main
cp index.js build

echo "Preparing Lambda bundle"
cd build
sed -e "/const exeName = /s/.*/const exeName = 'main';/g" index.js > index.tmp
sed -e "s/done(output, null);/done(output.errorMessage, null);/g" index.tmp > index.js
zip -r lambda.zip main index.js
echo "Uploading to AWS"

#TODO upload lambda.zip as lambda function
#aws lambda upload-function \
# --region $LAMBDA
# --function-name $LAMBDA_FUNCTION \
# --function-zip fileb://lambda.zip \
# --role $LAMBDA_ARN \
# --mode event \
# --handler handler \
# --timeout 30 \
# --runtime nodejs

echo "Cleaning up"
rm index.* main
mv ../config/config.tmp ../config/config.go
echo "Done!"
