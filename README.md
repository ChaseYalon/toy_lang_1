
# Welcome to toy lang

## The basic features are outlined bellow

## Feel free to reach out to <chaseyalon@gmail.com> with questions

### Build instructions
Either download a precompiled binary at the git repo under releases or
1. Clone git repo
2. Download go programming language
3. Cd into the directory
4. run ```go build -o "toy_lang"``` on linux/MacOs or ```go build -o "toy_lang.exe"``` on windows
5. Then you can run that binary raw to get a REPL or pass a file a .toy file and run it

### Documentation

- Declare a variable with let
- DONT YOU DARE FORGET A SEMICOLON, you will get a really confusing error message
- 3 Supported datatypes, int, bool, and string

```toy
let a = 2;
let b = "hello " + "world";
let c = true || false;
```

- Arithmetic is supported fully
    - Plus (+), also used for string concatenation
    - Minus (-)
    - Multiply (*)
    - Divide (/)
    - Modulo (%v)
    - Exponent (**)
    - And (&&)
    - Or (||)
    - Not (!)
    - Less than (<)
    - Less than or equal to (<=)
    - Greater than (>)
    - Greater than or equal to (>=)
    - Equal to (==)
    - Not equal to (!=)
- Toy Lang also supports the following compound expressions
    - Plus equals (+=)
    - Minus equals (-=)
    - Multiply equals (*=)
    - Divide equals (/=)
    - Plus plus (++)
    - Minus minus (--)
- You can use them in inline expressions

```toy
let x = 2;
let y = 2 < 3;
```

- There are 6 builtin functions
    - print(str) prints a value to the screen
    - println(str) prints a value and a newline to the screen
    - input(str) prints a prompt to the screen and returns the user input
    - str(bool | int) converts a bool or int to a string
    - bool(str | int) converts a string or an int to a bool
    - int(str | bool) converts a string or bool to an int

Get a user input, add 2 and print it like this
```toy
let uIn = input("Enter a number: ");
let uInInt = int(uIn);
println(uInINt + 2);
```

- Toy Lang supports if statements as follows
```toy
if |COND|{
    |BODY|
}
```
- You can do if else like this
```toy
if |COND|{
    |BODY|
} else {
    |ALT|
}
```
- To do elsif just nest the second (and any subsequent if's) inside the else or use the guard clause technique
- In toy lang, functions are second class citizens (will change later) and can be declared like this 
```toy
fn add(a, b){
    return a + b;
}
fn sayHello(){
    return "hello";
}
```
- You call them like this 
```toy
fn add(a, b){
    return a + b;
}
let five = add(2, 3); 
```

- Toy lang supports while loops with the continue and break key words like this
```toy

let x = 0;
while x < 10{
    println(x);
    x++;
}

let y = 100;
while y > 0{
    println(y);
    if y == 40{
        continue;
    }
    if y == 30{
        break;
    }
}

```