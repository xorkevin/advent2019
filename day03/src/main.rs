use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

enum Dir {
    Up,
    Down,
    Right,
    Left,
}

struct Segment {
    dir: Dir,
    len: usize,
}

impl Segment {
    const fn new(dir: Dir, len: usize) -> Self {
        Self { dir: dir, len: len }
    }
}

#[derive(PartialEq, Eq, Hash)]
struct Point {
    x: i32,
    y: i32,
}

impl Point {
    const fn new(x: i32, y: i32) -> Self {
        Self { x: x, y: y }
    }

    const fn dist(&self, other: &Point) -> i32 {
        (self.x - other.y).abs() + (self.y - other.y).abs()
    }
}

const ZERO_POINT: Point = Point::new(0, 0);

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let (wire1, wire2) = {
        let mut wire1 = Vec::new();
        let mut wire2 = Vec::new();
        let mut first = true;
        for line in reader.lines() {
            for i in line.expect("Failed to read line").split(",") {
                let d = match &i[0..1] {
                    "U" => Dir::Up,
                    "D" => Dir::Down,
                    "R" => Dir::Right,
                    "L" => Dir::Left,
                    _ => panic!("illegal direction"),
                };
                let l = i[1..].parse::<usize>().expect("Failed to parse num");
                if first {
                    wire1.push(Segment::new(d, l));
                } else {
                    wire2.push(Segment::new(d, l));
                }
            }
            if first {
                first = false;
            }
        }
        (wire1, wire2)
    };

    let (m1, _, _, _) = wire1.iter().fold(
        (HashMap::new(), 0, 0, 0),
        |(mut m, mut x, mut y, mut n), Segment { dir, len }| {
            for _ in 0..*len {
                match dir {
                    Dir::Up => y -= 1,
                    Dir::Down => y += 1,
                    Dir::Right => x += 1,
                    Dir::Left => x -= 1,
                };
                n += 1;
                m.insert(Point::new(x, y), n);
            }
            (m, x, y, n)
        },
    );

    let (dist1, dist2, _, _, _) = wire2.iter().fold(
        (None, None, 0, 0, 0),
        |(mut dist1, mut dist2, mut x, mut y, mut n), Segment { dir, len }| {
            for _ in 0..*len {
                match dir {
                    Dir::Up => y -= 1,
                    Dir::Down => y += 1,
                    Dir::Right => x += 1,
                    Dir::Left => x -= 1,
                };
                n += 1;
                let curr = Point::new(x, y);
                if let Some(w1) = m1.get(&curr) {
                    if let Some(d) = dist1 {
                        let k = curr.dist(&ZERO_POINT);
                        if k < d {
                            dist1 = Some(k);
                        }
                    } else {
                        dist1 = Some(curr.dist(&ZERO_POINT));
                    }

                    if let Some(d) = dist2 {
                        let k = w1 + n;
                        if k < d {
                            dist2 = Some(k);
                        }
                    } else {
                        dist2 = Some(w1 + n);
                    }
                }
            }
            (dist1, dist2, x, y, n)
        },
    );
    if let Some(d) = dist1 {
        println!("{}", d);
    }
    if let Some(d) = dist2 {
        println!("{}", d);
    }
}
