# Tavor file format

## Comments

```
// single line comment
```

```
/* this comment can be a single line */
/*
    or multiple
    lines
    long
    and can be inlined too
*/
```

## Token definition

The left side of a token definition defines the name of the token. The right site the format of the token. Both sides are separated by "=" and end at the end of the line.

Naming convention for tokens:
* Token names have to start with a letter
* Token names can only consist of letters, digits and "_"
* Reserved token names are
    * START - parsing of data will start from this token. Is required for every format definition.

```
Token = "I am a constant string" 123
```

"I am a constant string" is a constant string. 123 is a constant number. They will be used as is. Everything is concatenated by default. So "I am a constant string" and 123 in this example are parsed as "I am a constant string123".

```
AnotherToken = Token // AnotherToken embeds the token "Token".
MultiLineToken = "a", // Token definitions can have multiple lines if there is a comma at the end of the line.
                 "b",
                 "c" // There is no comma at the end of the line which means that this token definitions ends here.
```

```
Umläüt = "Umlauts can be used since these definitions have to be in utf8"
Quoting = "\"this is quoted\"" // Backspaces is used as escape character.
```

### Alternations and grouping

```
Alternation = 1 | 2 | 3 // The token "Alternation" can hold either 1, 2 or 3.
SameAlternationAsShortage = [123] // This is the same as the "Alternation" token except it is much shorter to define.
AnotherAlternation = "a string" | [123] | Token // Alternations can hold every kind of token.
Grouping = ("old" | "new") "letter" // Everything between parenthesis is a group. The token "Grouping" can therefore hold "oldletter" or "newletter". Groups can be nested too.
Permutations = @(1 | 2 | 3) // Alternation groups can become permutation groups with the "@" right before the opening parenthesis. Each entry will be used once but the order is nonrelevant. For example the token "Permutations" can hold 123, 132, 213, 231, 312 or 321.
```

### Optionals and repeats

```
Optional = "i am not optional" ?("but hey i am optional!") // The constant string "but hey i am optional!" is optional.
RepeatAtLeastOnce = "text" +("me") // "me" will be repeated at least once.
OptionalRepeat = "text me" *("or me") // "me" can be repeated zero, one or more times.
RepeatExactlyTwice = "text" +2("me") // "me" is repeated exactly twice.
RepeatAtLeastTwice = "text" +2,("me") // "me" is repeated at least twice.
RepeatAtMostTenTimes = "text" +,10("me") // "me" is repeated at most ten times.
RepeatTwoToTenTimes = "text" +2,10("me") // "me" is repeated two to ten times.
```

### Character classes

```
Letters  = [abc]
```

### Token attributes

Token attributes can be used in token definitions by prepending a dollar sign to their name and separate the token name from the attribute by a dot.

```
Letters = *(Letter)
Letter = "a" | "b" | "c"
LetterCount = $Letters.Count // LetterCount then holds the count of the repeater Letters
```

Possible token attributes are:
* Count - Holds the count of this token. Must be a repeater.
* Index - Holds the index of a token. Must be a token of a repeater.
* Unique - Chooses at random a token of a repeater.

### Special tokens

Special tokens can be defined by prepending a dollar sign to their name. Special tokens do not have a format on their right side like regular tokens, instead arguments written as key-value pairs, which are separated by a colon, define the token. At least the "type" argument must be defined.

```
$Number = type: Int
Arithmetic = Number "+" Number
```

Possible arguments are:
* type - Defines the type of the token. Can be "Int" or "Sequence"

Additional (optional) arguments for each type are:
* "Int"
    * from - First integer value
    * to - Last integer value
* "Sequence"
    * start - First sequence value. Default is 1.
    * step - Increment of the sequence. Default is 1.

Possible attributes for each type are:
* "Int"
    * Value - The value of the Int
* "Sequence"
    * Next - Indicates the next value of the sequence.
    * Existing - Indicates an available value of the sequence in the whole data.
    * Reset - The sequence is reseted when this token is reached.

```
$Id = type: Sequence,
      start: 0,
      step: 2
NextId = $Id.Next
```

### Expressions

Expressions can be used on the right side of a token definition.

```
Sum = ${1 + 2 + 3} // Sum will be interpreted as 6

SomeIdOrMore = $Id.Existing | ${Id.Existing + 1}

DoubleTheCount = ${Letter.Count + Letter.Count}
```

### Variables

Every token on the right side of a definition can be saved into a variable.


```
START = "text"<var> Print

Print = <var>.Value
```

This will save the string "text" into the variable "var" without preventing the relay of the string to the output stream.

Since there are circumstances where a token should be just saved into a variable but not relayed to the output stream a second syntax can be used.

```
START = "text"<=var> Print

Print = <var>.Value
```

### Set operators

Some attributes can be combined with set operators. For example

```
$Id = type: Sequence

Pair = $Id.Next<id> " " ${Id.Existing not in (id)}
```

This will search through the existing sequenced IDs without the one saved in the variable "id".

### If, If else and else

```
START = Choose<var> Print

Choose = 1 | 2 | 3

Print = {if var.Value == 1} "var is one" {else if var.Value == 2} "var is two" {else} "var is three" {endif}
```

### Condition operators

* "=="

  ```
  Print = (1 | 2 | 3)<var> {if var.Value == 1} "var is 1" {else} "var is not 1" {endif}
  ```

* "defined"

  ```
  START = Print "save this text"<var> Print

  Print = {if defined var} "var is: " $var.Value {else} "var is not defined" {endif}
  ```
