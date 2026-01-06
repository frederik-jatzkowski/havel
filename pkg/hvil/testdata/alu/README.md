# ALU Operation Test Guidelines

Each alu operation should be tested by at least the following test cases:

## Happy Paths

1. A basic test for each possible input type
2. At least one test for behavior at the edge of the valid range

## Error Cases

1. For each parameter:
   1. One case with an ill-typed parameter, checking for type errors
2. One case for invalid result assignment, checking for type errors
3. If the operation can error (division by 0), one case for each error case