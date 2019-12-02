use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";
const PUZZLEINPUT2: usize = 19690720;

struct Machine {
    pc: usize,
    mem: Vec<usize>,
}

impl Machine {
    fn new(mem: Vec<usize>) -> Self {
        return Self { pc: 0, mem: mem };
    }

    fn exec(&mut self) -> Result<bool, usize> {
        match self.mem[self.pc] {
            1 => {
                let arg1 = self.mem[self.pc + 1];
                let arg2 = self.mem[self.pc + 2];
                let dest = self.mem[self.pc + 3];
                self.mem[dest] = self.mem[arg1] + self.mem[arg2];
                self.pc += 4;
                Ok(true)
            }
            2 => {
                let arg1 = self.mem[self.pc + 1];
                let arg2 = self.mem[self.pc + 2];
                let dest = self.mem[self.pc + 3];
                self.mem[dest] = self.mem[arg1] * self.mem[arg2];
                self.pc += 4;
                Ok(true)
            }
            99 => {
                self.pc += 1;
                Ok(false)
            }
            _ => Err(self.mem[self.pc]),
        }
    }

    fn execute(&mut self) -> Result<(), usize> {
        loop {
            match self.exec() {
                Ok(k) => {
                    if !k {
                        return Ok(());
                    }
                }
                Err(k) => return Err(k),
            }
        }
    }

    fn mem_at(&self, offset: usize) -> usize {
        self.mem[offset]
    }

    fn mem_set(&mut self, offset: usize, val: usize) {
        self.mem[offset] = val;
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        for line in reader.lines() {
            for i in line.expect("Failed to read line").split(",") {
                nums.push(i.parse::<usize>().expect("Failed to parse num"));
            }
        }
        nums
    };

    {
        let mut m = Machine::new(nums.clone());
        m.mem_set(1, 12);
        m.mem_set(2, 2);
        if let Err(k) = m.execute() {
            panic!("invalid op code: {}", k);
        }
        println!("{}", m.mem_at(0));
    }

    'outer: for i in 0..100 {
        for j in 0..100 {
            let mut m = Machine::new(nums.clone());
            m.mem_set(1, i);
            m.mem_set(2, j);
            if let Err(k) = m.execute() {
                panic!("invalid op code: {}", k);
            }
            if m.mem_at(0) == PUZZLEINPUT2 {
                println!("{}", i * 100 + j);
                break 'outer;
            }
        }
    }
}
