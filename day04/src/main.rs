const PUZZLE_INPUT_MIN: i32 = 231832;
const PUZZLE_INPUT_MAX: i32 = 767346;

fn pass_to_digits(pass: i32) -> [i32; 6] {
    [
        pass / 100000 % 10,
        pass / 10000 % 10,
        pass / 1000 % 10,
        pass / 100 % 10,
        pass / 10 % 10,
        pass % 10,
    ]
}

fn is_valid_pass(pass: i32) -> bool {
    let (_, double, decr) = pass_to_digits(pass)
        .into_iter()
        .fold((-1, false, false), |(prev, has_double, decr), &i| {
            (i, has_double || i == prev, decr || i < prev)
        });
    double && !decr
}

fn is_valid_pass2(pass: i32) -> bool {
    let (_, run, has_run2) =
        pass_to_digits(pass)
            .into_iter()
            .fold((-1, 0, false), |(prev, run, has_run2), &i| {
                if i == prev {
                    (i, run + 1, has_run2)
                } else {
                    (i, 1, has_run2 || run == 2)
                }
            });
    (has_run2 || run == 2) && is_valid_pass(pass)
}

fn main() {
    let k1 =
        (PUZZLE_INPUT_MIN..=PUZZLE_INPUT_MAX)
            .fold(0, |s, i| if is_valid_pass(i) { s + 1 } else { s });
    println!("{}", k1);
    let k2 =
        (PUZZLE_INPUT_MIN..=PUZZLE_INPUT_MAX)
            .fold(0, |s, i| if is_valid_pass2(i) { s + 1 } else { s });
    println!("{}", k2);
}
