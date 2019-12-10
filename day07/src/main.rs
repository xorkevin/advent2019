use permutohedron;
use std::error::Error;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use tokio::sync::mpsc;

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
    id: u32,
    pc: usize,
    mem: Vec<i32>,
    inp: mpsc::Receiver<i32>,
    out: mpsc::Sender<i32>,
    out_gauge: i32,
}

impl Machine {
    fn new(id: u32, mem: Vec<i32>, inp: mpsc::Receiver<i32>, out: mpsc::Sender<i32>) -> Self {
        Self {
            id,
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

    async fn recv_inp(&mut self) -> Option<i32> {
        self.inp.recv().await
    }

    async fn send_out(&mut self, out: i32) -> Result<(), mpsc::error::SendError<i32>> {
        self.out_gauge = out;
        self.out.send(out).await
    }

    async fn exec(&mut self) -> Result<bool, Box<dyn Error>> {
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
                self.mem[arg1] = match self.recv_inp().await {
                    Some(k) => k,
                    None => {
                        return Err("Failed to read".into());
                    }
                };
                self.step_pc(2);
                Ok(true)
            }
            4 => {
                let arg1 = self.get_arg(1);
                match self.send_out(self.eval_arg(a1, arg1)).await {
                    Ok(_) => (),
                    Err(e) => {
                        return Err(format!("Failed to send: {}", e).into());
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

#[tokio::main(threaded_scheduler)]
async fn main() -> Result<(), Box<dyn Error>> {
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
        let mut max = -1;
        let mut permutations = (0..5).collect::<Vec<_>>();
        let permutations = permutohedron::Heap::new(&mut permutations);
        for phases in permutations {
            let mut k = 0;
            for phase in phases.into_iter() {
                let (mut tx, rx) = mpsc::channel(2);
                let (ntx, mut nrx) = mpsc::channel(2);
                let mut m = Machine::new(0, tokens.clone(), rx, ntx);
                tx.send(phase).await?;
                tx.send(k).await?;
                m.execute().await?;
                k = nrx.recv().await.unwrap();
            }
            if k > max {
                max = k;
            }
        }
        println!("{}", max);
    }

    {
        let mut max = -1;
        let mut permutations = (5..10).collect::<Vec<_>>();
        let permutations = permutohedron::Heap::new(&mut permutations);
        for phases in permutations {
            let len = phases.len();
            let mut send_chans = Vec::with_capacity(len);
            let mut recv_chans = Vec::with_capacity(len);
            for _ in 0..len {
                let (s, r) = mpsc::channel(2);
                send_chans.push(s);
                recv_chans.push(r);
            }
            recv_chans.rotate_left(1);

            for (i, phase) in phases.into_iter().enumerate() {
                let prev = (i + len - 1) % len;
                send_chans[prev].send(phase).await?;
            }
            send_chans[len - 1].send(0).await?;

            let mut threads = Vec::with_capacity(len);
            for (i, (send, recv)) in send_chans
                .into_iter()
                .zip(recv_chans.into_iter())
                .enumerate()
            {
                let tokens = tokens.clone();
                threads.push(tokio::spawn(async move {
                    let mut m = Machine::new(i as u32, tokens, recv, send);
                    let _ = m.execute().await;
                    m.out_gauge
                }));
            }

            let mut k = -1;
            for t in threads.iter_mut() {
                k = t.await?;
            }

            if k > max {
                max = k;
            }
        }
        println!("{}", max);
    }

    Ok(())
}
