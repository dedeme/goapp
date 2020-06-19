## Blank

It is a character (ch) such that :

`ch <= ' ' || ch == ';' || ch == ':' || ch == ','`

Blanks are separators of tokens and are skipped by the reader.

## Comment

There are two types:

  - Line comment. Start with `//` and extends to end of line.
  - Block comment. Starts with `/*` and ends with `*/`. There are not
    nested block comments.

Examples:

```c
35 : x =
// This is a line comment.
12 : y = // This other line comment.
/*
  This is a block command.
*/
x y + /* This is a block command too */ puts
```

---

_NOTE_ :

Comments must be separated with a blank from the precedent token. However the
next token can be without separation.

This is incorrect:

```c
12 : y =// This other line comment.
x y +/* This is a block command too */ puts
```

But this is correct:

```c
x y + /* This is a block command too */puts
```

---

## Boolean

Values `true` and `false`

## Integer

Its underline value is an integer of 64 bytes.

#### Decimal Integer

Digits sequence that can start with `-`. (`0`, `-0`, `142`, `-4532`).

#### Hexadecimal Integer

Hexadecimal digits (upper or lowercase) which must start with '0x'.

Its digits part can start with `-` or `+`. (`0x2f`, `0x-AA`, `0x+16`)

## Float

Its underline value is an float of 64 bytes.

It is a sequence of digits whih must contain a decimal point `.` and can
have `e` or `E` for scientific notation. (`-12.45`, `-0.`, `-11.e1`, `23.4E-3`,
`45.e+2`.

---

> _NOTE_ :

> A float can not start with a point. In such case is interpreted as an object
  key.

> The decimal point is mandatory. For example `-11e1` throws an error.

---

## String

#### Normal Strings

They are sequences of UTF-8 symbols between quotes.

The following escape sequences are allowed: `\"`, `\\\\`, `\t`, `\n`, `\r`,
`\b`, `\f` and `\uXXXX`(Hexadecimal unicode character).

Charancter with integer value less than ' '(space) are not allowed. Therefore
this strings can not be extended more than one line.

#### "HereDoc" Strings

This string starts with `` `symbol\n`` and ends with ``symbol` ``. `symbol`
is optional.

This strings have not escape sequences and allows any type of UTF-8 character.

**Examples**

```
`
a multiline
string.`
```
produces `"a multiline\nstring"`.

```
`TX--
a `multiline`
string.--TX`
```
produces ``"a `multiline`\nstring"``.

But
```
`a multiline
string.`
```
produces an error.

#### String Interpolation

Both string classes allow interpolation. A interpolation is a expresion of
type `${value}`, where `value` is a expresion which is calculated with `data`
and must returns no value or one only value which procedure `toStr` can be
applied to.

**Examples**

`"2 + 2 = ${2 2 +}"` -> `2 + 2 = 4`

`'${}'` -> `''`

---

NOTE:

> - If you want to use the literal `${`, you should use some resourece lake spliting
the string or using `${$}{`.

> > For example:

> > `"The '$" "{' is a symbol" +` -> `The '${' is a symbol`

> > or

> > `"The '${$}{' is a symbol"` -> `The '${' is a symbol`

> - Nither you can use `}` inside the interpolation.

---

## Symbol

A symbol is any sequence of ASCII characters which integer value is greater
than 32 (space) and different of `,`(comma), `:`(colon) and `;`(semicolon).

Symbols that contains `{`, `}`, `[`, `]`, `(` or `)`, apart of themselves, are
not allowed.

See [Symbols](../symbols).

## Procedure

A procedure is a sequence of tokens between parentheses.

**Examples**

```c
(4 5 +)
```

```c
(
  4 5 +
    6 ==
)
```

## List

A list is a sequence of tokens between square brackets.

A list is read as a Procedure and adding the reserverd word `data`. That is
`[...]` is replaced by `(...) data`. As result tokens in the list will be
evaluated by the `dmstack machine`.

**Examples**

```c
[4 5 +]
```
that is the same as `[9]`.

```c
[
  4 5 +
    6 "a"
]
```
that is the same as `[9 6 "a"]`.

For better visualization is recomended to use commas:

```c
[4 5 +, 6, "a"]
```

## Map

A map is a sequence of tokens between curly brackets. Tokens are alternatively
"string - token" so that its number is even.

It is a list that will be converted to `map`. `{...}` is replaced by
`(...) data map.fromList`.

**Examples**

```c
{ "a" 7 6 + "b" true }
```
that is the same as `{"a" 13 "b" true}`.

It is recomended to use `:` and `,` to improve the visualization.

```c
{ "a": 7 6 +, "b": true }
```
or
```c
{
  "a": 7 6 +
  "b": true
}
```

## Program

It is a sequence of tokens.

**Example**

```c
"Hello World!" puts
```




