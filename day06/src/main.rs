use hashbrown::{HashMap, HashSet};
use std::collections::VecDeque;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Orbits {
    orbits: HashMap<String, HashSet<String>>,
}

impl Orbits {
    fn new() -> Self {
        Self {
            orbits: HashMap::new(),
        }
    }

    fn add(&mut self, a: String, b: String) {
        self.orbits.entry(a).or_insert(HashSet::new()).insert(b);
    }

    fn checksum(&self, node: &str, depth: usize) -> usize {
        match self.orbits.get(node) {
            Some(k) => depth + k.iter().map(|i| self.checksum(i, depth + 1)).sum::<usize>(),
            None => depth,
        }
    }

    fn get_path(&self, from: &str, to: &str) -> Option<VecDeque<&String>> {
        if from == to {
            return Some(VecDeque::new());
        }

        let (k, v) = match self.orbits.get_key_value(from) {
            Some(k) => k,
            None => return None,
        };

        if v.contains(to) {
            let mut d = VecDeque::with_capacity(1);
            d.push_front(k);
            return Some(d);
        }

        for i in v.iter() {
            match self.get_path(i, to) {
                Some(mut d) => {
                    d.push_front(k);
                    return Some(d);
                }
                None => (),
            }
        }
        None
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let orbits = {
        let mut orbits = Orbits::new();
        for line in reader.lines() {
            let l = line.expect("Failed to read line");
            let k: Vec<_> = l.split(")").collect();
            orbits.add(String::from(k[0]), String::from(k[1]));
        }
        orbits
    };

    println!("{}", orbits.checksum("COM", 0));

    let path1 = match orbits.get_path("COM", "YOU") {
        Some(k) => k,
        None => {
            eprintln!("failed to find path from COM to YOU");
            return;
        }
    };
    let path2 = match orbits.get_path("COM", "SAN") {
        Some(k) => k,
        None => {
            eprintln!("failed to find path from COM to SAN");
            return;
        }
    };
    for (i, (x, y)) in path1.iter().zip(path2.iter()).enumerate() {
        if x != y {
            println!("{}", path1.len() + path2.len() - i - i);
            break;
        }
    }
}
