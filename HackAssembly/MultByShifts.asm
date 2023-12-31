// multiply(x, y):
// // Where x; yb 0
// sum ¼ 0
// shiftedX ¼ x
// for j ¼ 0 ...ðn  1Þ do
// if (j-th bit of y) ¼ 1 then
// sum ¼ sum þ shiftedX
// shiftedX ¼ shiftedX  2
// Figure 12.1 Multiplication of two n-bit number

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

// HEre is logic and example, like school math

//  x           1011
//  y            101    jth bit of y, 
//  _________________
//              1011    1
//  +          0000     0
//  +         101100    1  -> as we see we added 101100 shifted, thats why we shift each time
//  ____________________
//  x*y       110111

//set r2 to 0
@R2
M=0
//shiftedX = R0 first argument
@R0
D=M
@shiftedX
M=D

//check jth bit from 0 to n. we have 16 bits max
//loop from 0 t0 15
@jth_bit_Mask
M=1
//we will add it itself each time to move jth bit 
// 1+1 -. 0b10  , ob10+ob10 -> 0b100 

@i
M=0
(LOOP_BITS)

//if jth bit of R1 (number y is set) we will add shifted Of X  to the sum
//if not we will just shift it
@jth_bit_Mask
D=M
@R1
//and to see if it is set
D=D&M
@SHIFT
D;JEQ

//we will add shiftedX into sum
//sum +=shiftedX ,R2 is our sum 
//
@shiftedX
D=M
@R2
M=D+M

///////////////////////////////////////////
(SHIFT)
//shifting X , it means multiply by2 or just adding itself
@shiftedX
D=M
M=D+M

//shifting jth_bit_Mask
@jth_bit_Mask
D=M
M=D+M

//LOOPING bits
@i
MD=M+1
@16
D=A-D
@LOOP_BITS
D;JGT


(END)
@END
0;JMP // Infinite loop




