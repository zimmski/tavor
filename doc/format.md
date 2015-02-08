# The Tavor format

The [Tavor](/) format is an [EBNF-like notation](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_Form) which allows the definition of data (e.g. file formats and protocols) without the need of programming. It is the default format of the [Tavor framework](/) and supports every feature which the framework currently provides.

The format is Unicode text encoded in UTF-8 and consists of terminal and non-terminal symbols which are called `tokens` throughout the Tavor framework. An explanation of the general meaning can be found in the [What are tokens?](/#token) section.

> **Note**: The Tavor format is not stable and changes slightly with every version of Tavor. This is necessary to make the format more powerful and easier to understand and write. Please [submit an issue](https://github.com/zimmski/tavor/issues/new) if you find flaws, ambiguous content or definitions, or something that could be easier defined.

Every example of this page is a complete and syntactical correct Tavor format file. The content of each example can be for instance saved into a file called `file.tavor` and then fuzzed with the Tavor binary. To get a a better understanding of the format it is advised to do this with every example.

```bash
tavor --format-file file.tavor fuzz
```

Since some examples have more than one permutation, meaning there is more than one possible fuzzing generation, it is advisable to use the `AllPermutations` fuzzing strategy to print out every possible permutation of the fuzzed format.

```bash
tavor --format-file file.tavor fuzz --strategy AllPermutations
```

## <a name="table-of-content"></a>Table of content

- [Token definition](#token-definition)
- [Terminal tokens](#terminal-tokens)
	+ [Numbers](#terminal-tokens-numbers)
	+ [Strings](#terminal-tokens-strings)
- [Concatenation](#concatenation)
- [Multi line token definitions](#multi-line)
- [Comments](#comments)
- [Token embedding](#embedding)
- [Alternation](#alternation)
- [Grouping](#grouping)
	+ [Optional group](#grouping-optional)
	+ [Repeat groups](#grouping-repeats)
	+ [Permutation group](#grouping-permutation)
- [Difference between token reference and token usage](#reference-usage)
- [Character classes](#character-classes)
	+ [Escape characters](#character-classes-escapes)
	+ [Ranges](#character-classes-ranges)
	+ [Special escape characters](#character-classes-special-escapes)
- [Token attributes](#attributes)
	+ [General attributes](#attributes-general)
	+ [Scope of attributes](#attributes-scope)
- [Typed tokens](#typed-tokens)
	+ [Type `Int`](#typed-tokens-Int)
	+ [Type `Sequence`](#typed-tokens-Sequence)
- [Expressions](#expressions)
	+ [Arithemtic operators](#expressions-arithmetic)
	+ [Graph operators (experimental)](#expressions-graph)
	+ [Set operators (experimental)](#expressions-set)
- [Variables](#variables)
	+ [ Token attributes](#variables-token-attributes)
	+ [Just-save operator](#variables-just-save)
- [Statements](#statements)
	+ [`if` statement](#statements-if)

## <a name="token-definition"></a>Token definition

Every token in the format belongs to a non-terminal token definition which consists of a unique case-sensitive name and its definition part. Both are separated by exactly one equal sign. Syntactical white spaces are ignored. Every token definition must be declared by default in one line. A line ends with a new line character.

To give an example, the following format declares the token `START` with the constant string token "Hello World" as its definition.

```tavor
START = "Hello World"
```

Token names have the following rules:

- Token names have to start with a letter.
- Token names can only consist of letters, digits and the underscore sign `_`.
- Token names have to be unique in the [global scope](#attributes-scope).

Additional to these rules it is not allowed to declare a token without any reference. Except if it is the `START` token which is used as the entry point of the format. Meaning it defines the beginning of the format and is therefore required for every format definition.

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

Since Tavor is using Go's text parser as foundation of its format parsing, the same rules for `interpreted string literals` apply. These rules can be looked up in [Go's language specification](https://golang.org/ref/spec#String_literals).

> **Note**: Empty strings are forbidden and lead to a format parse error. The reasons are explained in more detail in the [Repeat groups section](#grouping-repeats).

## <a name="concatenation"></a>Concatenation

Sequential tokens in the definition part are automatically concatenated.

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

The token definition ends at the constant string "definition" since there is no comma before the new line character. This example also underlines that syntactical white spaces are ignored and can be used to make the format definition more human readable.

## <a name="comments"></a>Comments

The comments of the Tavor format follow the same rules as Go's comments which are specified in [Go's language specification](https://golang.org/ref/spec#Comments).

There are two types of comments:

- **Line comments** start with the character sequence `//` and end at the next new line character.
- **General comments** start with the character sequence `/*` and end at the character sequence `*/`. A general comment can therefore contain new line characters.

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

Non-terminal tokens can be embedded in the definition part by using the name of the referenced token. The following example embeds the token `Text` into the `START` token.

```tavor
START = Text

Text = "This is some text"
```

Token names declared in the global scope of a format definition can be used throughout the format regardless of their declaration position.

Terminal and non-terminal tokens can be also mixed	.

```tavor
Dot = "."

First  = 1 Dot
Second = 2 Dot
Third  = 3 Dot

START = First ", " Second " and " Third
```

## <a name="alternation"></a>Alternation

Alternations are defined by the pipe character `|`. The following example defines that the token `START` can either hold `1`, `2` or `3`.

```tavor
START = 1 | 2 | 3
```

An alternation term has its own scope which means that a sequence of tokens can be used.

```tavor
START = 1 "green apple" | 2 "orange oranges" | 3 "yellow bananas"
```

Alternation terms can be empty which allows more advanced definitions of formats. For example the next definition defines the possibility of a loop.

```tavor
A = "a" A | B |
B = "b"

START = A
```

This example can either hold the strings "", "a", "b", "ab", "aab" or any amount of "a" characters ending with one or no "b" character.

## <a name="grouping"></a>Grouping

Tokens can be grouped using parenthesis beginning with the opening parenthesis `(` and ending with the closing parenthesis `)`. A group is a token on its own. This means that it can be mixed with other tokens. Additionally, a group starts a new scope between its parenthesis and can therefore hold a sequence of tokens. The tokens between the parenthesis are called the `group body`.

The following example declares that the token `START` either holds the string "old news" or "new news".

```tavor
START = ("old" | "new") " news"
```

Groups can be nested too. For example the following can be used to define that the `START` token can either hold "a", "b", "1" or "2".

```tavor
START = (("a" | "b") | (1 | 2))
```

An even more complicated example is the definition of an one to three digits integer.

```tavor
Digit = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

START = Digit | Digit Digit | Digit Digit Digit
```

This could be also written with the following format definition.

```tavor
Digit = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

START = Digit (Digit (Digit | ) | )
```

Group parenthesis can have modifiers which give the group additional abilities. The following sections will introduce these modifiers.

### <a name="grouping-optional"></a>Optional group

The optional group has the question mark `?` as modifier and allows the whole group token to be optional. In the next example the `START` token can hold either hold strings "funny" or "very funny".

```tavor
START = ?("very ") "funny"
```

### <a name="grouping-repeats"></a>Repeat groups

The default modifier for the repeat group is the plus character `+`. The repetition is executed by default at least once. In the next example the string "a" is repeated and the `START` token can therefore hold the strings "a", "aa", "aaa" or any amount of "a" characters.

```tavor
START = +("a")
```

> **Note:** It is forbidden to repeat an optional group or an alternation group with an optional term. The reason becomes obvious in terms of parsing and delta-debugging data. Since optional groups must not parse anything, they can repeatedly parse nothing and still conform to the defined format. This is a waste of resources and leads to an enormous amount of unneeded reducing steps while delta-debugging.
>
> For example the following leads to a format parse error.
>
> ```tavor
> START = +(?("a"))
> ```
>
> The following example does lead to an error too, since alternation groups with an optional term are also forbidden.
>
> ```tavor
> START = +("a" | )
> ```

Although the format definition allows the repetition to go on forever there are bounds since there is only a finite amount of memory available. The Tavor framework does set a maximum repetition which can be altered by the `--max-repeat` option of the Tavor binary or the `MaxRepeat` variable exported by the `github.com/zimmski/tavor` package.

If no maximum repetition is set the repetition modifier repeats by default from one to infinite which can be altered with arguments to the modifier. The next example repeats the string "a" exactly twice meaning the `START` token does only hold the string "aa".

```tavor
START = +2("a")
```

It is also possible to define a repetition range. The next example repeats the string "a" at least twice but at most 4 times. This means that the `START` token can either hold the strings "aa", "aaa" or "aaaa".

```tavor
START = +2,4("a")
```

The `from` and `to` arguments can be empty too. They are then set to their default values. For example the next definition repeats the string "a" at least once and at most 4 times.

```tavor
START = +,4("a")
```

The following example repeats the string "a" at least twice.

```tavor
START = +2,("a")
```

Since the repetition zero, once or more is very common the modifier `*` exists. In the next example the token `START` can either hold the string "a", "ab", "abb" or any amount of "b" characters prepended by an "a" character.

```tavor
START = "a" *("b")
```

### <a name="grouping-permutation"></a>Permutation group

The `@` is the permutation modifier which is combined with an alternation in the group body. Each alternation term will be executed exactly once but the order of execution is non-relevant. In the next example the `START` token can either hold 123, 132, 213, 231, 312 or 321.

```tavor
START = @(1 | 2 | 3)
```

## <a name="reference-usage"></a>Difference between token reference and token usage

The following example demonstrates the difference between a **token reference** and a **token usage**.

```tavor
Choice = "a" | "b" | "c"

List = +2(Choice)

START = "1. list: " List "\n",
        "2. list: " List "\n"
```

This format defines two tokens called `Choice` and `List`.

A **token reference** is the embedding of a token in a definition. There exists one token reference of `Choice`, which can be found in the `List` definition, and two for `List`, which are both in the `START` definition. Even though `Choice` is in a repeater group it is only referenced once.

A **token usage** is the execution of a token during an operation like fuzzing or delta-debugging. `List` has two token usages in this format while `Choice` has 4. Every `List` token does have two `Choice` usages because of the repeat group in the definition of `List`.

## <a name="character-classes"></a>Character classes

Character classes are a special kind of token and can be directly compared to character classes of regular expressions used in most programming languages such as Perl's implementation which is documented [here](http://perldoc.perl.org/perlre.html#Character-Classes-and-other-Special-Escapes). They behave like terminal tokens meaning that they cannot include others tokens but they are, unlike constant integers and constant strings, not single but multiple constants at once. A character class starts with the left bracket `[` and ends with the right bracket `]`. Character classes are like terminal tokens in that they are tokens on their own and can be therefore mixed with other tokens. The content between the brackets is called a pattern and can consists of almost any UTF8 encoded character, escape character, special escape and range. In general the character class token can be seen as a shortcut for a string alternation.

For example the following definition lets the `START` token hold either the strings "a", "b" or "c".

```tavor
START = "a" | "b" | "c"
```

With a character class this can be written as the following.

```tavor
START = [abc]
```

### <a name="character-classes-escapes"></a>Escape characters

The following table holds UTF8 encoded characters which are not directly allowed within a character class pattern. Their equivalent escape sequence has to be used instead.

| Character       | Escape sequence   |
| :-------------- | :---------------- |
| `-`             | `\-`              |
| `\`             | `\\`              |
| form feed       | `\f`              |
| newline         | `\n`              |
| return          | `\r`              |
| tab             | `\t`              |

For example the following defines that the `START` token can hold only white space characters.

```tavor
START = +([ \n\t\n\r])
```

Since some characters can be hard to type and read the `\x` escape sequence can be used to define them with their hexadecimal code points. This is also needed to explicitly define specific character independent of the text encoding. There are two options to do this. Either only two hexadecimal characters are used in the form of `\x0A` or when more then two hexadecimal digits are needed the form `\x{0AF}` has to be used. The second form allows up to 8 digits and is therefore fully Unicode ready.

To give an example the following definition holds either the Unicode character "/" or "😃".

```tavor
START = [\x2F\x{1F603}]
```

### <a name="character-classes-ranges"></a>Ranges

Ranges can be defined using the `-` character. A range holds all characters starting at the character before the `-` and ending at the character after the `-`. Both characters have to be either an UTF8 encoded or an escaped character. The starting character must have a lower value than the ending character.

For example the following defines a decimal digit.

```tavor
START = [0123456789]
```

This can be easier defined using a range.

```tavor
START = [0-9]
```

It is also possible to use hexadecimal code points, since either range characters can be escape characters.

```tavor
START = [\x23-\x5B]
```

### <a name="character-classes-special-escapes"></a>Special escape characters

Special escape characters combine many characters into one escape character and can also hold additional functionality. The following table is an overview of all currently implemented special escape characters.

| Special escape character | Character class           | Description                     |
| :----------------------- | :------------------------ | :------------------------------ |
| `\d`                     | `[0-9]`                   | Holds a decimal digit character |
| `\s`                     | `[ \f\n\r\t]`             | Holds a white space character   |
| `\w`                     | `[a-zA-Z0-9_]`            | Holds a word character          |

## <a name="attributes"></a>Token attributes

Some tokens define attributes which can be used in a definition by prepending a dollar sign to their name and appending a dot followed by the attribute name.

All list tokens have for example the `Count` attribute which holds the count of the token's direct child entries.

```tavor
Number = +([0-9])
START = "The number " Number " has " $Number.Count " digits"
```

When fuzzed this example will generate for example the string "The number 56 has 2 digits".

Some attributes can have arguments. An argument list begins with the opening parenthesis `(` and ends with the closing parenthesis `)`. Each argument is an [expression](#expressions) without the expression frame `${...}`. Attributes are separated by a comma.

All list tokens have for example the `Item` attribute which holds a child entry of the token. `Item` has one argument which is the index to the child entry.

```tavor
Letters = "a" "b" "c"
START = "The letter with the index 1 is " $Letters.Item(1)
```

When fuzzed this example will generate the string "The letter with the index 1 is b".

### <a name="attributes-general"></a>General attributes

The following enumeration defines and describes currently implemented general token attributes.

**List token**

A list token is a token which has in its definition either only a sequence of tokens or exactly one repeat group token.

| Attribute           | Arguments | Description                                                                                                         |
| :------------------ | :-------- | :------------------------------------------------------------------------------------------------------------------ |
| `Count`             | \-        | Holds the count of the token's direct child entries.                                                                |
| `Item`              | `i`       | Holds a child entry of the token with the index `i`.                                                                |
| `Unique`            | \-        | Chooses at random a direct child of the token and embeds it. The choice is unique for every reference of the token. |

### <a name="attributes-scope"></a>Scope of attributes

The Tavor format allows the usage of token attributes as long as the referenced token exists in the current scope.

Two main types of scopes exists:

- **Global scope** which is the scope of the whole format definition. An entry of the global scope is set by the nearest token reference to the `START` token.
- **Local scope** which is the scope held by a definition, group or any other token which opens up a new scope. Local scopes are initialized with entries from their parent scope at the time of the creation of the new local scope.

To give an example the following format definition is used.

```tavor
List = +,10("a")

Inner = "\tInner.1.Print: " $List.Count "\n",
        "\tInner.1.List: " List "\n",
        "\tInner.2.Print: " $List.Count "\n",
        "\tInner.3.Print: " $List.Count "\n",
        "\tInner.2.List: " List "\n",
        "\tInner.4.Print: " $List.Count "\n"

START = "Outer.1.Print: " $List.Count "\n",
        "Outer.1.List: " List "\n",
        "Outer.2.Print: " $List.Count "\n",
        Inner,
        "Outer.3.Print: " $List.Count "\n",
        "Outer.2.List: " List "\n",
        "Outer.4.Print: " $List.Count "\n",
        Inner,
        "Outer.5.Print: " $List.Count "\n",
        "Outer.3.List: " List "\n",
        "Outer.6.Print: " $List.Count "\n"
```

The format can result in the following fuzzing generation.

```
Outer.1.Print: 1
Outer.1.List: a
Outer.2.Print: 1
    Inner.1.Print: 1
    Inner.1.List: aa
    Inner.2.Print: 2
    Inner.3.Print: 2
    Inner.2.List: aaa
    Inner.4.Print: 3
Outer.3.Print: 1
Outer.2.List: aaaa
Outer.4.Print: 4
    Inner.1.Print: 4
    Inner.1.List: aaaaa
    Inner.2.Print: 5
    Inner.3.Print: 5
    Inner.2.List: aaaaaa
    Inner.4.Print: 6
Outer.5.Print: 4
Outer.3.List: aaaaaaa
Outer.6.Print: 7
```

This example generation shows that the first `$List.Count` usage attributed as `Outer.1.Print` uses the list `Outer.1.List` because it is the first usage of the token `List` next to the `START` token.

Additional observations can be made:

- Every new `List` reference overwrites the current entry of the current scope (e.g. `Outer.2.Print` uses `Outer.1.List`, the first `Inner.2.Print` uses the first `Inner.1.List`)
- An inner scope inherits from its parent scope (e.g. first `Inner.1.Print` uses `Outer.1.List`, second `Inner.1.Print` uses `Outer.2.List`)
- Parent scopes are not overwritten by their child scopes (e.g. `Outer.3.Print` uses `Outer.1.List`, `Outer.5.Print` uses `Outer.2.List`)

## <a name="typed-tokens"></a>Typed tokens

Typed tokens are an functional addition to regular token definitions of the Tavor format. They provide specific functionality which can be utilized by embedding them like regular tokens or through their additional token attributes. Typed tokens can be defined by prepending a dollar sign to their name. They do not have a format definition on their right-hand side. Instead, a type and optional arguments written as key-value pairs, which are separated by a colon, define the token.

A simple example for a typed token is the definition of an integer token.

```tavor
$Number Int

Addition = Number " + " (Number | Addition)

START = Addition
```

This format definition generates additions with random integers as numbers like for example `47245 + 6160 + 6137`. Note that since arguments of typed tokens are optional, the right hand side is optional.

The number of the `Int` type can be bounded in its range using arguments for the definition.

```tavor
$Number Int = from: 1,
              to:   10

Addition = Number " + " (Number | Addition)

START = Addition
```

Which generates for example `10 + 5 + 8 + 9`.

The following sections describe the currently implemented typed tokens with their arguments and attributes.

### <a name="typed-tokens-Int"></a>Type `Int`

The `Int` type implements a random integer.

#### Optional arguments

| Argument   | Description                                  |
| :--------- | :------------------------------------------- |
| `from`     | First integer value (defaults to 0)        |
| `to`       | Last integer value (defaults to 2<sup>31</sup> - 1) |

#### Token attributes

| Attribute | Arguments | Description                            |
| :-------- | :-------- | :------------------------------------- |
| `Value`   | \-        | Embeds a new token based on its parent |

### <a name="typed-tokens-Sequence"></a>Type `Sequence`

The `Sequence` type implements a generator for integers.

#### Optional arguments

| Argument   | Description                               |
| :--------- | :---------------------------------------- |
| `start`    | First sequence value (defaults to 1)      |
| `step`     | Increment of the sequence (defaults to 1) |

#### Token attributes

| Attribute  | Arguments | Description                                                 |
| :--------- | :-------- | :---------------------------------------------------------- |
| `Existing` | \-        | Embeds a new token holding one existing value of the parent |
| `Next`     | \-        | Embeds a new token holding the next value of the parent     |
| `Reset`    | \-        | Embeds a new token which on execution resets the parent     |

#### Example usages

The following example defines a sequence called `Id` which generates integers starting from 0 incremented by 2. It will therefore generate the sequence 0, 2, 4, 6 and so on. The example starts of by generating the first three values of the sequence using the token attribute `Next`. Afterwards the sequence is reseted using the token attribute `Reset` and then again three values are generated. Since the sequence got reseted before that the same values as in the beginning are generated. Ending the definition are three usages of the `Existing` token attribute which chooses at random one value out of all currently in use values of the sequence. Meaning it is possible that `Existing` chooses the same number more than once.

```tavor
$Id Sequence = start: 0,
               step:  2

START = +3("First Next: " $Id.Next "\n"),
        $Id.Reset,
        +3("Second Next: " $Id.Next "\n"),
        +3("Existing: " $Id.Existing "\n")
```

Will generate for example:

```
First Next: 0
First Next: 2
First Next: 4
Second Next: 0
Second Next: 2
Second Next: 4
Existing: 4
Existing: 2
Existing: 4
```

## <a name="expressions"></a>Expressions

Expressions can be used in token definitions and allow dynamic and complex operations using operators who can have different numbers of operands. An expressions starts with the dollar sign `$` and the opening curly brace `{` and ends with the closing curly brace `}`.

A simple example for an expression is an addition.

```tavor
START = ${1 + 2}
```

Every operand of an operator can be a token. The usual dollar sign for a token attribute can be omitted.

```tavor
$Number Int

START = ${Number.Value + 1}
```

The following sections describe the currently implemented operators.

### <a name="expressions-arithmetic"></a>Arithmetic operators

Arithmetic operators have two operands between the operator sign. Note that operators currently always embed the right side. This means that `2 * 3 + 4` will result into `2 * (3 + 4)` and not `(2 * 3) + 4`.

#### Operators

| Operator | Description    |
| :------- | :------------- |
| `+`      | Addition       |
| `-`      | Subtraction    |
| `*`      | Multiplication |
| `/`      | Division       |

#### Example usages

```tavor
START = ${9 + 8 + 7} "\n",
        ${6 - 5} "\n",
        ${4 * 3} "\n",
        ${10 / 2} "\n"
```

### <a name="expressions-graph"></a>Graph operators (experimental)

Graph operators are currently experimental since only a specific case has been implemented. The `path` operator traverses a list token based on the described structure. The structure defines the starting value of the traversal, the value which identifies each entry of the list token, how the entries are connected and which values are ending values for the traversal.

The `path` operator has the following format:

`path from (<starting value>) over (<entry identifier>) connected by (<entry connections>) without(<ending values>)`

All values are expressions. Furthermore, the `entry connections` and `ending values` are expressions lists. The `entry identifier`, `entry connections` and `ending values` have the variable `e` in their scope which holds the currently traversed entry of the token list.

> **Note**: Since the `path` operator acts on a list token, it might be necessary to use a variable reference, to avoid loops in the token definition.

The following example defines a list of connections called `Pairs`. Each entry in the list `Pairs` defines the identifier as its first token and the connection as its second token. The used `path` operator arguments define that all entries are traversed beginning from the value `2` and ending at the value `0`.

```tavor
START = Pairs "->" Path

Path = ${Pairs path from (2) over (e.Item(0)) connect by (e.Item(1)) without (0)}

Pairs = (,
	(1 0),
	(3 1),
	(2 3),
)
```
This example generates:

```
103123->231
```

> **Note**: The `path` operator can also traverse trees and graphs which have loops, and it can be combined with set operators.

### <a name="expressions-set"></a>Set operators (experimental)

Set operators are currently experimental since only a specific case has been implemented. The `not in` operator queries the `Existing` token attribute of a sequence to not include the given expression list. An expression list begins with the opening parenthesis `(` and ends with the closing parenthesis `)`. Each [expression](#expressions) is defined without the expression frame `${...}`. Expressions are separated by a comma.

```tavor
$Id Sequence

Pair = $Id.Next<id> " " ${Id.Existing not in (id)} "\n"

START = $Id.Reset +2(Pair)
```

This example generates:

```
1 2
2 1
```

The `Existing` token attribute can choose only between the values `1` and `2`, since the sequence generates only two values in this format definition. The `not in` operator excludes the given expression which is the variable `id` that holds the current sequence value. Hence if the current value is `1` only `2` can be used by `Existing` and if the value is `2` only `1` can be used.

## <a name="variables"></a>Variables

Every token of a token definition can be saved into a variable which consists of a name and a reference to a token usage. Variables follow the [same scope rules](#attributes-scope) as token attributes. It is therefore possible to for example define the same variable name more than once in one token sequence. They also do not overwrite variables definitions of parent scopes. Variables can be defined by using the `<` character after the token which should be saved, then defining the name of the variable and closing with the `>` character. They have a range of token attributes such as `Value`, which embeds a new token based on the current state of the referenced token.

In the following example the string token "text" will be saved into the variable `var`. The `Print` token uses this variable by embedding the referenced token.

```tavor
START = "text"<var> "->" Print

Print = $var.Value
```

This generates the string "text->text"

### <a name="variables-token-attributes"></a> Token attributes

Variables have the following token attributes:

| Attribute   | Arguments | Description                                                       |
| :---------- | :-------- | :---------------------------------------------------------------- |
| `Count`     | \-        | Holds the count of the referenced token's direct child entries    |
| `Index`     | \-        | Holds the index of the referenced token in relation to its parent |
| `Item`      | `i`       | Holds a child entry of the referenced token with the index `i`    |
| `Reference` | \-        | Holds a reference to a token which is needed for some operators   |
| `Value`     | \-        | Embeds a new token based on the referenced token                  |

### <a name="variables-just-save"></a>Just-save operator

Tokens which are saved to a variables are by default relayed to the generation. This means that their usage generates data as usual. Since this is sometimes unwanted, the just-save operator can be used to omit the relay. This is accomplished by adding the equal sign `=` after the `<` character.

```tavor
$Number Int

START = Number<=a> Number<=b>,
        a " + " b " = " ${a.Value + b.Value} "\n",
        a " * " b " = " ${a.Value * b.Value} "\n"
```

This format definition will generate for example:

```
5 + 3 = 8
5 * 3 = 15
```

The two usages of the `Number` token are hence only saved as variables and not relayed to the generation.

## <a name="statements"></a>Statements

Statements allow the Tavor format to have a control flow in its token definitions depending on the generated tokens and values. All statements start with the opening curly brace `{` and end with closing curly brace `}`. Right after `{`, the statement operator must be defined.

### <a name="statements-if"></a>`if` statement

The `if` statement allows to embed conditions into token definitions and defines an if body which is a scope on its own. The body lies between an opening `{if condition}` statement and an ending `{endif}` statement. The condition can be formed using the if statement's operators.

The following example will generate the character "A" if the variable `var` is equal to `1`.

```tavor
Choose = 1 | 2 | 3

START = Choose<var> "->" {if var.Value == 1}"A"{endif}
```

Additional to the `if` statement, the statements `else` and `else if` can be used. They can only be defined inside an if body and always belong to the `if` statement which the body belongs to. Both statement operators create a new scope.

The following example will generate the character "A" if the variable `var` is equal to `1`, "B" if its equal to `2` and "C" if its some other value.

```tavor
Choose = 1 | 2 | 3

START = Choose<var> "->" {if var.Value == 1}"A"{else if var.Value == 2}"B"{else}"C"{endif}
```


#### <a name="statements-if-operators"></a>Operators

Operands can be (if not otherwise described) defined tokens of all kind, variables or terminal tokens.

| Operator  | Usage        | Description                              |
| :-------- | :----------- | :--------------------------------------- |
| `==`      | `op1 == op2` | Returns true if op1 is equal to op2      |
| `defined` | `defined op` | Returns true if op is a defined variable |
