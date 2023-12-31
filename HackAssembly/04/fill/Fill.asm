// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

 
(LOOP)
//reset pixel to clear
@pixels
M=0

//check kbd
@KBD
D=M
//if 0 skip setting -1
@START
D;JEQ

//set -1
@pixels
M=-1
//reset keyboard
//@KBD
//M=0
(START)


@i
M=0

(DRAW)

//dest = SCREEN + i
@SCREEN
D=A
@i
D=D+M
@dest
M=D

//dest = screen +i ; *dest = pixels
@pixels
D=M
@dest
A=M
M=D

@i
MD=M+1
//i + 1 < 8192 then DRAW
@8192
D=D-A
@DRAW
D;JLT

@LOOP
0, JMP


(END)
@END
0;JMP // Infinite loop
