# Specification of the Havel Intermediate Language (HVIL)

HVIL (Havel Intermediate Language) is a low-level intermediate representation used in the Havel compiler toolchain.
It is a typed, block-based language designed for program analysis and code generation.

## Types

HVIL is statically typed. The following types are supported:

- `1 byte`: 8-bit unsigned integer.
- `2 byte`: 16-bit unsigned integer.
- `4 byte`: 32-bit unsigned integer.
- `8 byte`: 64-bit unsigned integer.
- `ref`: A reference (pointer) to a memory location.

## Variables and Registers
 
 Variables in HVIL can be either:
 - **Declared Variables**: Declared in a `declare` block at the beginning of a function. They typically represent stack-allocated variables.
```
declare (
    n : 4 byte
);
```
 - **Registers**: Identified by a `$` prefix (e.g., `$n`). They are temporary values used within blocks.
```
$n : 4 byte = 0b1;
```
 
 ## Program Structure
 
 ### Functions
 A program consists of one or more functions. The entry point of a program is the function named `main`.
 Functions can take parameters and return multiple named values.

```
func main() {
  block entry {
      $val1 : 1 byte = 0b1;
      $val2 : 1 byte = 0b10;
      $res : 1 byte = local.add($val1, $val2);
      debug.dump($res);
  } return;
}

func add(a : 1 byte, b : 1 byte) -> (result : 1 byte) {
  block entry {
      result = alu.add_u(a, b);
  } return;
}
```
 
 ### Blocks
 A function body consists of an optional `declare` block followed by one or more `basic_block`s. The first block in a function is conventionally named `entry`.
 A basic block is a sequence of instructions that ends with a jump target.
 
 ### Jump Targets
 - `goto <label>;`: Unconditional jump to the block with the given label.
```
goto next_block;
```
 - `if <cond> then <label1> else <label2>;`: Conditional jump based on the value of `<cond>` (which must be a `1 byte` type).
```
if $is_less then loop_start else exit;
```
 - `return;`: Exit the function.
```
return;
```
 
 ## Built-in Functions

Built-in functions are organized into namespaces.

**Note:** Literals cannot be used as function arguments. They must be assigned to a register first.

### `alu` Namespace
Provides Arithmetic and Logic Unit operations.
 
 - `alu.add_u(a, b)`: Unsigned addition of `a` and `b`.
 - `alu.sub_u(a, b)`: Unsigned subtraction of `a` and `b`.
 - `alu.mul_u(a, b)`: Unsigned multiplication of `a` and `b`.
 - `alu.div_u(a, b)`: Unsigned division of `a` by `b`.
 - `alu.mod_u(a, b)`: Unsigned modulo of `a` by `b`.
 - `alu.lt_u(a, b)`: Returns `0b1` if `a < b` (unsigned), else `0b0`.
 - `alu.eq(a, b)`: Returns `0b1` if `a == b`, else `0b0`.
 - `alu.move(a)`: Returns the value of `a`.
 
**Example:**
```
$sum : 4 byte = alu.add_u($a, $b);
$ten : 4 byte = 0b1010;
$is_equal : 1 byte = alu.eq($sum, $ten);
```

 ### `mem` Namespace
 Provides memory access operations.
 
 - `mem.alloc(size)`: Allocates `size` bytes on the heap and returns a `ref`.
 - `mem.free(ptr)`: Frees the memory pointed to by `ptr`.
 - `mem.load(ptr)`: Loads a value from the memory location pointed to by `ptr`. The size of the load is determined by the result type.
 - `mem.store(ptr, value)`: Stores `value` at the memory location pointed to by `ptr`.
 - `mem.ptr(var)`: Returns a `ref` to the declared variable `var`.
 
**Example:**
```
$size : 1 byte = 0b100;
$ptr : ref = mem.alloc($size);
$val : 4 byte = 0b1010;
mem.store($ptr, $val);
$loaded : 4 byte = mem.load($ptr);
mem.free($ptr);
```

 ### `debug` Namespace
 Provides debugging operations.
 
 - `debug.dump(val)`: Dumps the value of `val` to the debug output.
 
**Example:**
```
debug.dump($val);
```

 ### `local` Namespace
 Used for calling other functions defined within the same program.
 
 - `local.<func_name>(args...)`: Calls the function `<func_name>`.

**Example:**
```
$res : 4 byte = local.my_function($arg1, $arg2);
```
