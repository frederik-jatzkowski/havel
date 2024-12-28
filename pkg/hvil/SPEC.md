# Specification of the Havel Intermediate Language (HVIL)

## Examples

### n-th fibonacci number

```
func main() {
    block entry {
        r0 : 4 byte = 0x1101;
        r1 : 4 byte = fib(r0);
        builtin.print_i_32(r1);
    } => return;
}


// returns the a0-ths fibonacci number
func fib(a0 : 4 byte) : (v0 : 4 byte) {
    declare (
        s1 : 4 byte,
        s2 : 4 byte,
        s3 : 4 byte
    );
    
    block entry {
        s1 = 0x0;
        s2 = 0x0;
        s3 = 0x1;
        v0 = 0x0;
    } => b1;
    
    block b1 {
        r0 : 4 byte = i_lt(s1, a0);
    } => r0 ? b2 : b3;
    
    block b2 {
        v0 = i_add(s2, s3);
        s2 = s3; s3 = v0;
        r1 : 4 byte = 0x1;
        s1 = i_add(s1, r1);
    } => b1;
    
    block b3 {} return;
}
```

## Formal Grammar

### General Definitions

```ebnf
identifier ::= letter, { letter | digit | "_" };
letter ::= a-z;
digit ::= 0-9;
hex_literal ::= 0x[0-9a-f]+;

type ::= ( digit, { digit } ) | ( "[", type, { ",", type }, "]" );
```

### Program Structure

```ebnf
program ::= function, { function }
```

### Function Structure

```ebnf
function ::= "func", identifier, function_head, function_body, ";";

function_head ::= "(", variable_declarations, ")", [ "=>", "(", variable_declarations, ")" ];

function_body ::= "{", [ declare_block ], "}";

declare_block ::= "declare", "(", variable_declarations, ")";

variable_declarations ::= [ variable_declaration, { ",", variable_declaration } ];
variable_declaration ::= identifier, ":", type;
```

### Block Structure

```ebnf
basic_block ::= "block", identifier, "{", instruction_list, "}", jump_target;
instruction_list ::= { instruction, ";" };
jump_target ::= "=>", ( identifier | identifier, "?", identifier, ":", identifier );
```

### Instruction Structure

```ebnf
instruction ::= instruction_literal;
instruction_literal ::= instruction_result, "=", hex_literal;
instruction_unop ::= instruction_result, "=", "(", identifier, ")";
instruction_binop ::= instruction_result, "=", "(", identifier, ",", identifier ")";
instruction_result ::= variable_declaration | identifier;
```

## Semantics

### Predefined Builtin Functions

#### Debug
