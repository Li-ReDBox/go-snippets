# Calculate days between two dates from scratch

The request is to calculate the distance in whole days between two dates, counting only
the days in between those dates, i.e. 01/01/2001 to 03/01/2001 yields “1”. The valid
date range is between 01/01/1900 and 31/12/2999, all other dates should be
rejected.

## Notes
There are two implmentatons of how a date is defined: [one](demo1/) is to use ints in an array, [another one](demo2/) is to define a type.
Using type is eaisier to understand as it is typed.

The executalbe can be build from [cmd](cmd/) directory.
1. [demo1](cmd/demo1/main.go) is for array implementaion;
2. [demo2](cmd/demo2/main.go) is for type implementation.

In both demos, there is a shell script: `run_demos.sh`. To run it:
```shell
# for example to run run_demos.sh in cmd/demo1
go build -o demo
sh run_demos.sh
```
The output should look like:
```shell
Days between 2/6/1983 and 22/6/1983 is 19
Days between 4/7/1984 and 25/12/1984 is 173
Days between 3/1/1989 and 3/8/1983 is 1979
Days between 1/3/1989 and 3/8/1983 is 2036
```

If using Go's `time` module, the calculation is redueced to [convert `Duration` to days](cmd/native/main.go). 