// Potentially use later
// use std::io::{self, BufRead};
// use std::io::{self, Write};
// use serde_json::{self, Value};
// add more cargo bins?

use std::env;
use std::fs;
use std::path::Path;
use std::io;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct Workout {
    date: String,
    exercise: String,
    weight: f64,
    reps: u32,
    #[serde(skip_serializing_if = "Option::is_none")]
    estimated_1rm: Option<f64>,
}

fn epley_1rm(weight: f64, reps: u32) -> f64 {
    if reps <= 1 {
        weight
    } else {
        weight * ( 1.0 + (reps as f64 / 30.0 ))
    }
}

fn main() -> io::Result <()> {
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        eprintln!("Usage: stats <path_to_workouts.json>");
        std::process::exit(1);
    }
    
    let path = Path::new(&args[1]);
    let file = fs::File::open(path)?;
    let reader = io::BufReader::new(file);

    let mut workouts: Vec<Workout> = serde_json::from_reader(reader)?;

    for workout in &mut workouts {
        let one_rm = epley_1rm(workout.weight, workout.reps);
        workout.estimated_1rm = Some(one_rm.round()); //round to whole number
    }

    // For V1, I am just printing the result, I will add writing back or returning JSON later
    println!("Stats Calculation");
    for w in &workouts {
        println!(
            "{}  |  {}  | {} lb x {} reps -> est. 1RM: {:.0} lb",
            w.date, w.exercise, w.weight, w.reps, w.estimated_1rm.unwrap()
        );
    }

    Ok(())
}
