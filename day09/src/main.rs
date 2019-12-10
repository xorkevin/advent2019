use std::error::Error;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use tokio::sync::mpsc;

const PUZZLEINPUT: &str = "input.txt";
const RAM_SIZE: usize = 8192;

enum Mode {
    Pos,
    Imm,
    Rel,
}

impl Mode {
    fn i_to_a(mode: i64) -> Result<Self, i64> {
        match mode {
            0 => Ok(Self::Pos),
            1 => Ok(Self::Imm),
            2 => Ok(Self::Rel),
            _ => Err(mode),
        }
    }
}

struct Machine {
    id: u32,
    pc: i64,
    mem: Vec<i64>,
    inp: mpsc::Receiver<i64>,
    out: mpsc::Sender<i64>,
    out_gauge: i64,
    rel_base: i64,
}

impl Machine {
    fn new(id: u32, mem: Vec<i64>, inp: mpsc::Receiver<i64>, out: mpsc::Sender<i64>) -> Self {
        Self {
            id,
            pc: 0,
            mem,
            inp,
            out,
            out_gauge: 0,
            rel_base: 0,
        }
    }

    fn decode_op(code: i64) -> Result<(i64, Mode, Mode, Mode), i64> {
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

    fn get_mem(&self, pos: i64) -> i64 {
        self.mem[pos as usize]
    }

    fn eval_arg(&self, mode: Mode, arg: i64) -> i64 {
        match mode {
            Mode::Pos => self.get_mem(arg),
            Mode::Imm => arg,
            Mode::Rel => self.get_mem(arg + self.rel_base),
        }
    }

    fn get_arg(&self, mode: Mode, offset: i64) -> i64 {
        let arg = self.get_mem(self.pc + offset);
        self.eval_arg(mode, arg)
    }

    fn set_mem(&mut self, pos: i64, val: i64) {
        self.mem[pos as usize] = val;
    }

    fn set_arg(&mut self, mode: Mode, offset: i64, val: i64) -> Result<(), String> {
        let arg = self.get_mem(self.pc + offset);
        match mode {
            Mode::Pos => self.set_mem(arg, val),
            Mode::Imm => {
                return Err("Illegal mem write imm mode".into());
            }
            Mode::Rel => self.set_mem(arg + self.rel_base, val),
        }
        Ok(())
    }

    fn step_pc(&mut self, offset: i64) {
        self.pc += offset;
    }

    async fn recv_inp(&mut self) -> Option<i64> {
        self.inp.recv().await
    }

    async fn send_out(&mut self, out: i64) -> Result<(), mpsc::error::SendError<i64>> {
        self.out_gauge = out;
        self.out.send(out).await
    }

    async fn exec(&mut self) -> Result<bool, Box<dyn Error>> {
        let (op, a1, a2, a3) = match Self::decode_op(self.get_mem(self.pc)) {
            Ok(k) => k,
            Err(e) => return Err(format!("Invalid arg mode: {}", e).into()),
        };

        match op {
            1 => {
                self.set_arg(a3, 3, self.get_arg(a1, 1) + self.get_arg(a2, 2))?;
                self.step_pc(4);
            }
            2 => {
                self.set_arg(a3, 3, self.get_arg(a1, 1) * self.get_arg(a2, 2))?;
                self.step_pc(4);
            }
            3 => {
                let inp = match self.recv_inp().await {
                    Some(k) => k,
                    None => {
                        return Err("Failed to read".into());
                    }
                };
                self.set_arg(a1, 1, inp)?;
                self.step_pc(2);
            }
            4 => {
                match self.send_out(self.get_arg(a1, 1)).await {
                    Ok(_) => (),
                    Err(e) => {
                        return Err(format!("Failed to send: {}", e).into());
                    }
                };
                self.step_pc(2);
            }
            5 => {
                if self.get_arg(a1, 1) != 0 {
                    self.pc = self.get_arg(a2, 2);
                } else {
                    self.step_pc(3);
                }
            }
            6 => {
                if self.get_arg(a1, 1) == 0 {
                    self.pc = self.get_arg(a2, 2);
                } else {
                    self.step_pc(3);
                }
            }
            7 => {
                if self.get_arg(a1, 1) < self.get_arg(a2, 2) {
                    self.set_arg(a3, 3, 1)?;
                } else {
                    self.set_arg(a3, 3, 0)?;
                }
                self.step_pc(4);
            }
            8 => {
                if self.get_arg(a1, 1) == self.get_arg(a2, 2) {
                    self.set_arg(a3, 3, 1)?;
                } else {
                    self.set_arg(a3, 3, 0)?;
                }
                self.step_pc(4);
            }
            9 => {
                self.rel_base += self.get_arg(a1, 1);
                self.step_pc(2);
            }
            99 => {
                self.step_pc(1);
                return Ok(false);
            }
            _ => {
                return Err(
                    format!("Invalid op code: {}: {}", self.pc, self.get_mem(self.pc)).into(),
                );
            }
        }
        Ok(true)
    }

    async fn execute(&mut self) -> Result<(), Box<dyn Error>> {
        loop {
            match self.exec().await {
                Ok(ok) => {
                    if !ok {
                        return Ok(());
                    }
                }
                Err(k) => {
                    return Err(format!("Machine {}: {}", self.id, k).into());
                }
            }
        }
    }
}

fn copy_vec<T: Copy>(dest: &mut Vec<T>, src: &Vec<T>) {
    let k = std::cmp::min(dest.len(), src.len());
    for i in 0..k {
        dest[i] = src[i]
    }
}

#[tokio::main(threaded_scheduler)]
async fn main() -> Result<(), Box<dyn Error>> {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let tokens = {
        let mut tokens = Vec::new();
        for line in reader.lines() {
            for i in line.expect("Failed to read line").split(",") {
                tokens.push(i.parse::<i64>().expect("Failed to parse num"));
            }
        }
        tokens
    };

    {
        let mut mem = vec![0; RAM_SIZE];
        copy_vec(&mut mem, &tokens);
        let (mut tx, rx) = mpsc::channel(2);
        let (ntx, mut nrx) = mpsc::channel(2);
        tx.send(1).await?;
        let thread = tokio::spawn(async move {
            let mut m = Machine::new(0, mem, rx, ntx);
            let _ = m.execute().await;
        });
        while let Some(k) = nrx.recv().await {
            println!("{}", k);
        }
        thread.await?;
    }
    {
        let mut mem = vec![0; RAM_SIZE];
        copy_vec(&mut mem, &tokens);
        let (mut tx, rx) = mpsc::channel(2);
        let (ntx, mut nrx) = mpsc::channel(2);
        tx.send(2).await?;
        let thread = tokio::spawn(async move {
            let mut m = Machine::new(0, mem, rx, ntx);
            let _ = m.execute().await;
        });
        while let Some(k) = nrx.recv().await {
            println!("{}", k);
        }
        thread.await?;
    }

    Ok(())
}
