use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Machine {
    pc: usize,
    mem: Vec<i32>,
    inp: i32,
}

enum Mode {
    Pos,
    Imm,
}

impl Mode {
    fn i_to_a(mode: i32) -> Result<Self, i32> {
        match mode {
            0 => Ok(Self::Pos),
            1 => Ok(Self::Imm),
            _ => Err(mode),
        }
    }
}

impl Machine {
    fn new(mem: Vec<i32>, inp: i32) -> Self {
        Self {
            pc: 0,
            mem: mem,
            inp: inp,
        }
    }

    fn decode_op(code: i32) -> Result<(i32, Mode, Mode, Mode), i32> {
        let op = code % 100;
        let m1 = match Mode::i_to_a(code / 100 % 10) {
            Ok(m) => m,
            Err(e) => return Err(e),
        };
        let m2 = match Mode::i_to_a(code / 1000 % 10) {
            Ok(m) => m,
            Err(e) => return Err(e),
        };
        let m3 = match Mode::i_to_a(code / 10000 % 10) {
            Ok(m) => m,
            Err(e) => return Err(e),
        };
        Ok((op, m1, m2, m3))
    }

    fn eval_arg(&self, mode: Mode, arg: i32) -> i32 {
        match mode {
            Mode::Pos => self.mem[arg as usize],
            Mode::Imm => arg,
        }
    }

    fn get_arg(&self, offset: usize) -> i32 {
        self.mem[self.pc + offset]
    }

    fn step_pc(&mut self, offset: usize) {
        self.pc += offset;
    }

    fn exec(&mut self) -> Result<bool, (bool, i32)> {
        let (op, a1, a2, _) = match Self::decode_op(self.mem[self.pc]) {
            Ok(k) => k,
            Err(e) => return Err((true, e)),
        };

        match op {
            1 => {
                let arg1 = self.get_arg(1);
                let arg2 = self.get_arg(2);
                let dest = self.get_arg(3) as usize;
                self.mem[dest] = self.eval_arg(a1, arg1) + self.eval_arg(a2, arg2);
                self.step_pc(4);
                Ok(true)
            }
            2 => {
                let arg1 = self.get_arg(1);
                let arg2 = self.get_arg(2);
                let dest = self.get_arg(3) as usize;
                self.mem[dest] = self.eval_arg(a1, arg1) * self.eval_arg(a2, arg2);
                self.step_pc(4);
                Ok(true)
            }
            3 => {
                let arg1 = self.get_arg(1) as usize;
                self.mem[arg1] = self.inp;
                self.step_pc(2);
                Ok(true)
            }
            4 => {
                let arg1 = self.get_arg(1);
                println!("{}", self.eval_arg(a1, arg1));
                self.step_pc(2);
                Ok(true)
            }
            5 => {
                let arg1 = self.get_arg(1);
                let arg2 = self.get_arg(2);
                if self.eval_arg(a1, arg1) != 0 {
                    self.pc = self.eval_arg(a2, arg2) as usize;
                } else {
                    self.step_pc(3);
                }
                Ok(true)
            }
            6 => {
                let arg1 = self.get_arg(1);
                let arg2 = self.get_arg(2);
                if self.eval_arg(a1, arg1) == 0 {
                    self.pc = self.eval_arg(a2, arg2) as usize;
                } else {
                    self.step_pc(3);
                }
                Ok(true)
            }
            7 => {
                let arg1 = self.get_arg(1);
                let arg2 = self.get_arg(2);
                let arg3 = self.get_arg(3) as usize;
                if self.eval_arg(a1, arg1) < self.eval_arg(a2, arg2) {
                    self.mem[arg3] = 1;
                } else {
                    self.mem[arg3] = 0;
                }
                self.step_pc(4);
                Ok(true)
            }
            8 => {
                let arg1 = self.get_arg(1);
                let arg2 = self.get_arg(2);
                let arg3 = self.get_arg(3) as usize;
                if self.eval_arg(a1, arg1) == self.eval_arg(a2, arg2) {
                    self.mem[arg3] = 1;
                } else {
                    self.mem[arg3] = 0;
                }
                self.step_pc(4);
                Ok(true)
            }
            99 => {
                self.step_pc(1);
                Ok(false)
            }
            _ => Err((false, self.mem[self.pc])),
        }
    }

    fn execute(&mut self) {
        loop {
            match self.exec() {
                Ok(ok) => {
                    if !ok {
                        return;
                    }
                }
                Err((mode, k)) => {
                    if mode {
                        panic!("invalid arg mode: {}", k);
                    } else {
                        panic!("invalid op code: {}", k);
                    }
                }
            }
        }
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        for line in reader.lines() {
            for i in line.expect("Failed to read line").split(",") {
                nums.push(i.parse::<i32>().expect("Failed to parse num"));
            }
        }
        nums
    };

    {
        let mut m = Machine::new(nums.clone(), 1);
        m.execute();
    }
    {
        let mut m = Machine::new(nums.clone(), 5);
        m.execute();
    }
}
