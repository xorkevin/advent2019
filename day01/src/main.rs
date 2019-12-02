use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Fuel {
    i: i32,
}

impl Fuel {
    fn new(i: i32) -> Self {
        return Self { i: i };
    }
}

impl Iterator for Fuel {
    type Item = i32;

    fn next(&mut self) -> Option<i32> {
        self.i = self.i / 3 - 2;
        Some(self.i)
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        for line in reader.lines() {
            nums.push(
                line.expect("Failed to read line")
                    .parse::<i32>()
                    .expect("Failed to parse num"),
            );
        }
        nums
    };

    let sum1 = nums.iter().fold(0, |k, i| k + i / 3 - 2);
    println!("{}", sum1);

    let sum2 = nums.iter().fold(0, |k, &i| {
        k + Fuel::new(i).take_while(|&x| (x / 3 - 2) > 0).sum::<i32>()
    });
    println!("{}", sum2);
}
