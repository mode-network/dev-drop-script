# dev-drop-script
generate transactions batch to send via safe wallet

## To build the executable
```shell
go build -o create-dev-drop ./cmd

```

## To run the script
```
cp .env.example .env
# set the envvars in .env file
# make sure you have the input csv file ready 
# run the generator 
./create-dev-drop sample-photons.csv transactions-batch.json

```