'''
Assmebler written for the hack assembly language
@author Abdelrauf
'''


import sys
import os
CASE_SENSITIVE=True

class SymbolTable:
    table = []

    def __init__(self):
        #initialize
        self.addEntry("SP", 0)
        self.addEntry("LCL", 1)
        self.addEntry("ARG", 2)
        self.addEntry("THIS", 3)
        self.addEntry("THAT", 4)
        for h in range(16):
            self.addEntry("R"+str(h), h)
        self.addEntry("screen", 16384)
        self.addEntry("kbd", 24576)



    def get(self, k):
        kk = k if CASE_SENSITIVE else k.upper()
        if len(self.table) > 0:
            for idx,  (x,address) in enumerate(self.table):
                if x == k:
                    return (idx, address)
        return None

    def getAddress(self, k):
        ret = self.get(k)
        return ret[1] if ret is not None else None

    def addEntry(self, k, val):
        kk = k if CASE_SENSITIVE else k.upper()
        ret =  self.get(kk)
        if ret is None:
            self.table.append((kk,val))
        else:
            print(f"warning : modification {k}'s old value {ret[1]} -> {val}")
            self.table[ret[0]]=val


class Logger:

    def __init__(self, log =True):
        self.doLog=log
        pass

    def info(self,x):
        if self.doLog==True:
            print(x)

#global logger
log = Logger()

class Token:
    Text = 1
    NonText = 2
    EOL = 88
    EOF = 99
    spaceList = " \n\r\t"

    def __init__(self):
        pass


    def skipSpace(self, line , beg, end):
        idx = beg
        while idx<end and Token.isSpace(line[idx]):
            idx+=1
        return idx

    def isAlphaVariableNameChar(c):
        cc = ord(c)
        if cc>=ord('a') and cc <= ord('z'):
            return True
        if cc>=ord('A') and cc <= ord('Z'):
            return True
        if cc>=ord('0') and cc <= ord('9'):
            return True
        if cc  == ord("_")  or cc == ord(".") or cc==ord("$")  or cc==ord("%"):
            return True
        return False
    
    def isSpace(c):
        return c in Token.spaceList

    def isSymbolCandidate(c):
        if Token.isSpace(c):
            return False
        elif Token.isAlphaVariableNameChar(c):
            return False
        return True
    
    def isMathSymbolItself(c):
        if c in "+-!&|=":
            return True
        else:
            return False

    def checkAlphaNumeric(self, line, beg, end):
        idx = self.skipSpace(line, beg, end )
        beg = idx
        while idx<end and Token.isAlphaVariableNameChar(line[idx]) :
            idx+=1
        return (beg, idx)



    def checkSymbol(self, line, beg, end):
        idx = self.skipSpace(line, beg, end )
        beg = idx
        while idx<end and Token.isSymbolCandidate(line[idx]):
            c = line[idx]
            idx+=1
            if Token.isMathSymbolItself(c):
                break
        return (beg, idx)

    def tokenize(self, filename):
        _line_pos = 1;
        with open(filename , "r") as f:
            for line in f:
                end = len(line)
                pos = _line_pos
                _line_pos+=1
                idx = 0 
                while idx<end:
                    beg,idx = self.checkAlphaNumeric(line, idx, end)

                    if idx>beg:
                        yield (pos, beg,  idx ,Token.Text, line[beg:idx].upper())
                    beg,idx = self.checkSymbol(line, idx, end)
                    if idx>beg:
                        #ignore comment
                        if line[beg:idx] == "//":
                            log.info("TOKENIZER:: comment " + line[beg:end].strip())
                            break
                        yield (pos, beg,  idx , Token.NonText, line[beg:idx])
                yield(_line_pos, end-1, end-1, Token.EOL, "")
            yield(_line_pos, 0, 0, Token.EOF, "")


