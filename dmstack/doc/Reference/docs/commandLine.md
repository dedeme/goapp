## User Interface

The command line has the following structure:

```dmstack [Options] dmsProgram [-- dmsProgramOptions]```

  - Options:
    - `-d`: Running in 'debug mode'. The only difference between 'debug mode'
            and 'normal mode' is that in 'debug mode' are read and interpreted
            'at directives' and in 'normal mode' not.
  - dmsProgram: Path to ".dms" file. If this path not ends in ".dms", such
                extension will be added.
  - dmsProgramOptions: Arguments of ".dms" code. They can be read through
                       the procedure "sys.args".

## Examples

```
  dmstack myProg
```
```
  dmstack -d src/myProg.dms
```
```
  dmstack myProg -- -i "data/data 2.tb" -o result/book.txt
```
