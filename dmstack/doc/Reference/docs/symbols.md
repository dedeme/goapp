A symbol is any sequence of ASCII characters which integer value is greater
than 32 (space) and different of `,`(comma), `:`(colon) and `;`(semicolon).

Symbols that contains `{`, `}`, `[`, `]`, `(` or `)`, apart of themselves, are
not allowed.

## Index

There is an index where every symbol is registered.

To convert a string in a symbol is necessary call `symbol.New` that returns
returns the corresponding symbol (If the string is already registered
it returns its symbol, otherwise, it adds the string to index, creates a new
symbol and returns it).

There are a number of symbols intially registered: if, else, nop, +, ...

## Symbols With Dot

A symbol which does not start whith `this.` only can have one dot. If it start
with `this.`, is allowed:

  - one dot. That of `this.`
  - two dots. That of `this.` and a final dot (e.g. `this.fn.`). This
    construction is thought to make "private" procedures.

If the dot is between two blanks, it is consedered a blank too.

---

- If the symbol start with `this.`, the next replacement is done:

> `this.right` -> `this right`

> The symbol `right` can have a dot if it is the last character.

---

- If the symbol start with a dot the following replacement is done:

> `.right` -> `"right" map get`

---

- If the symbol ends with a dot, it is treated like another normal symbol.

---

- If the dot is inner the symbol the next replacement is done:

> `left.right` ->  `left right`

> In this case if `left` is an import symbol, it will be substituted  by its
> path.

> Note that, in general, you ever can write `left right` directly.

---

## Symbols With `! + Int`

These symbols are sustituted by:

> `!number` -> `right lst get`

`number` must be an not negative Int (e.g. `l !3` -> `l 3 lst get`).

## @ Symbols

They are symbols wich start with `@`. Its semantics is shown in
[Stack Contol](../stackControl).

**@**

It is substituted by:

> `@` -> `stk empty?`

**@?Types**

It is substituted by:

> `@?types` -> `"types" stackCheckSymbol`

Where

> `types` is a types defition (e.g. `@?is`)

and

> `stackCheckSymbol` is the reserverd symbol `StackCheck` = `<= @?>`

**@+Types**

It is substituted by

> `@+types` -> `"types" stackOpenSymbol`

Where

> `types` is a types defition (e.g. `@+is`)

and

> `stackOpenSymbol` is the reserverd symbol `StackOpen` = `<= @+>`

**@-Types**

It is substituted by

> `@-types` -> `"types" stackCloseSymbol`

Where

> `types` is a types defition (e.g. `@-is`)

and

> `stackCloseSymbol` is the reserverd symbol `StackClose` = `<= @->`

**@Types**

It is substituted by

> `@types` -> `"types" stackSymbol`

Where

> `types` is a types defition (e.g. `@is`)

and

> `stackSymbol` is the reserverd symbol `Stack` = `<= @>`

## "import" Symbol

This symbol receive a special treatment when reading.
See [Imports](../imports)

## "this" Symbol

This symbol is substituted by the module path symbol.

## Reserverd Symbols

**`=` Symbol**

Used after a new symbol to save a token. If the symbol is not a new one,
an error will be raised.

**`=>` Symbol**

Used after a new symbol to save an exported token. If the symbol is not a new
one, an error will be raised.

**Global Module**

See [API](http://localhost/dmcgi/DmsDoc/?_system@global).

**Other Primitive Modules**

## # Symbols

Dmstack virtual machine keep a counter for singling symbols.

The symbol `#` increments the counter.

A symbol type `xxx#` is replaced by `xxx·counter` (e.g. `xxx·183`).

## Other Symbols

They are processed in two steps:

1. Verification if the symbol is an import. In this case it is processed as
   such import.
2. It it is not an import is normaly processed.


