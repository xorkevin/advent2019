const PUZZLE_INPUT_MIN: i32 = 231832;
const PUZZLE_INPUT_MAX: i32 = 767346;

fn is_valid_pass(pass: i32) -> bool {
    let (_, double, decr) =
        pass.to_string()
            .chars()
            .fold((None, false, false), |(prev, has_double, decr), i| {
                if let Some(p) = prev {
                    (Some(i), has_double || i == p, decr || i < p)
                } else {
                    (Some(i), has_double, decr)
                }
            });
    double && !decr
}

fn is_valid_pass2(pass: i32) -> bool {
    let (_, run, has_run2) =
        pass.to_string()
            .chars()
            .fold((None, 0, false), |(prev, run, has_run2), i| {
                if let Some(p) = prev {
                    if i == p {
                        (Some(i), run + 1, has_run2)
                    } else {
                        (Some(i), 1, has_run2 || run == 2)
                    }
                } else {
                    (Some(i), 1, has_run2)
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
