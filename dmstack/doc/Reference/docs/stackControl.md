Stack control has a symbol ever allowed (`@?xxx`) and others (`@`, `@xxx`,
`@+xxx`, `@-xxx`) only actives if `dmsatck` has been called with the option
'-d'.

### Stack Sequence Descriptor

It is a sequence of characters which especify value types in the machine stack.
A sequence descriptor is a sequence (possibly empty) of descriptor of type.

For example:

  - `is` specify `Int-String`-TopStack.
  - `b<= File>i` specfiy `Bool-Object File-Int`-TopStack

Valid values for descriptors are:

  * `b` Boolean
  * `i` Integer
  * `f` Float
  * `s` String
  * `p` Procedure
  * `l` List
  * `m` Map
  * `y` Symbol
  * `<= Ob>` Native object "Ob"

### @?xxx

The reader replaces it by `"xxx" <= @?>`, where "xxx" is a sequence descriptor.

`<= @?>` returns true if the machine stack last positions match "xxx".

For example:

---

```c
4 "a" @?is
```

Returns `true`.

---

`@?` without sequence descriptor (an empty sequence descriptor) returns `true`
if the machine stack is empty.

### @+xxx

Only active with the option '-d'.

The reader replaces it by `"xxx" <= @+>`, where "xxx" is a sequence descriptor.

`<= @!>` throws an exception if the machine stack last positions does not
match "xxx". Otherwise push the symbol `<= @!>` on the stack before the values
which match the sequence descriptor.

`@+` without sequence descriptor (an empty sequence descriptor) throws an
exception if the machine stack is empty. Otherwise push the symbol `<= @!>`
on the stack.

### @-xxx

Only active with the option '-d'.

The reader replaces it by `"xxx" <= @->`, where "xxx" is a sequence descriptor.

`<= @!>` throws an exception if the machine stack positions from `<= @!>`
does notmatch "xxx" or if `<= @!>` is not found. Otherwise remove the symbol
`<= @!>`.

`@-` without sequence descriptor (an empty sequence descriptor) throws an
exception if the machine stack is empty  or if `<= @!>` is not found. Otherwise
remove the symbol `<= @!>`.

### @xxx

Only active with the option '-d'.

The reader replaces it by `"xxx" <= @>`, where "xxx" is a sequence descriptor.

`<= @>` throws an exception if the machine stack last positions does not
match "xxx".

`@` without sequence descriptor (an empty sequence descriptor) throws an
exception if the machine stack is empty.

