extern crate chrono;

use std::time::Duration;
use std::process::Command;
use std::thread;

use chrono::Utc;

const DURATION_MS: Duration = Duration::from_millis(1000);
const SEPARATOR: &str = "   ";

fn main() {
    loop {
        update_status_bar(&aggregate_values());
        thread::sleep(DURATION_MS);
    }
}

fn get_time() -> String {
    Utc::now().format("%H:%MZ").to_string()
}

fn aggregate_values() -> String {
    let time_value = get_time();
    let values = vec![time_value];

    values.join(SEPARATOR)
}

fn update_status_bar(text: &str) {
    let padded_text = &format!(" {} ", text);
    let _ = Command::new("xsetroot")
        .args(&["-name", padded_text])
        .output();
}
