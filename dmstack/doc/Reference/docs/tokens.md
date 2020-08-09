
Tokens are syntax units. Each token is a serie of valid characters separated
by blanks.

Every token has a type and a position.

The type determines which operations can be done with it.

The position is a reference to the module (file) and line in which appears.

## Types

There are the following types of tokens and its literals:

  - **Bool**: Only two values. (`true`, `false`).
  - **Int**: [-]digits. This literals have not decimal point. The character '_'
             is allowed to grouping. (`0`, `34`, `-115`, `2_345 -> 2345`)
  - **Float**: [-]digits and decimal point. The character '_' is allowed to
               grouping. (`0.0`, `34.41`, `-115.16`, `2_345.2` -> `2345.2`)
  - **String**: It has two formats:
    - `"`_characters_`"`. This format can not be multiline, quotes and slashes
      should be escaped and allows other escaped symbols (\n\t...).
      (`""`, `"abc"`, `"a\"b\"d"`, `"33€"`)
    - `` ` ``_characters_``  ` ``. This format is multiline, quotes and slashes
      have not to be escaped. No escape symbol is processed. See more in
      [Literals-String](../literals/#string).
  - **Procedure**: `(`_token_ _token_ _..._`)`.
                   (`()`, `(1 == ("a") ("b") elif)`, `(2 +)`).
  - **List**: `[`_token_`,`_token_`,`_..._`]`. Commas are optional.
              (`[]`, `[1, "a", true]`, `[1 "a" true]`). Every element is
              evaluated with `data`.
  - **Map**: `{`_key_`:`_token_`,`_key_`:`_token_`,`_..._`}`. _key_ must
              evaluate to String. Semicolons and commas are optional.
              (`{}`, `{"a": 1, "b": true}`, `{"a" 1 "b" true}`). Every element
              is evaluated with `data`.
  - **Symbol**: _characters_. They can not contain `{}[]():;,` nor can start
                with `"-` or a digit (0-9). Neither can be `true` or `false`.
  - **Native**: Used to represent native objects. They have not literals and
                must be created programmatically.

