# The Tavor format

The [Tavor](/) format is an [EBNF-like notation](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_Form) which allows the definition of data (e.g. file formats and protocols) without the need of programming. It is the default format of the [Tavor platform](/) and supports every feature which the platform currently provides.

The format is Unicode text encoded in UTF-8 and consists of terminal and non-terminal symbols which are called <code>tokens</code> throughout the Tavor framework. An explanation of the general meaning can be found in the [What are tokens?](/#token) section.

## <a name="table-of-content"></a>Table of content

- [Token definition](#token-definition)
- [Terminal tokens](#terminal-tokens)
	+ [Numbers](#terminal-tokens-numbers)
	+ [Strings](#terminal-tokens-strings)
- [Concatenation](#concatenation)

TODO update this

## <a name="token-definition"></a>Token definition

Every token in the format belongs to a non-terminal token definition which consists of a unique case-sensitive name and its definition part. Both are separated by exactly one equal sign. Syntactical white spaces are ignored. Every token definition must be declared by default in one line. A line ends with a new line character.

To give an example, the following format declares the token <code>START</code> with the constant string token "Hello World" as its definition.

```tavor
START = "Hello World"
```

Token names have the following rules:
- Token names have to start with a letter.
- Token names can only consist of letters, digits and the underscore sign "_".
- Token names have to be unique in the format definition scope.

Additional to these rules it is not allowed to declare a token without any usage in the format definition scope except if it is the <code>START</code> token which is used as the entry point of the format, meaning it defines the beginning of the format. Hence, it is required for every format definition.

## <a name="terminal-tokens"></a>Terminal tokens

Terminal tokens are the constants of the Tavor format.

### <a name="terminal-tokens-numbers"></a>Numbers

Currently only positive decimal integers are allowed. They are written as a sequence of digits.

```tavor
START = 123
```

### <a name="terminal-tokens-strings"></a>Strings

Strings are character sequences between double quotes and can consist of any UTF8 encoded character except new lines, the double quote and the backslash which have to be escaped with a backslash.

```tavor
START = "The next word is \"quoted\" and here is a new line\n"
```

Since Tavor is using Go's text parser as foundation of its format parsing, the same rules for <code>interpreted string literals</code> apply. These rules can be looked up in [Go's language specification](https://golang.org/ref/spec#String_literals).

## <a name="concatenation"></a>Concatenation

Tokens in the definition part are automatically concatenated.

```tavor
START = "This is a string token and this " 123 " was a number token"
```

This example will be concatenated to the string "This is a string token and this 123 was a number token".

## <a name="multi-line"></a>Multi line token definitions

A token definition can be sometimes too long or poorly readable. It can be therefore split into multiple lines by using a comma before the newline character.

```tavor
START = "This",
        "is",
        "a",
        "multi line",
        "definition"
```

The token definition ends at the string "definition" since there is no comma before the new line character. This example also underlines that syntactical white spaces are ignored and can be used to make the format definition more human readable.

## <a name="comments"></a>Comments

The comments of the Tavor format follow the same rules as Go's comments which are specified in [Go's language specification](https://golang.org/ref/spec#Comments).

There are two types of comments:
- **Line comment** which starts with the character sequence <code>//</code> and ends at the next new line character.
- **General comment** which starts with the character sequence <code>/\*</code> and ends at the character sequence <code>\*/</code>. A general comment can contain new line characters.

```tavor
/*

This is a general comment
which can have
multiple lines

*/

START = "This is a string" // this is a line comment

// this is also a line comment
```

General comments can be used, like white space characters, between token definitions and tokens.

```tavor
START /* this is */ = "an" /* extreme */ "example" /* but
it should make it clear how general comments */ "work"
```

## <a name="embedding"></a>Token embedding

Non-terminal tokens can be embedded in the definition part by using the name of the referenced token. The following example embeds the token <code>String</code> into the <code>START</code> token.

```tavor
START = String

String = "this is a string"
```

Token names declared in the global scope of a format definition can be used throughout the format regardless of their declaration position.

Terminal and non-terminal tokens can be mixed.

```tavor
First  = "1."
Second = "2."
Third  = "3."

START = First ", " Second " and " Third
```

-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------

# TODO rewrite everything down below

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
Digits  = [\d]
Hex  = [\x20]
Unicode = [\x{10FFFF}] // Up to 8 Hex digits
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

