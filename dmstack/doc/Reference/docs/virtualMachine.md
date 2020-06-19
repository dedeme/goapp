## Structure

A virtual machine processes a token of type Procedure.

Each virtual machine have:

  - A stack (implemented as a pointer to a slice of tokens).
  - A heap (implemented as a map from symbol to token).
  - A list of modules imported (implemented as a slice of symbols).

## Working

The virtual machine process a procedure token by token with the following
steps:

- The token is a symbol:

    - Is `=`?: Raise an error.

    - Is `&`?: Skip the token.

    - Is a global symbol?: It is run.

    - It is a primitive module?: Next token, that must be a symbol, is read
      and run.

    - It is an imported module?: Next token, that must be a symbol, is read
      and then:
        - If the last symbol is defined in the import:
            - If it references a procedure:
                - If an `&` fallows: The procedure is pushed into the stack.

                - Otherwise: The procedure is run.

            - Otherwise the token referenced is pushed into the stack.

        - Otherwise an error is raised.

    - It is a module in heap?:
        - If it references a procedure:
            - If an `&` fallows: The procedure is pushed into the stack.

            - Otherwise: The procedure is run.

        - Otherwise the token referenced is pushed into the stack.

    - It is followed by `=`?: The token is put into the heap.

    - Otherwise an error is raised ("unknown symbol")

- Otherwise is pushed into the stack.

---

NOTE: Observe that is not possible to redefine (push into heap) a symbol.

---

## Normal And Isolate Virtual Machine

A normal virtual machine share its stack with the virtual machine which
call it.

An isolate virtual machine has its own stack.

## Heap Access

Virtual machines create an call other virtual machines to execute procedures.

This sets a hierarchy m1 -> m2 -> ... -> mN.

When a machine finds reference for a symbol, it accesses orderly the heap of
every ancestor.

A problem can apear when a procedure is send to be executed by other virtual
machine. Let see the following program.

```c
( fn =
  5 n =
  5 fn puts
) sub = // Execute fn with 'n' as argument and shows the result.

( 3 s =
  (s +) sub
) pr0 =

( 3 n =
  (n +) sub
) pr1 =

pr0 // Shows 8 as expected.
pr1 // Shows unexpectedly 10.

```

`pr1` fails because when in `sub` the reference of `n` is searched, the
value `5` is found before `3`.

For fix that there are two ways (see [Symbol #](../symbols/#symbols_1)):

1

```c
( fn =
  # 5 n# =  // Make n unique.
  n# fn puts
) sub =

( 3 n =
  (n +) sub
) pr1 =

pr1 // Shows 8.
```

2
```c
( fn =
  5 n =
  n fn puts
) sub =

( # 3 n# = // Make n unique.
  (n# +) sub
) pr1 =

pr1 // Shows 8.
```

The first way is better because avoid complications to the procedure client.

Also is posible a more complicated form:

```c
( fn =
  # 5 n# =  // Make n unique.
  n# fn puts
) sub =

( # 3 n# = // Make n unique.
  (n# +) sub
) pr1 =

pr1 // Shows 8.
```