class Parser:
    '''
     Parsing 
    '''


    def __init__(self, filename ):
        self.filename =filename
        pass


    def parse_comp(self, comp, next_token, tokenizer):
        #we expect comp
        # next_token = tokenizer.__next__()
        l, beg, end, type, val  = next_token
        if type==Token.EOL or type==Token.EOF:
            return next_token
        elif type != Token.Text:
            if not val in "!+&|-":
                return next_token
            comp[1] = val
            #expect last operand
            next_token = tokenizer.__next__()
            l, beg, end, type, val  = next_token
            if type != Token.Text:
                raise Exception(f"syntax error in the line {l}: col: {beg}    {val}")
            comp[2] = val
            next_token = tokenizer.__next__()
        else:
            comp[0] = val
            next_token = tokenizer.__next__()
            next_token = self.parse_comp(comp,next_token, tokenizer)
        return next_token


    def parse(self):
        #first pass. preprocess labels
        last_line = 0
        parse_type = CodeGenerator.Undef
        tokenizer = Token().tokenize(self.filename)
        next_token = tokenizer.__next__()
        infinite_loop = 0
        while True :
            l, beg, end, type, val  = next_token
            if type == Token.EOF:
                break
            if type == Token.EOL:
                parse_type = CodeGenerator.Undef
                next_token =  tokenizer.__next__()
                continue
            if parse_type == CodeGenerator.A:
                #this should be text in the same line
                if l!=last_line or  type != Token.Text:
                    raise Exception(f"PARSER:: syntax error in the line {l}: col: {beg}    {val}")
                yield ( CodeGenerator.A, val)
                parse_type = CodeGenerator.Undef
                next_token = tokenizer.__next__()
            elif parse_type == CodeGenerator.L:
                # it should follow by text and end with ")" in the sme line
                if l!=last_line or  type != Token.Text:
                    raise Exception(f"PARSER:: syntax error in the line {l}: col: {beg}    {val}")
                text = val
                next_token = tokenizer.__next__()
                l, beg, end, type, val  = next_token
                if l!=last_line or  type != Token.NonText or val!=")":
                    raise Exception(f"PARSER:: syntax error in the line {l}: col: {beg}    {val}")
                yield( CodeGenerator.L, text)
                parse_type = CodeGenerator.Undef
                next_token = tokenizer.__next__()
            else:
                #c intruction
                if type==Token.NonText and  val.startswith("@"):
                        parse_type = CodeGenerator.A
                        last_line = l
                        next_token = tokenizer.__next__()
                elif  type==Token.NonText and  val=="(":
                        parse_type = CodeGenerator.L
                        last_line = l
                        next_token = tokenizer.__next__()

                else:
                    parse_type = CodeGenerator.C
                    last_line = l

                    #parse C type
                    #dest=comp; jmp
                    #(dest, comp) or (comp,jmp) or (dest,comp,jmp)
                    dest = None
                    # operand1 op operand2
                    comp = [None, None, None]
                    jmp = None
                    if type == Token.Text:
                        t = val
                        next_token = tokenizer.__next__()
                        l, beg, end, type, val  = next_token
                        if type == Token.NonText and val=="=":
                            dest = t
                            next_token = tokenizer.__next__()
                        else:
                            comp[0] =t
                    #parse
                    next_token = self.parse_comp( comp, next_token, tokenizer)

                    l, beg, end, type, val  = next_token
                    if type==Token.NonText and val == ";":
                        #expect jmp
                        next_token = tokenizer.__next__()
                        l, beg, end, type, val  = next_token
                        if type == Token.Text:
                            if val.startswith("J"):
                                jmp = val
                                next_token = tokenizer.__next__()
                            else:
                                raise Exception(f"PARSER:: syntax error in the line {l}: col: {beg}    {val}")
                    if dest==None and comp==[None,None,None] and jmp==None:
                        raise Exception(f"PARSER:: syntax error in the line {l}: col: {beg}    {val}")
                    #just yield
                    yield (CodeGenerator.C, (dest, comp, jmp))
                    parse_type = CodeGenerator.Undef



