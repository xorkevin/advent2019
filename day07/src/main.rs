use crossbeam::channel;
use crossbeam::channel::{Receiver, Sender};
use crossbeam::thread;
use permutohedron;
use std::error::Error;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

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

struct Machine {
    pc: usize,
    mem: Vec<i32>,
    inp: Receiver<i32>,
    out: Sender<i32>,
    out_gauge: i32,
}

impl Machine {
    fn new(mem: Vec<i32>, inp: Receiver<i32>, out: Sender<i32>) -> Self {
        Self {
            pc: 0,
            mem,
            inp,
            out,
            out_gauge: 0,
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

    fn recv_inp(&self) -> Result<i32, channel::RecvError> {
        self.inp.recv()
    }

    fn send_out(&mut self, out: i32) -> Result<(), channel::SendError<i32>> {
        self.out_gauge = out;
        self.out.send(out)
    }

    fn exec(&mut self) -> Result<bool, Box<dyn Error>> {
        let (op, a1, a2, _) = match Self::decode_op(self.mem[self.pc]) {
            Ok(k) => k,
            Err(e) => return Err(format!("Invalid arg mode: {}", e).into()),
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
                self.mem[arg1] = match self.recv_inp() {
                    Ok(k) => k,
                    Err(e) => {
                        return Err(format!("Failed to read: {}", e).into());
                    }
                };
                self.step_pc(2);
                Ok(true)
            }
            4 => {
                let arg1 = self.get_arg(1);
                match self.send_out(self.eval_arg(a1, arg1)) {
                    Ok(_) => (),
                    Err(_) => {
                        return Err("Failed to send".into());
                    }
                };
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
            _ => Err(format!("Invalid op code: {}: {}", self.pc, self.mem[self.pc]).into()),
        }
    }

    fn execute(&mut self) -> Result<(), Box<dyn Error>> {
        loop {
            match self.exec() {
                Ok(ok) => {
                    if !ok {
                        return Ok(());
                    }
                }
                Err(k) => {
                    return Err(k);
                }
            }
        }
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let tokens = {
        let mut tokens = Vec::new();
        for line in reader.lines() {
            for i in line.expect("Failed to read line").split(",") {
                tokens.push(i.parse::<i32>().expect("Failed to parse num"));
            }
        }
        tokens
    };

    {
        let max =
            permutohedron::Heap::new(&mut (0..5).collect::<Vec<_>>()).try_fold(0, |out, phases| {
                let k = phases.iter().try_fold(0, |out, &phase| {
                    let (tx, rx) = channel::unbounded();
                    let (ntx, nrx) = channel::unbounded();
                    let mut m = Machine::new(tokens.clone(), rx, ntx);
                    match tx.send(phase) {
                        Ok(_) => (),
                        Err(e) => {
                            eprintln!("{}", e);
                            return None;
                        }
                    };
                    match tx.send(out) {
                        Ok(_) => (),
                        Err(e) => {
                            eprintln!("{}", e);
                            return None;
                        }
                    };
                    match m.execute() {
                        Ok(_) => (),
                        Err(e) => {
                            eprintln!("{}", e);
                            return None;
                        }
                    };
                    match nrx.recv() {
                        Ok(k) => Some(k),
                        Err(e) => {
                            eprintln!("{}", e);
                            None
                        }
                    }
                });
                match k {
                    Some(k) => {
                        if k > out {
                            Some(k)
                        } else {
                            Some(out)
                        }
                    }
                    None => None,
                }
            });
        let max = match max {
            Some(k) => k,
            None => return Err("Failed to run".into()),
        };
        println!("{}", max);
    }

    {
        let max = permutohedron::Heap::new(&mut (5..10).collect::<Vec<_>>()).try_fold(
            0,
            |out, phases| {
                let chans = phases
                    .iter()
                    .map(|_| channel::unbounded())
                    .collect::<Vec<_>>();

                let mut machines = Vec::new();
                for (i, &phase) in phases.iter().enumerate() {
                    let prev = (i + chans.len() - 1) % chans.len();
                    match chans[prev].0.send(phase) {
                        Ok(_) => (),
                        Err(e) => {
                            eprintln!("{}", e);
                            return None;
                        }
                    };
                    let rcv = chans[prev].1.clone();
                    let send = chans[i].0.clone();
                    machines.push(Machine::new(tokens.clone(), rcv, send));
                }
                match chans[chans.len() - 1].0.send(0) {
                    Ok(_) => (),
                    Err(e) => {
                        eprintln!("{}", e);
                        return None;
                    }
                };

                thread::scope(|s| {
                    for m in machines.iter_mut() {
                        s.spawn(move |_| match m.execute() {
                            Ok(_) => Ok(()),
                            Err(e) => {
                                eprintln!("{}", e);
                                Err(())
                            }
                        });
                    }
                })
                .unwrap();

                let k = machines[machines.len() - 1].out_gauge;
                if k > out {
                    Some(k)
                } else {
                    Some(out)
                }
            },
        );
        let max = match max {
            Some(k) => k,
            None => return Err("Failed to run".into()),
        };
        println!("{}", max);
    }

    Ok(())
}
