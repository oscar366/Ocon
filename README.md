# OCON

## About
OCON (Oscar's Console Language) is my programming language, possibly the worst ever made. 🐟
It is made mostly by me using golang

OCON does not use an [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree), and I probably won't add one in the future because switching to one would be difficult at this point. When I started the project, I didn't even know ASTs existed, so OCON's interpreter developed around a different approach.


> **Platform:** Ocon currently is tested Windows only.  
> macOS and Linux may have bugs or other things.

## Why use OCON
I have no idea why you would use this insted of other languages.
If you find a reason, please tell me.

## Current (semi-working) Features

- Variables
- Arithmetic
- Boolean operations
- Conditionals
- Sections
- Goto command
- Math commands
- Comments
- Return commands
- Imports & fancy imports
- Flags (eg `!debugOFF` turns off console debug output)
- Updating 
- Functions

## Hopefuly new Features to come

- Arrays (currently working on)
- String handling
- Better loops 
- User input
- Error handling
- Files
- Reflection

## Egg samples of Ocon

### Hello world:
```ocon
echo Hello_world!
```

### increment:
```ocon
var "i '0

§ "loop

increment "i '1
echo $i

if [ isless $i '5 ] goto "loop else continue
```

### full test script:
See Test script.ocon

## How to run Ocon

```shell
ocon execute "filename.ocon"
```

## More
This README was made only by me no AI. Sorry if there were any spelling mistakes I can't spell good.

## Contact me

Email: [oscarfoxpatterson@gmail.com](mailto:oscarfoxpatterson@gmail.com)

Discord: **0thisisauser** ← use this one


## Icon
![OCON logo](ocon.png)