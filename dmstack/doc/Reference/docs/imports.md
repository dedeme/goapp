## Syntax

There are two forms:

`"modulePath" import`

and

`("modulePath" symbol) import`

Where:

_"modulePath"_ is the relative path of the file to import, without the `.dms`
extension (relative to the calling file). The path is transformed to its
canonical representation.

and

_symbol_ is the symbol to represent "modulePath". When it is not indicated,
the base of "modulePath" is used as symbol.

**Example**

We have the following structure,
```text
|- main.dms
|- lib
    |- inc.dms
```

If `inc.dms` has a procedure `print` with print an Int on screen, we can
call this procedure from `main.dms` in any of the next ways:

```c
"lib/inc" import
4 inc.print
```
```c
("lib/inc" i) import
4 i.print
```

#### Syntax Implementation

The reader has a map for each file from symbols to path.

Every import generates a new entry "symbol: path". After that following
cases of "symbol" are replaced by "path".

Note that a redefinition in the same file of "symbol" hiddes the previous one.

The efective importation is made by _dmstack machine_, when code is
processed.

## Working

- If the import has been imported, it is added to the list of imports of
  the current virtual machine.

- If the import is 'on way' an error is raised (cyclic import)

- Otherwaise the correponding module is read, runned in as isolate virtual
  machine, saved in the imports list and added to the list of imports of the
  current virtual machine.

#### OnWay

When a file import starts, its file path is marked as "on way" with the
function `imports.PutOnWay`.

When the import ends, such mark is removed with `imports.QuitOnWay`.

To test if a path is "on way' can be used `imports.IsOnWay`.


