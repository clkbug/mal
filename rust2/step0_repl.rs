use std::io::{self, Write};

fn read() -> String {
  let stdin = io::stdin();
  print!("user> ");
  io::stdout().flush().unwrap();

  let mut s = String::new();
  stdin.read_line(&mut s).expect("ERROR: can't read");

  s
}

fn eval(s: String) -> String {
  s
}

fn print(s: String) {
  print!("{}", s)
}

fn main() {
  loop {
    let s = read();
    if s.is_empty() {
      break;
    }

    let r = eval(s);
    print(r);
  }
  println!();
}