class CodeGenerator:
    A = 1
    C = 2
    L = 3
    Undef = 99

    def __init__(self, parser):
        self.table = SymbolTable()
        self.parser = parser
        self.var_addr = 16

    def generateA(self, val):
        addr = 0
        if not val.isdigit():
            ret = self.table.getAddress(val)
            if ret is None:
                self.table.addEntry(val, self.var_addr)
                addr = self.var_addr
                self.var_addr +=1
            else:
                addr = int(ret)
        else:
            addr = int(val)

        return '0{:015b}'.format(addr).strip()
    
    def generateC(self, dest, comp, jmp):
        ret = list('1110000000000000')
        #111accccccdddjjj
        #setting d bits
        if dest is not None:
            for c in dest:
                if c=='A':
                    ret[10]='1'
                elif c=='D':
                    ret[11]='1'
                elif c=='M':
                    ret[12] = '1'
        #setting j bits
        if jmp is not None:
            m=[None, 'JGT', 'JEQ', 'JGE', 'JLT', 'JNE', 'JLE', 'JMP']
            for idx, x in enumerate(m):
                if x==jmp:
                    if idx & 4 == 4:
                        ret[13] ='1'
                    if idx & 2 ==2:
                        ret[14] ='1'
                    if idx & 1 == 1:
                        ret[15] = '1'

                    break
        if comp is not None:
            a,op, b = comp
            if a=="M" or b=="M":
                #make a 1
                ret[3] = "1"
            #Additional cases: 
            # ignore/correct mathematically some cases where the one of operand is 0
            if b=="0":
                op=None
                b=None
            if a=="0":
                if op=="+" or op =="|":
                    op=None
                    a=b
                    b=None
                elif op=="-":
                    a=None
                elif op == "&":
                    op=None
                    b=None

            if op == "&":
                pass
            elif op=="|":
                #010101
                #it means !(!a & !b) which is() (a|b)')'
                ret[5]="1"
                ret[7]="1"
                ret[9]="1"
            elif op=="+":
                ret[8]="1"
                if a=="1" or b=="1":
                    #we should set 1 all besides below
                    for i in range(4,10):
                        ret[i] = '1'
                    if a=="D" or b=="D":
                        #011111
                        #~( ~d + ~0) . 
                        #reset first to zero
                        ret[4]='0'
                    else:
                        ret[6] ='0'
            elif op == "-":
                ret[8] ='1'
                if b=="1":
                    if a==None:
                        ret[4]="1"
                        ret[5]="1"
                        ret[6]="1"
                    elif a=="D":
                        ret[6]="1"
                        ret[7]="1"
                    elif a=="A" or a=="M":
                        ret[4]="1"
                        ret[5]="1"
                elif a==None:
                    ret[9]='1'
                    if b=="D":
                        ret[6]="1"
                        ret[7]="1"
                    else:
                        ret[4]="1"
                        ret[5]="1"
                else:
                    ret[8]="1"
                    ret[9]="1"
                    if a=="D":
                        ret[5]="1"
                    elif b=="D":
                        ret[7]="1"
            elif op==None:
                if a=="0":
                    #101010
                    ret[4]="1"
                    ret[6]="1"
                    ret[8]="1"
                elif a=="1":
                    for i in range(4,10):
                        ret[i]="1"
                elif a=="D":
                    ret[6]="1"
                    ret[7]="1"
                elif a=="A" or a=="M":
                    ret[4]="1"
                    ret[5]="1"
            elif op=="!":
                ret[9]="1"
                if b=="D":
                    ret[6]="1"
                    ret[7]="1"
                else:
                    ret[4]="1"
                    ret[5]="1"
                    
        return "".join(ret).strip()

    def preProcessLcommands(self):
        nn = 0
        log.info("------------preProcessLcommands----------------")
        for h in self.parser.parse():
            log.info("PARSER:: " + str(h))
            if h[0] == CodeGenerator.L:
                self.table.addEntry( h[1], nn)
            else:
                nn+=1
        log.info("---------------------------------")

    def code_generate(self):
        #first pass: build labels
        self.preProcessLcommands()
        #second pass:
        #generate codes
        #for variables 16+
        log.info("------------generate Machine Code----------------")
        for h in self.parser.parse():
            log.info("PARSER:: " + str(h))
            if h[0] == CodeGenerator.A:
                tt = self.generateA(h[1])
                log.info("CODE GEN:: "+tt[0] + " -- " + tt[1:])
                yield tt
            elif h[0] == CodeGenerator.C:

                tt  = self.generateC(*h[1])
                log.info("CODE GEN:: "+tt[:3] + " -- "+tt[3]+" -- "+ tt[4:10] + " -- " + tt[10:13] + " -- " + tt[13:])
                yield tt



def start():
    filename = "" 
    if len(sys.argv)>1:
        filename = sys.argv[1]
    base = os.path.splitext(filename)[0]
    with open(base+".hack", "w") as f:
        try:
            c = CodeGenerator(Parser(filename)) 
            for h in c.code_generate():
                print(h,  file=f)
        except Exception as e:
            print(e)

if __name__ == "__main__":
    start()