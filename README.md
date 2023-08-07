# taxyfare-cli
## What is this about
A taxy fare calculator CLI application that calculates the fare and simulates taxy distance meter & fare meter.

add new record
```
add 00:00:00.123 100.0
add 00:01:00.123 300.0
```

get current fare
```
sum
```

## Assumptions and limitations
- calculations based on distance only, meaning if trip distance value is stuck but elapsed time keep changing, the fare will constants.
- program will not exit if no error, unless got terminated explicitly .

## How to use
### Manual
1. From the project root directory, run below commands to setup and build the application as executable
```
make setup
make build
``` 
2. After successfully run previous commands, run by executing the executable (the example provide taxyfare as executable name)
```
./taxyfare
```

### Docker
1. From the project root directory, run below command to build the docker image
```
docker image build -t <image_name> .
```
2. Run the program inside the docker container using docker interactive mode
```
docker run -it <image_name>
```
3. Execute the command from inside docker image container


### Test scenarios
To run all test scenarios , run below command
```
make test
```